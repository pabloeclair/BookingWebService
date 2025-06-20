package utils

import (
	"context"
	"cu_coworking_book/go/internal/db"
	"cu_coworking_book/go/internal/models"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
)

// Генерирует JWT-токен указанному пользователю
func GenerateJWT(ctx context.Context, email string) (string, error) {
	secretKey := []byte(os.Getenv("JWT_SECRET_KEY"))
	if secretKey == nil {
		log.Fatal("значение переменной JWT_SECRET_KEY обязана быть указанной")
	}

	res, err := db.GetUserByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	durationStr := os.Getenv("JWT_USER_DURATION")
	durationInt := CheckEnvJWTDuration(durationStr, false)

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
			return nil, errors.New("отсутствует secret key")
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
func CheckEnvJWTDuration(duration string, isAdmin bool) int {
	var user string
	if isAdmin {
		user = "ADMIN"
	} else {
		user = "USER"
	}

	if duration == "" {
		log.Println("ПРЕДУПРЕЖДЕНИЕ: в переменной окружения JWT_" + user + "_DURATION, которая указывает на " +
			"продолжительность в минутах сохранения JWT-токенов, было пропущено значение, в следствии чего будет " +
			"использоваться значение по умолчанию - 60 минут")
		return 60
	} else {
		durationInt, err := strconv.Atoi(duration)
		if err != nil {
			log.Fatalf("ошибка окружения: необходимо указывать значение JWT_" + user + "_DURATION в формате int (в минутах)")
		}
		return durationInt
	}
}
