package utils

import (
	"context"
	"cu_coworking_book/go/internal/db"
	"cu_coworking_book/go/internal/pb"
	"errors"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Сверяет тип ошибки и возвращает NotFound или Internal.
func CompareErrAndErrNotFound(err error) error {
	if errors.Is(err, db.ErrNotFound) {
		return status.Error(codes.NotFound, err.Error())
	} else {
		return status.Error(codes.Internal, err.Error())
	}
}

// Сверяет пароль указанного пользователя, а также проверяет роль при выполнении
// команд администратора.
func ComparePassword(ctx context.Context, email string, password string, admin bool) (db.User, error) {
	user, err := db.GetUserByEmail(ctx, email)

	if err != nil {
		return user, CompareErrAndErrNotFound(err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	adminPermission := admin && (user.Role != pb.Role_ADMIN.String() && user.Role != pb.Role_MAIN_ADMIN.String())

	if err != nil || adminPermission {
		return user, status.Error(codes.PermissionDenied, "доступ запрещен")
	}
	return user, nil
}

// Преобразовывает db.User в pb.GetResponse, а также возвращает ФИО с Большой Буквы.
func ParseToResult(res db.User) *pb.GetResponse {

	c := cases.Title(language.Russian)

	return &pb.GetResponse{
		Id:         res.ID,
		Email:      res.Email,
		FirstName:  c.String(res.FirstName),
		SecondName: c.String(res.SecondName),
		Patronymic: c.String(res.Patronymic),
		Password:   res.Password,
		Role:       pb.Role(pb.Role_value[res.Role]),
	}
}
