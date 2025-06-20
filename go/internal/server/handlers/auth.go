package handlers

import (
	"cu_coworking_book/go/internal/db"
	"cu_coworking_book/go/internal/server/models"
	"errors"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// Регистрация нового пользователя
func SignupUser(w http.ResponseWriter, req *http.Request) {
	// Первый пользователь автоматически становится MAIN_ADMIN, а остальные последующие — USER.

	// - Если пользователь с указанной почтой уже существует, то вернется ошибка Conflict.

	// При недопустимом поле запроса – Bad Request.
	// При любых других ошибках – Internal Server.
	// Показатель успеха — статус Created и id пользователя.

	ctx := req.Context()
	var user models.UserSignupRequest

	// некорректное тело запроса
	if err := models.JsonToStruct(&user, req.Body); err != nil {
		errDto := models.NewExceptionDto(
			http.StatusBadRequest,
			fmt.Sprintf("%s, допустимые поля: %s", models.ErrBadBody.Error(), "email, first_name, second_name, patronymic, password"),
		)
		errDto.WriteException(w)
		return
	}

	if err := user.Validation(); err != nil {
		errDto := models.NewExceptionDto(
			http.StatusBadRequest,
			err.Error(),
		)
		errDto.WriteException(w)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		errDto := models.NewExceptionDto(
			http.StatusBadRequest,
			err.Error(),
		)
		errDto.WriteException(w)
		return
	}

	// обращение к бд
	id, err := db.CreateUser(ctx, user.ToStructForDB(hashedPassword))
	if err != nil {
		errDto := models.NewExceptionDto(
			http.StatusInternalServerError,
			err.Error(),
		)
		if errors.Is(err, db.ErrConflict) {
			errDto.StatusCode = http.StatusConflict
		}
		errDto.WriteException(w)
		return
	}

	res, err := models.StructToJson(models.UserSignupResponse{Id: id})
	if err != nil {
		errDto := models.NewExceptionDto(
			http.StatusInternalServerError,
			err.Error(),
		)
		errDto.WriteException(w)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}

// Авторизация пользователя и генерация JWT-токена
func LoginUser(w http.ResponseWriter, req *http.Request) {
	// - Если пользователь с указанной почтой не найден, вернется ошибка NotFound.
	// - Если пароль пользователя не совпал — Forbidden.

	// При любых других ошибках – InternalServer или log.Fatal, если отсутствует JWT_SECRET_KEY
	// или не указано число в JWT_USER_DURATION.
	// Показатель успеха — статус Created и JWT-токен в заголовке Authorization.

	ctx := req.Context()
	var user models.UserLoginRequest
	if err := models.JsonToStruct(&user, req.Body); err != nil {
		errDto := models.NewExceptionDto(
			http.StatusBadRequest,
			fmt.Sprintf("%s, допустимые поля: %s", models.ErrBadBody.Error(), "email, password"),
		)
		errDto.WriteException(w)
		return
	}

	if _, err := user.ComparePassword(ctx); err != nil {
		errDto := models.NewExceptionDto(
			http.StatusForbidden,
			fmt.Sprintf("%s: неверный пароль", models.ErrPermissionDenied.Error()),
		)
		errDto.WriteException(w)
		return
	}

	token, err := user.GenerateJWT(ctx)
	if err != nil {
		errDto := models.NewExceptionDto(
			http.StatusInternalServerError,
			err.Error(),
		)
		errDto.WriteException(w)
		return
	}
	w.Header().Set("Authorization", token)
	w.WriteHeader(http.StatusCreated)
}

func NotFoundError(w http.ResponseWriter, req *http.Request) {
	errDto := models.NewExceptionDto(
		http.StatusNotFound,
		"Адресный путь не найден",
	)
	errDto.WriteException(w)
}
