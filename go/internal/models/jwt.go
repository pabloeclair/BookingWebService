package models

import (
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
