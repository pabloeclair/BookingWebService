package utils

import (
	"cu_coworking_book/go/internal/db"
	"cu_coworking_book/go/internal/pb"
	"errors"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func CompareErrAndErrNotFound(err error) error {
	if errors.Is(err, db.ErrNotFound) {
		return status.Error(codes.NotFound, err.Error())
	} else {
		return status.Error(codes.Internal, err.Error())
	}
}

func ComparePassword(email string, password []byte, admin bool) (db.User, error) {
	user, err := db.GetUserByEmail(email)

	if err != nil {
		return user, CompareErrAndErrNotFound(err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), password); err != nil {
		return user, status.Error(codes.Unauthenticated, "неверный пароль")
	}

	if admin && (user.Role != pb.Role_ADMIN.String() || user.Role != pb.Role_MAIN_ADMIN.String()) {
		return user, status.Error(codes.PermissionDenied, "доступ запрещен")
	}
	return user, nil
}

func ParseToResult(res db.User) *pb.UserResponse {

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
