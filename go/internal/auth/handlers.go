package auth

import (
	"context"
	"cu_coworking_book/go/internal/db"
	"cu_coworking_book/go/internal/pb"
	"cu_coworking_book/go/internal/utils"
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

// Записывает логи применения всех хандлеров. Если какой-то из хандлеров
// запустился  неудачно, то сообщает о типе и сообщении ошибки.
func LogInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Printf("Auth server: %s", info.FullMethod)

	resp, err := handler(ctx, req)

	if err != nil {
		log.Printf("Error: %v", err)
	}

	return resp, err
}

// Хэндлер регистрации нового пользователя. Если пользователь уже
// существует, возвращает ошибку AlreadyExists.
func (s *AuthServer) SignupUser(ctx context.Context, req *pb.SignupRequest) (*pb.Empty, error) {

	user := db.User{
		Email:      req.GetEmail(),
		FirstName:  strings.ToLower(req.GetFirstName()),
		SecondName: strings.ToLower(req.GetSecondName()),
		Patronymic: strings.ToLower(req.GetPatronymic()),
		Password:   req.GetPassword(),
	}

	s.mu.Lock()
	err := db.CreateUser(user)
	s.mu.Unlock()
	if err != nil {
		if errors.Is(err, db.ErrBadRequest) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return nil, nil
}

func (s *AuthServer) LoginUser(ctx context.Context, req *pb.Email) (*pb.GetResponse, error) {

	s.mu.RLock()
	res, err := utils.ComparePassword(req.GetEmail(), req.GetPassword(), false)
	s.mu.RUnlock()
	if err != nil {
		return nil, err
	}

	return utils.ParseToResult(res), nil
}

func (s *AuthServer) UpdateUser(ctx context.Context, req *pb.UpdateRequest) (*pb.Empty, error) {

	s.mu.RLock()
	_, err := utils.ComparePassword(req.GetEmail(), req.GetPassword(), false)
	s.mu.RUnlock()
	if err != nil {
		return nil, err
	}

	user := db.User{
		Email:      req.GetEmail(),
		FirstName:  strings.ToLower(req.GetFirstName()),
		SecondName: strings.ToLower(req.GetSecondName()),
		Patronymic: strings.ToLower(req.GetPatronymic()),
	}

	s.mu.Lock()
	err = db.UpdateUser(req.GetId(), user)
	s.mu.Unlock()
	if err != nil {
		if errors.Is(err, db.ErrBadRequest) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, utils.CompareErrAndErrNotFound(err)
	}

	return nil, nil
}

func (s *AuthServer) UpdatePassword(ctx context.Context, req *pb.UpdatePasswordRequest) (*pb.Empty, error) {

	s.mu.RLock()
	_, err := utils.ComparePassword(req.GetEmail(), req.GetOldPassword(), false)
	s.mu.RUnlock()
	if err != nil {
		return nil, err
	}

	s.mu.Lock()
	err = db.UpdatePassword(req.GetId(), req.GetNewPassword())
	s.mu.Unlock()

	if err != nil {
		return nil, utils.CompareErrAndErrNotFound(err)
	}
	return nil, nil
}

func (s *AuthServer) DeleteUser(ctx context.Context, req *pb.Email) (*pb.Empty, error) {

	s.mu.RLock()
	_, err := utils.ComparePassword(req.GetEmail(), req.GetPassword(), false)
	s.mu.RUnlock()

	if err != nil {
		return nil, err
	}

	s.mu.Lock()
	err = db.DeleteUser(req.Id)
	s.mu.Unlock()
	if err != nil {
		return nil, utils.CompareErrAndErrNotFound(err)
	}
	return nil, nil
}
