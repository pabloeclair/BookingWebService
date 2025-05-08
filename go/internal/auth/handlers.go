package auth

import (
	"context"
	"cu_coworking_book/go/internal/db"
	"cu_coworking_book/go/internal/pb"
	"database/sql"
	"errors"
	"sync"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServer struct {
	pb.UnimplementedAuthenticationServer
	mu sync.RWMutex
}

func (s *AuthServer) CreateUser(ctx context.Context, req *pb.User) (*pb.UserId, error) {

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(req.GetPswd()), 12)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create user error: generate hash password: %v", err)
	}
	user := db.User{
		Email:      req.GetEmail(),
		FirstName:  req.GetFirstName(),
		SecondName: req.GetSecondName(),
		MidName:    req.GetMidName(),
		Password:   string(hashPassword),
	}

	s.mu.Lock()
	id, err := db.CreateUser(user)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	s.mu.Unlock()

	return &pb.UserId{Id: id}, nil
}

func (s *AuthServer) GetUserByEmail(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {

	email := req.GetEmail()
	user, err := db.GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.InvalidArgument, "get user by email error: user with email %s does not exist", email)
		}
	}
	return user, nil
}
