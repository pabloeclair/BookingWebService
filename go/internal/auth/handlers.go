package auth

import (
	"context"
	"cu_coworking_book/go/internal/db"
	"cu_coworking_book/go/internal/pb"
	"database/sql"
	"encoding/base64"
	"errors"
	"log"
	"strings"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServer struct {
	pb.UnimplementedAuthenticationServer
	mu sync.RWMutex
}

func LogInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Printf("Server call: %s", info.FullMethod)

	resp, err := handler(ctx, req)

	if err != nil {
		log.Printf("RPC failed with error: %v", err)
	}

	return resp, err
}

func parseToResult(res db.User) *pb.UserResponse {

	return &pb.UserResponse{
		Id:         res.ID,
		Email:      res.Email,
		FirstName:  res.FirstName,
		SecondName: res.SecondName,
		Patronymic: &res.Patronymic,
		AuthHeader: GenerateAuthHeader(res.Email, res.Password),
		Role:       pb.Role(pb.Role_value[res.Role]),
	}
}

func GenerateAuthHeader(email string, password string) string {

	emailPassword := []byte(email + ":" + password)
	authHeader := base64.StdEncoding.EncodeToString(emailPassword)
	return authHeader
}

func (s *AuthServer) SignupUser(ctx context.Context, req *pb.SignupRequest) (*pb.UserResponse, error) {

	user := db.User{
		Email:      req.GetEmail(),
		FirstName:  req.GetFirstName(),
		SecondName: req.GetSecondName(),
		Patronymic: req.GetPatronymic(),
		Password:   req.GetPassword(),
	}

	s.mu.Lock()
	res, err := db.CreateUser(user)
	s.mu.Unlock()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "sign up error: db error: %v", err)
	}

	return parseToResult(res), nil
}

func (s *AuthServer) GetUserByEmail(ctx context.Context, req *pb.Email) (*pb.UserResponse, error) {

	s.mu.RLock()
	res, err := db.GetUserByEmail(req.GetEmail())
	s.mu.RUnlock()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "log in error: db error: %v", err)
		} else {
			return nil, status.Errorf(codes.Internal, "log in error: db error: %v", err)
		}
	}

	return parseToResult(res), nil
}

func (s *AuthServer) GetUserById(ctx context.Context, req *pb.Id) (*pb.UserResponse, error) {

	s.mu.RLock()
	res, err := db.GetUserById(req.GetId())
	s.mu.RUnlock()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "log in error: db error: %v", err)
		} else {
			return nil, status.Errorf(codes.Internal, "log in error: db error: %v", err)
		}
	}

	return parseToResult(res), nil
}

func (s *AuthServer) UpdateUser(ctx context.Context, req *pb.UpdateRequest) (*pb.UserResponse, error) {

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
	res, err := db.UpdateUser(oldEmail[0], user)
	s.mu.Unlock()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "update user error: db error: %v", err)
	}

	return parseToResult(res), nil
}

func (s *AuthServer) DeleteUser(ctx context.Context, req *pb.Email) (*pb.Empty, error) {

	s.mu.Lock()
	err := db.DeleteUser(req.Email)
	s.mu.Unlock()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "delete user error: db error: %v", err)
	}
	return nil, nil
}
