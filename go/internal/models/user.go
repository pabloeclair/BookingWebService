package models

import (
	"context"
	"cu_coworking_book/go/internal/db"
	"cu_coworking_book/go/internal/utils"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type Role int32

const (
	Role_ADMIN      Role = 0
	Role_USER       Role = 1
	Role_MAIN_ADMIN Role = 2
)

var (
	Role_name = map[Role]string{
		0: "ADMIN",
		1: "USER",
		2: "MAIN_ADMIN",
	}
	Role_value = map[string]Role{
		"ADMIN":      0,
		"USER":       1,
		"MAIN_ADMIN": 2,
	}
)

func (r Role) String() string {
	if name, ok := Role_name[r]; ok {
		return name
	}
	return "UNKNOWN"
}

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Сверяет пароль пользователя
func (u *UserLoginRequest) ComparePassword(ctx context.Context) (*db.User, error) {
	user, err := db.GetUserByEmail(ctx, u.Email)
	if err != nil {
		return nil, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(u.Password)); err != nil {
		return nil, fmt.Errorf("%w: неверный пароль", ErrPermissionDenied)
	}
	return user, nil
}

func (u *UserLoginRequest) GenerateJWT(ctx context.Context) (string, error) {
	res, err := u.ComparePassword(ctx)
	if err != nil {
		return "", nil
	}

	secretKey := []byte(os.Getenv("JWT_SECRET_KEY"))
	if secretKey == nil {
		log.Fatal("значение переменной JWT_SECRET_KEY обязана быть указанной")
	}

	durationStr := os.Getenv("JWT_USER_DURATION")
	durationInt := utils.CheckJWTDuration(durationStr, false)

	userClaim := UserClaim{
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

type UserLoginResponse struct {
	Id         string `json:"id"`
	Email      string `json:"email"`
	FirstName  string `json:"first_name"`
	SecondName string `json:"second_name"`
	Patronymic string `json:"patronymic"`
	Password   string `json:"password"`
	Role       string `json:"role"`
}

type UserSignupRequest struct {
	Email      string `json:"email"`
	FirstName  string `json:"first_name"`
	SecondName string `json:"second_name"`
	Patronymic string `json:"patronymic"`
	Password   string `json:"password"`
}

func (u *UserSignupRequest) Validation() error {
	if u.Email == "" || len(strings.Split(u.Email, " ")) > 1 {
		return fmt.Errorf("%w: поле email должен быть не пустым и содержать лишь только адрес почты", ErrBadBody)
	}
	if u.FirstName == "" || len(strings.Split(u.FirstName, " ")) > 1 {
		return fmt.Errorf("%w: поле имени должно быть не пустым и содержать лишь только само имя", ErrBadBody)
	}
	if u.SecondName == "" || len(strings.Split(u.SecondName, " ")) > 1 {
		return fmt.Errorf("%w: поле фамилии должно быть не пустым и содержать лишь только саму фамилию", ErrBadBody)
	}
	if len(strings.Split(u.Patronymic, " ")) > 1 {
		return fmt.Errorf("%w: поле отчества должно содержать лишь только само отчество, если оно имеется", ErrBadBody)
	}
	return nil
}

func (userDto *UserSignupRequest) ToStructForDB(hashedPassword []byte) db.User {
	return db.User{
		Email:      userDto.Email,
		FirstName:  strings.ToLower(userDto.FirstName),
		SecondName: strings.ToLower(userDto.SecondName),
		Patronymic: strings.ToLower(userDto.Patronymic),
		Password:   string(hashedPassword),
		Role:       Role_USER.String(),
	}
}

type UserSignupResponse struct {
	Id uint32 `json:"id"`
}
