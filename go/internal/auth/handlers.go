package auth

import (
	"context"
	"cu_coworking_book/go/internal/db"
	"cu_coworking_book/go/internal/pb"
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
		Token:      GenerateToken(res.Email, res.Password, res.Role),
		Role:       pb.Role(pb.Role_value[res.Role]),
	}
}

func GenerateToken(email string, password string, role string) string {

	str := []byte(email + ":" + password + ":" + role)
	authHeader := base64.StdEncoding.EncodeToString(str)
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
		if errors.Is(err, db.ErrBadRequest) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return parseToResult(res), nil
}

func (s *AuthServer) GetUserByEmail(ctx context.Context, req *pb.Email) (*pb.UserResponse, error) {

	s.mu.RLock()
	res, err := db.GetUserByEmail(req.GetEmail())
	s.mu.RUnlock()
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		} else {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return parseToResult(res), nil
}

func (s *AuthServer) GetUserById(ctx context.Context, req *pb.Id) (*pb.UserResponse, error) {

	s.mu.RLock()
	res, err := db.GetUserById(req.GetId())
	s.mu.RUnlock()
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		} else {
			return nil, status.Error(codes.Internal, err.Error())
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

	emailPassword, err := base64.StdEncoding.DecodeString(req.GetToken())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid token")
	}
	oldEmail := strings.Split(string(emailPassword), ":")

	s.mu.Lock()
	res, err := db.UpdateUser(oldEmail[0], user)
	s.mu.Unlock()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return parseToResult(res), nil
}

func (s *AuthServer) DeleteUser(ctx context.Context, req *pb.Email) (*pb.Empty, error) {

	s.mu.Lock()
	err := db.DeleteUser(req.Email)
	s.mu.Unlock()
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		} else {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	return nil, nil
}
