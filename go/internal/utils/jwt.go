package utils

import (
	"context"
	"cu_coworking_book/go/internal/db"
	"cu_coworking_book/go/internal/models"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
)

var (
	DefaultUserJWTDuration  int   = 60
	DefaultAdminJWTDuration int   = 60
	ErrNotFoundSecretKey    error = errors.New("отсутствует секретный ключ для JWT")
	ErrNotFoundJWTDuration  error = errors.New("пропущено значение времени действия токена")
	ErrInvalidJWTDuration   error = errors.New("значение времени действия токена некорректно; должно быть целым числом")
	ErrInvalidRole          error = errors.New("неизвестная роль пользователя")
)

// Генерирует JWT-токен указанному пользователю
func GenerateJWT(ctx context.Context, email string) (string, error) {
	secretKey := []byte(os.Getenv("JWT_SECRET_KEY"))
	if secretKey == nil {
		return "", fmt.Errorf("ошибка окружения: %w", ErrNotFoundSecretKey)
	}

	res, err := db.GetUserByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	var durationInt int
	if res.Role == models.Role_USER.String() {
		if durationInt, err = CheckEnvUserJWTDuration(); err != nil && !errors.Is(err, ErrNotFoundJWTDuration) {
			return "", err
		}
	} else if res.Role == models.Role_ADMIN.String() || res.Role == models.Role_MAIN_ADMIN.String() {
		if durationInt, err = CheckEnvAdminJWTDuration(); err != nil && !errors.Is(err, ErrNotFoundJWTDuration) {
			return "", err
		}
	} else {
		return "", fmt.Errorf(
			"jwt error: %w; при попытке сгенерировать JWT для пользователя с email %s",
			ErrInvalidRole, email)
	}

	userClaim := models.UserClaim{
		Id:         res.Id,
		Email:      res.Email,
		FirstName:  res.FirstName,
		SecondName: res.SecondName,
		Patronymic: res.Patronymic,
		Role:       res.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * time.Duration(durationInt)).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaim)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("jwt error: %w", err)
	}
	return tokenString, nil
}

func ParseJWT(tokenString string) (*models.UserClaim, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.UserClaim{}, func(t *jwt.Token) (interface{}, error) {
		secretKey := []byte(os.Getenv("JWT_SECRET_KEY"))
		if secretKey == nil {
			return nil, ErrNotFoundSecretKey
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}
	return token.Claims.(*models.UserClaim), nil
}

// Проверяет наличие и корректность переменных окружения, указывающих продолжительность
// JWT-токенов и в случае успеха возвращает указанное число или число по умолчанию (60) при отсутствии
func CheckEnvUserJWTDuration() (int, error) {
	duration := os.Getenv("JWT_USER_DURATION")

	if duration == "" {
		return DefaultUserJWTDuration, fmt.Errorf(
			"ошибка окружения: в JWT_USER_DURATION %w, из-за чего использовано значение по умолчанию – %d",
			ErrNotFoundJWTDuration, DefaultUserJWTDuration)
	} else {
		durationInt, err := strconv.Atoi(duration)
		if err != nil {
			return 0, fmt.Errorf("ошибка окружения: %w", ErrInvalidJWTDuration)
		}
		return durationInt, nil
	}
}

func CheckEnvAdminJWTDuration() (int, error) {
	duration := os.Getenv("JWT_ADMIN_DURATION")

	if duration == "" {
		return DefaultAdminJWTDuration, fmt.Errorf(
			"ошибка окружения: в JWT_ADMIN_DURATION %w, из-за чего использовано значение по умолчанию – %d",
			ErrNotFoundJWTDuration, DefaultUserJWTDuration)
	} else {
		durationInt, err := strconv.Atoi(duration)
		if err != nil {
			return 0, fmt.Errorf("ошибка окружения: %w", ErrInvalidJWTDuration)
		}
		return durationInt, nil
	}
}
