package admin

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

type AdminService struct {
	pb.UnimplementedAdminServiceServer
	mu sync.RWMutex
}

func LogInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Printf("Admin server: %s", info.FullMethod)

	resp, err := handler(ctx, req)

	if err != nil {
		log.Printf("Error: %v", err)
	}

	return resp, err
}

func (s *AdminService) CreateUser(ctx context.Context, req *pb.CreateRequestAdmin) (*pb.UserResponse, error) {

	s.mu.RLock()
	_, err := utils.ComparePassword(req.AdminEmail, req.AdminPassword, true)
	s.mu.RUnlock()

	if err != nil {
		return nil, err
	}

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

	return utils.ParseToResult(res), nil
}

func (s *AdminService) GetUser(ctx context.Context, req *pb.GetRequestAdmin) (*pb.GetUserResponseAdmin, error) {

	s.mu.RLock()
	_, err := utils.ComparePassword(req.AdminEmail, req.AdminPassword, true)
	s.mu.RUnlock()
	if err != nil {
		return nil, err
	}

	s.mu.RLock()
	res, err := db.GetUserByKey(&req.SortBy, req.SortKey)
	s.mu.RUnlock()
	if err != nil {
		return nil, utils.CompareErrAndErrNotFound(err)
	}

	var users []*pb.UserResponse
	for r := range res {
		users = append(users, utils.ParseToResult(res[r]))
	}

	return &pb.GetUserResponseAdmin{Users: users}, nil
}

func (s *AdminService) UpdateUser(ctx context.Context, req *pb.UpdateRequestAdmin) (*pb.UserResponse, error) {

	s.mu.RLock()
	_, err := utils.ComparePassword(req.AdminEmail, req.AdminPassword, true)
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
	res, err := db.UpdateUser(req.GetId(), user)
	s.mu.Unlock()
	if err != nil {
		if errors.Is(err, db.ErrBadRequest) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, utils.CompareErrAndErrNotFound(err)
	}

	return utils.ParseToResult(res), nil
}

func (s *AdminService) DeleteUser(ctx context.Context, req *pb.DeleteRequestAdmin) (*pb.Empty, error) {

	s.mu.RLock()
	_, err := utils.ComparePassword(req.AdminEmail, req.AdminPassword, true)
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
