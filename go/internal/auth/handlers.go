package auth

import (
	"context"
	"cu_coworking_book/go/internal/db"
	"cu_coworking_book/go/internal/pb"
	"cu_coworking_book/go/internal/utils"
	"errors"
	"log"
	"sync"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServer struct {
	pb.UnimplementedAuthenticationServer
	mu sync.RWMutex
}

func LogInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Printf("Auth server: %s", info.FullMethod)

	resp, err := handler(ctx, req)

	if err != nil {
		log.Printf("Error: %v", err)
	}

	return resp, err
}

func (s *AuthServer) SignupUser(ctx context.Context, req *pb.SignupRequest) (*pb.UserResponse, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	user := db.User{
		Email:      req.GetEmail(),
		FirstName:  req.GetFirstName(),
		SecondName: req.GetSecondName(),
		Patronymic: req.GetPatronymic(),
		Password:   string(hashedPassword),
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

	return utils.ParseToResult(res), nil
}

func (s *AuthServer) LoginUser(ctx context.Context, req *pb.Email) (*pb.UserResponse, error) {

	s.mu.RLock()
	res, err := utils.ComparePassword(req.GetEmail(), []byte(req.GetPassword()), false)
	s.mu.RUnlock()
	if err != nil {
		return nil, err
	}

	return utils.ParseToResult(res), nil
}

func (s *AuthServer) UpdateUser(ctx context.Context, req *pb.UpdateRequest) (*pb.UserResponse, error) {

	s.mu.RLock()
	_, err := utils.ComparePassword(req.GetEmail(), []byte(req.GetPassword()), false)
	s.mu.RUnlock()
	if err != nil {
		return nil, err
	}

	user := db.User{
		Email:      req.GetEmail(),
		FirstName:  req.GetFirstName(),
		SecondName: req.GetSecondName(),
		Patronymic: req.GetPatronymic(),
	}

	s.mu.Lock()
	res, err := db.UpdateUser(req.GetOldEmail(), user)
	s.mu.Unlock()
	if err != nil {
		if errors.Is(err, db.ErrBadRequest) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, utils.CompareErrAndErrNotFound(err)
	}

	return utils.ParseToResult(res), nil
}

func (s *AuthServer) DeleteUser(ctx context.Context, req *pb.Email) (*pb.Empty, error) {

	s.mu.RLock()
	_, err := utils.ComparePassword(req.GetEmail(), []byte(req.GetPassword()), false)
	s.mu.RUnlock()

	if err != nil {
		return nil, err
	}

	s.mu.Lock()
	err = db.DeleteUser(req.Email)
	s.mu.Unlock()
	if err != nil {
		return nil, utils.CompareErrAndErrNotFound(err)
	}
	return nil, nil
}
