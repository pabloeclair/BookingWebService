package auth

import (
	"context"
	"cu_coworking_book/go/internal/db"
	"cu_coworking_book/go/internal/pb"
	"encoding/base64"
	"sync"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServer struct {
	pb.UnimplementedAuthenticationServer
	mu sync.RWMutex
}

func (s *AuthServer) SignUpUser(ctx context.Context, req *pb.SignUpRequest) (*pb.SignUpResponse, error) {

	var result *pb.SignUpResponse

	password, err := bcrypt.GenerateFromPassword([]byte(req.GetPassword()), 14)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "sign up error: generate password hash: %v", err)
	}

	user := db.User{
		Email:          req.GetEmail(),
		FirstName:      req.GetFirstName(),
		SecondName:     req.GetSecondName(),
		Patronymic:     req.GetPatronymic(),
		HashedPassword: string(password),
	}

	s.mu.Lock()
	res, err := db.CreateUser(user)
	s.mu.Unlock()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "sign up error: generate password hash: %v", err)
	}
	result.Id = res.ID
	result.Role = res.Role

	emailPassword := req.GetEmail() + ":" + req.GetPassword()
	result.AuthHeader = base64.StdEncoding.EncodeToString([]byte(emailPassword))

	return result, nil
}
