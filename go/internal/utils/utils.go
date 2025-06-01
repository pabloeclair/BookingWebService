package utils

import (
	"crypto/hmac"
	"cu_coworking_book/go/internal/db"
	"cu_coworking_book/go/internal/pb"
	"encoding/hex"
	"errors"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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

func ComparePassword(email string, password string, admin bool) (db.User, error) {
	user, err := db.GetUserByEmail(email)

	if err != nil {
		return user, CompareErrAndErrNotFound(err)
	}

	hash1, err := hex.DecodeString(password)
	if err != nil {
		return user, status.Error(codes.Internal, "ошибка декодирования sha256 (1)")
	}

	hash2, err := hex.DecodeString(user.Password)
	if err != nil {
		return user, status.Error(codes.Internal, "ошибка декодирования sha256 (2)")
	}

	if !hmac.Equal(hash1, hash2) {
		return user, status.Error(codes.Unauthenticated, "неверный пароль")
	}

	if admin && (user.Role != pb.Role_ADMIN.String() || user.Role != pb.Role_MAIN_ADMIN.String()) {
		return user, status.Error(codes.PermissionDenied, "доступ запрещен")
	}
	return user, nil
}

func ParseToResult(res db.User) *pb.GetResponse {

	c := cases.Title(language.Russian)
	patronymic := c.String(res.Patronymic)

	return &pb.GetResponse{
		Id:         res.ID,
		Email:      res.Email,
		FirstName:  c.String(res.FirstName),
		SecondName: c.String(res.SecondName),
		Patronymic: &patronymic,
		Password:   res.Password,
		Role:       pb.Role(pb.Role_value[res.Role]),
	}
}
