package models

import (
	"errors"
	"log"
	"os"
	"strconv"

	"github.com/golang-jwt/jwt"
)

type UserClaim struct {
	Id         uint32 `json:"id"`
	Email      string `json:"email"`
	FirstName  string `json:"first_name"`
	SecondName string `json:"second_name"`
	Patronymic string `json:"patronymic"`
	Role       string `json:"role"`
	jwt.StandardClaims
}

// Проверяет наличие и корректность переменных окружения, указывающих продолжительность
// JWT-токенов и в случае успеха возвращает указанное число или число по умолчанию (60) при отсутствии
func CheckJEnvWTDuration(duration string, isAdmin bool) int {
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

func ParseJWT(tokenString string) (*UserClaim, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaim{}, func(t *jwt.Token) (interface{}, error) {
		secretKey := []byte(os.Getenv("JWT_SECRET_KEY"))
		if secretKey == nil {
			return nil, errors.New("отсутствует secret key")
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}
	return token.Claims.(*UserClaim), nil
}
