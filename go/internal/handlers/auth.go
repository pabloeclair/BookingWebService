package handlers

import (
	"cu_coworking_book/go/internal/db"
	"cu_coworking_book/go/internal/models"
	"cu_coworking_book/go/internal/utils"
	"errors"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// Регистрация нового пользователя.
func SignupUser(w http.ResponseWriter, req *http.Request) {
	// Первый пользователь автоматически становится MAIN_ADMIN, а остальные последующие — USER.
	//
	// - Если пользователь с указанной почтой уже существует, то вернется ошибка Conflict.
	//
	// При недопустимом поле запроса – Bad Request.
	// При любых других ошибках – Internal Server.
	// Показатель успеха — статус Created и id пользователя.

	if req.Method != http.MethodPost {
		errDto := models.NewExceptionDto(
			http.StatusMethodNotAllowed,
			"допустимым методом является только POST",
		)
		errDto.WriteException(w)
		return
	}

	ctx := req.Context()
	var user models.UserSignupRequest

	// получение тела запроса
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

	// шифрование пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		errDto := models.NewExceptionDto(
			http.StatusBadRequest,
			err.Error(),
		)
		errDto.WriteException(w)
		return
	}

	// сохранение пользователя
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

	// отправление ответа
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

// Авторизация пользователя и генерация JWT-токена.
func LoginUser(w http.ResponseWriter, req *http.Request) {
	// - Если пользователь с указанной почтой не найден, вернется ошибка NotFound.
	// - Если пароль пользователя не совпал — Unauthorized.
	//
	// При любых других ошибках – InternalServer или log.Fatal, если отсутствует JWT_SECRET_KEY
	// или указано не число в JWT_USER_DURATION.
	// Показатель успеха — статус Created и JWT-токен в заголовке Authorization.

	// проверка метода запроса
	if req.Method != http.MethodPost {
		errDto := models.NewExceptionDto(
			http.StatusMethodNotAllowed,
			"допустимым методом является только POST",
		)
		errDto.WriteException(w)
		return
	}

	ctx := req.Context()
	var user models.UserEmailPassword
	// получение тела запроса
	if err := models.JsonToStruct(&user, req.Body); err != nil {
		errDto := models.NewExceptionDto(
			http.StatusBadRequest,
			fmt.Sprintf("%s, допустимые поля: %s", models.ErrBadBody.Error(), "email, password"),
		)
		errDto.WriteException(w)
		return
	}

	// сравнение пароля
	if _, err := user.ComparePassword(ctx); err != nil {
		var errDto models.ExceptionDto
		if errors.Is(err, db.ErrNotFound) {
			errDto = models.NewExceptionDto(
				http.StatusNotFound,
				err.Error(),
			)
		} else {
			errDto = models.NewExceptionDto(
				http.StatusUnauthorized,
				fmt.Sprintf("%s: неверный пароль", models.ErrPermissionDenied.Error()),
			)
		}
		errDto.WriteException(w)
		return
	}

	// генерация jwt и отправление ответа
	token, err := utils.GenerateJWT(ctx, user.Email)
	if err != nil {
		errDto := models.NewExceptionDto(
			http.StatusInternalServerError,
			err.Error(),
		)
		errDto.WriteException(w)
		return
	}
	w.Header().Set("Authorization", token)
	w.Header().Del("Content-Type")
	w.WriteHeader(http.StatusCreated)
}
