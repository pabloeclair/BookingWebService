package auth

import (
	"context"
	"cu_coworking_book/go/internal/db"
	"cu_coworking_book/go/internal/pb"
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
	log.Printf("Server call: %s", info.FullMethod)

	resp, err := handler(ctx, req)

	if err != nil {
		log.Printf("RPC failed with error: %v", err)
	}

	return resp, err
}

func compareErrAndErrNotFound(err error) error {
	if errors.Is(err, db.ErrNotFound) {
		return status.Error(codes.NotFound, err.Error())
	} else {
		return status.Error(codes.Internal, err.Error())
	}
}

func comparePassword(email string, password []byte) (db.User, error) {
	user, err := db.GetUserByEmail(email)

	if err != nil {
		return user, compareErrAndErrNotFound(err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), password); err != nil {
		return user, status.Error(codes.Unauthenticated, "password didn't compare")
	}
	return user, nil
}

func parseToResult(res db.User) *pb.UserResponse {

	return &pb.UserResponse{
		Id:         res.ID,
		Email:      res.Email,
		FirstName:  res.FirstName,
		SecondName: res.SecondName,
		Patronymic: &res.Patronymic,
		Password:   res.Password,
		Role:       pb.Role(pb.Role_value[res.Role]),
	}
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

	return parseToResult(res), nil
}

func (s *AuthServer) GetUserByEmail(ctx context.Context, req *pb.Email) (*pb.UserResponse, error) {

	s.mu.RLock()
	res, err := comparePassword(req.GetEmail(), []byte(req.GetPassword()))
	s.mu.RUnlock()
	if err != nil {
		return nil, err
	}

	return parseToResult(res), nil
}

func (s *AuthServer) GetUserById(ctx context.Context, req *pb.Id) (*pb.UserResponse, error) {

	s.mu.RLock()
	res, errGet := db.GetUserByEmail(req.GetAdminEmail())
	s.mu.RUnlock()

	if !errors.Is(errGet, db.ErrNotFound) && errGet != nil {
		return nil, status.Error(codes.Internal, errGet.Error())
	}

	err := bcrypt.CompareHashAndPassword([]byte(res.Password), []byte(req.AdminPassword))
	if err != nil || errors.Is(errGet, db.ErrNotFound) || (res.Role != "ADMIN" && res.Role != "MAIN_ADMIN") {
		return nil, status.Error(codes.PermissionDenied, "access denied")
	}

	s.mu.RLock()
	res, err = db.GetUserById(req.GetId())
	s.mu.RUnlock()
	if err != nil {
		return nil, compareErrAndErrNotFound(err)
	}

	return parseToResult(res), nil
}

func (s *AuthServer) UpdateUser(ctx context.Context, req *pb.UpdateRequest) (*pb.UserResponse, error) {

	s.mu.RLock()
	_, err := comparePassword(req.GetEmail(), []byte(req.GetPassword()))
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
		return nil, compareErrAndErrNotFound(err)
	}

	return parseToResult(res), nil
}

func (s *AuthServer) DeleteUser(ctx context.Context, req *pb.Email) (*pb.Empty, error) {

	s.mu.RLock()
	_, err := comparePassword(req.GetEmail(), []byte(req.GetPassword()))
	s.mu.RUnlock()

	if err != nil {
		return nil, err
	}

	s.mu.Lock()
	err = db.DeleteUser(req.Email)
	s.mu.Unlock()
	if err != nil {
		return nil, compareErrAndErrNotFound(err)
	}
	return nil, nil
}
