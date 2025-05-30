package admin

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

type AdminService struct {
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

func (s *AdminService) CreateUser(ctx context.Context, req *pb.CreateRequestAdmin) (*pb.UserResponse, error) {

	s.mu.RLock()
	_, err := utils.ComparePassword(req.AdminEmail, []byte(req.AdminPassword), true)
	s.mu.RUnlock()

	if err != nil {
		return nil, err
	}

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
