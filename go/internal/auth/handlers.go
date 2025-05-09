package auth

import (
	"context"
	"cu_coworking_book/go/internal/db"
	"cu_coworking_book/go/internal/pb"
	"database/sql"
	"encoding/base64"
	"errors"
	"strings"
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

	password, err := bcrypt.GenerateFromPassword([]byte(req.GetPassword()), 14)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "sign up error: generate password hash: %v", err)
	}

	user := db.User{
		Email:      req.GetEmail(),
		FirstName:  req.GetFirstName(),
		SecondName: req.GetSecondName(),
		Patronymic: req.GetPatronymic(),
		Password:   string(password),
	}

	s.mu.Lock()
	res, err := db.CreateUser(user)
	s.mu.Unlock()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "sign up error: db error: %v", err)
	}

	emailPassword := req.GetEmail() + ":" + req.GetPassword()
	authHeader := base64.StdEncoding.EncodeToString([]byte(emailPassword))

	return &pb.SignUpResponse{Id: res.ID, Role: res.Role, AuthHeader: authHeader}, nil
}

func (s *AuthServer) LogInUser(ctx context.Context, req *pb.EmailPassword) (*pb.AuthHeader, error) {

	s.mu.RLock()
	expectedLogin, err := db.GetUserByEmail(req.GetEmail())
	s.mu.RUnlock()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "log in error: db error: %v", err)
		} else {
			return nil, status.Errorf(codes.Internal, "log in error: db error: %v", err)
		}
	}

	if expectedLogin.Password == req.GetPassword() {
		authHeader := base64.StdEncoding.EncodeToString([]byte(req.GetEmail() + ":" + req.GetPassword()))
		return &pb.AuthHeader{AuthHeader: authHeader}, nil
	} else {
		return nil, status.Error(codes.InvalidArgument, "log in error: invalid password")
	}
}

func (s *AuthServer) UpdateUser(ctx context.Context, req *pb.UpdateRequest) (*pb.AuthHeader, error) {

	user := db.User{
		Email:      req.GetEmail(),
		FirstName:  req.GetFirstName(),
		SecondName: req.GetSecondName(),
		Patronymic: req.GetPatronymic(),
		Password:   req.GetPassword(),
	}

	emailPassword, err := base64.StdEncoding.DecodeString(req.GetAuthHeader())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "update user error: invalid auth header")
	}
	oldEmail := strings.Split(string(emailPassword), ":")

	s.mu.Lock()
	err = db.UpdateUser(oldEmail[0], user)
	s.mu.Unlock()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "update user error: db error: %v", err)
	}

	newAuthHeader := base64.StdEncoding.EncodeToString([]byte(req.Email + ":" + req.Password))
	return &pb.AuthHeader{AuthHeader: newAuthHeader}, nil
}
