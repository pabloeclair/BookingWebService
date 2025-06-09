package auth

import (
	"context"
	"cu_coworking_book/go/internal/db"
	"cu_coworking_book/go/internal/pb"
	"cu_coworking_book/go/internal/utils"
	"errors"
	"log"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServer struct {
	pb.UnimplementedAuthenticationServer
}

// Записывает логи применения всех хэндлеров. Если какой-то из хэндлеров запустился неудачно, то сообщает о типе и сообщении ошибки.
func LogInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Printf("Auth server: %s", info.FullMethod)

	resp, err := handler(ctx, req)

	if err != nil {
		log.Printf("Error: %v", err)
	}

	return resp, err
}

// Регистрирует нового пользователя.
func (s *AuthServer) SignupUser(ctx context.Context, req *pb.SignupRequest) (*pb.Id, error) {
	// Первый пользователь автоматически становится MAIN_ADMIN, а остальные последующие — USER.

	// - Если пользователь с указанной почтой уже существует, то вернется ошибка AlreadyExists.

	// При любых других ошибках – Internal.
	// Показатель успеха — id пользователя.

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	user := db.User{
		Email:      req.Email,
		FirstName:  strings.ToLower(req.FirstName),
		SecondName: strings.ToLower(req.SecondName),
		Patronymic: strings.ToLower(req.GetPatronymic()),
		Password:   string(hashPassword),
		Role:       pb.Role_USER.String(),
	}

	id, err := db.CreateUser(ctx, user)
	if err != nil {
		if errors.Is(err, db.ErrBadRequest) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.Id{Id: id}, nil
}

// Сверяет пароль пользователя и при успехе возвращает информацию о нем.
func (s *AuthServer) LoginUser(ctx context.Context, req *pb.Email) (*pb.GetResponse, error) {
	// - Если пользователь с указанной почтой не найден, вернется ошибка NotFound.
	// - Если пароль пользователя не совпал — Unauthenticated.

	// При любых других ошибках – Internal.
	// Показатель успеха — информация о пользователе.

	res, err := utils.ComparePassword(ctx, req.Email, req.Password, false)
	if err != nil {
		return nil, err
	}

	return utils.ParseToResult(res), nil
}

// Обновляет информацию о пользователе по заданным полям, кроме пароля. Пароль обновляет хэндлер UpdatePassword.
func (s *AuthServer) UpdateUser(ctx context.Context, req *pb.UpdateRequest) (*pb.Empty, error) {
	// - Если пользователь с указанной id или почтой не найден, вернется ошибка NotFound.
	// - Если новая почта изменяемого пользователя уже существует — BadRequest.
	// - Если пароль пользователя не совпал — Unauthenticated.

	// При любых других ошибках – Internal.
	// Показатель успеха — отсутствие ошибки.

	_, err := utils.ComparePassword(ctx, req.Email, req.Password, false)
	if err != nil {
		return nil, err
	}

	user := db.User{
		ID:         req.Id,
		Email:      req.Email,
		FirstName:  strings.ToLower(req.FirstName),
		SecondName: strings.ToLower(req.SecondName),
		Patronymic: strings.ToLower(req.GetPatronymic()),
	}

	if err := db.UpdateUser(ctx, user); err != nil {
		if errors.Is(err, db.ErrBadRequest) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, utils.CompareErrAndErrNotFound(err)
	}

	return nil, nil
}

// Обновляет пароль пользователя.
func (s *AuthServer) UpdatePassword(ctx context.Context, req *pb.UpdatePasswordRequest) (*pb.Empty, error) {
	// - Если пользователь с указанной id или почтой не найден, вернется ошибка NotFound.
	// - Если пароль пользователя не совпал — Unauthenticated.

	// При любых других ошибках – Internal.
	// Показатель успеха — отсутствие ошибки.

	if _, err := utils.ComparePassword(ctx, req.Email, req.OldPassword, false); err != nil {
		return nil, err
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), 12)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if err = db.UpdatePassword(ctx, req.Id, string(hashPassword)); err != nil {
		return nil, utils.CompareErrAndErrNotFound(err)
	}
	return nil, nil
}

// Удаляет пользователя из базы данных.
func (s *AuthServer) DeleteUser(ctx context.Context, req *pb.Email) (*pb.Empty, error) {
	// - Если пользователь с указанным id или почтой не найден, вернется ошибка NotFound.
	// - Если пароль пользователя не совпал — Unauthenticated.

	// При любых других ошибках – Internal.
	// Показатель успеха — отсутствие ошибки.

	if _, err := utils.ComparePassword(ctx, req.Email, req.Password, false); err != nil {
		return nil, err
	}

	if err := db.DeleteUser(ctx, req.Id); err != nil {
		return nil, utils.CompareErrAndErrNotFound(err)
	}
	return nil, nil
}
