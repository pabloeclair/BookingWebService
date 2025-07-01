package models

import (
	"context"
	"cu_coworking_book/go/internal/db"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// Модель с почтой и паролем пользователя.
type UserEmailPassword struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Сверяет пароль пользователя. Может вернуть ErrNotFound и ErrPermissionDenied.
func (u *UserEmailPassword) ComparePassword(ctx context.Context) (*db.User, error) {

	// получение актуального пароля пользователя
	user, err := db.GetUserByEmail(ctx, u.Email)
	if err != nil {
		return nil, err
	}

	// сравнение паролей
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(u.Password)); err != nil {
		return nil, fmt.Errorf("%w: неверный пароль", ErrPermissionDenied)
	}
	return user, nil
}

// Модель запроса на регистрацию.
type UserSignupRequest struct {

	// Почта
	Email string `json:"email"`

	// Имя
	FirstName string `json:"first_name"`

	// Фамилия
	SecondName string `json:"second_name"`

	// Отчество (необязательное поле)
	Patronymic string `json:"patronymic"`

	// Пароль
	Password string `json:"password"`
}

// Проверка, что все поля заполнены корректно.
func (u *UserSignupRequest) Validation() error {
	if u.Email == "" || len(strings.Split(u.Email, " ")) > 1 {
		return fmt.Errorf("%w: поле email должен быть не пустым и не содержать пробелов", ErrBadBody)
	}
	if u.FirstName == "" || len(strings.Split(u.FirstName, " ")) > 1 {
		return fmt.Errorf("%w: поле имени должно быть не пустым и содержать лишь только само имя без пробелов", ErrBadBody)
	}
	if u.SecondName == "" || len(strings.Split(u.SecondName, " ")) > 1 {
		return fmt.Errorf("%w: поле фамилии должно быть не пустым и содержать лишь только саму фамилию без пробелов", ErrBadBody)
	}
	if len(strings.Split(u.Patronymic, " ")) > 1 {
		return fmt.Errorf("%w: поле отчества должно содержать лишь только само отчество без пробелов, если оно имеется", ErrBadBody)
	}
	if u.Password == "" {
		return fmt.Errorf("%w: поле пароля должно быть не пустым", ErrBadBody)
	}
	return nil
}

// Преобразование модели запроса в модель для сохранения в базе данных.
func (userDto *UserSignupRequest) ToStructForDB(hashedPassword []byte) *db.User {
	return &db.User{
		Email:      userDto.Email,
		FirstName:  strings.ToLower(userDto.FirstName),
		SecondName: strings.ToLower(userDto.SecondName),
		Patronymic: strings.ToLower(userDto.Patronymic),
		Password:   string(hashedPassword),
		Role:       Role_USER.String(),
	}
}

// Модель успешного ответа на запрос регистрации.
type UserSignupResponse struct {
	Id uint32 `json:"id"`
}

// Модель запроса на обновление пользователя.
type UserUpdateRequest struct {

	// Почта (необязательно поле)
	Email string `json:"email"`

	// Имя (необязательное поле)
	FirstName string `json:"first_name"`

	// Фамилия (необязательное поле)
	SecondName string `json:"second_name"`

	// Отчество (необязательное поле)
	Patronymic string `json:"patronymic"`
}

// Проверка, что все поля заполнены корректно.
func (u *UserUpdateRequest) Validation() error {
	if len(strings.Split(u.Email, " ")) > 1 {
		return fmt.Errorf("%w: поле email должен содержать лишь только адрес почты", ErrBadBody)
	}
	if len(strings.Split(u.FirstName, " ")) > 1 {
		return fmt.Errorf("%w: поле имени должно содержать лишь только само имя", ErrBadBody)
	}
	if len(strings.Split(u.SecondName, " ")) > 1 {
		return fmt.Errorf("%w: поле фамилии должно содержать лишь только саму фамилию", ErrBadBody)
	}
	if len(strings.Split(u.Patronymic, " ")) > 1 {
		return fmt.Errorf("%w: поле отчества должно содержать лишь только само отчество", ErrBadBody)
	}
	return nil
}

// Модель обновления пароля пользователя
type UserUpdatePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func (u *UserUpdatePasswordRequest) Validation() error {
	if u.OldPassword == "" {
		return fmt.Errorf("%w: поле старого пароля обязательно", ErrBadBody)
	}
	if u.NewPassword == "" {
		return fmt.Errorf("%w: поле нового пароля обязательно", ErrBadBody)
	}
	return nil
}
