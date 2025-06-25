package handlers

import (
	"cu_coworking_book/go/internal/db"
	"cu_coworking_book/go/internal/models"
	"cu_coworking_book/go/internal/utils"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// Получение информации о пользователе по JWT
func GetUserByJWT(w http.ResponseWriter, req *http.Request) {
	// Если произошла любая ошибка с JWT – Unauthorized
	// Если пользователь из JWT не найден – Not Found
	// При любых других ошибках – InternalServer
	// Показатель успеха – удачно переданный claim

	// парсинг токена
	tokenString := req.Header.Get("Authorization")
	if tokenString == "" {
		errDto := models.NewExceptionDto(
			http.StatusUnauthorized,
			"необходим JWT токен для получения пользователя",
		)
		errDto.WriteException(w)
		return
	}

	token, err := utils.ParseJWT(tokenString)
	if err != nil {
		errDto := models.NewExceptionDto(
			http.StatusInternalServerError,
			err.Error(),
		)
		if !errors.Is(err, utils.ErrNotFoundSecretKey) {
			errDto.StatusCode = http.StatusUnauthorized
		}
		errDto.WriteException(w)
		return
	}

	tokenJson, err := models.StructToJson(&token)
	if err != nil {
		errDto := models.NewExceptionDto(
			http.StatusInternalServerError,
			err.Error(),
		)
		errDto.WriteException(w)
		return
	}

	// сравнение с реальными данными
	ctx := req.Context()
	actualUser, err := db.GetUserById(ctx, token.Id)
	if err != nil {
		errDto := models.NewExceptionDto(
			http.StatusInternalServerError,
			err.Error(),
		)
		if errors.Is(err, db.ErrNotFound) {
			errDto.StatusCode = http.StatusNotFound
		}
		errDto.WriteException(w)
		return
	}

	isCorrect := actualUser.Id != token.Id || actualUser.Email != token.Email || actualUser.FirstName != token.FirstName || actualUser.SecondName != token.SecondName || actualUser.Patronymic != token.Patronymic || actualUser.Role != token.Role
	if !isCorrect {
		errDto := models.NewExceptionDto(
			http.StatusUnauthorized,
			"данные не совпадают с данными реального пользователя",
		)
		errDto.WriteException(w)
		return
	}

	w.Write(tokenJson)
}

// Обновление пользователя, кроме роли и пароля.
func UpdateUser(w http.ResponseWriter, req *http.Request) {
	// Принимает новые значения по полям, которые необходимо обновить и
	// в случае успеха генерирует новый jwt токен.
	//
	// - Если отправлено некорректное тело запроса – BadRequest.
	// - Если отправлен некорректный jwt-токен – Unauthorized.
	// - Если отправлена уже существующая почта – Conflict.
	// - Если id пользователя из jwt не найден – NotFound.
	// Успех – NoContent.

	ctx := req.Context()
	var user models.UserUpdateRequest

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

	// получение jwt токена и старой информации о пользователе
	tokenString := req.Header.Get("Authorization")
	claims, err := utils.ParseJWT(tokenString)
	if err != nil {
		errDto := models.NewExceptionDto(
			http.StatusUnauthorized,
			err.Error(),
		)
		errDto.WriteException(w)
		return
	}

	userDb := db.User{
		Id:         claims.Id,
		Email:      claims.Email,
		FirstName:  strings.ToLower(claims.FirstName),
		SecondName: strings.ToLower(claims.SecondName),
		Patronymic: strings.ToLower(user.Patronymic),
	}
	if user.Email != "" {
		userDb.Email = strings.ToLower(user.Email)
	}
	if user.FirstName != "" {
		userDb.FirstName = strings.ToLower(user.FirstName)
	}
	if user.SecondName != "" {
		userDb.SecondName = strings.ToLower(user.SecondName)
	}

	// обновление пользователя в бд
	if err := db.UpdateUser(ctx, &userDb); err != nil {
		errDto := models.NewExceptionDto(
			http.StatusInternalServerError,
			err.Error(),
		)
		if errors.Is(err, db.ErrConflict) {
			errDto.StatusCode = http.StatusConflict
		}
		if errors.Is(err, db.ErrNotFound) {
			errDto.StatusCode = http.StatusNotFound
		}
		errDto.WriteException(w)
		return
	}

	// генерация jwt токена и отправление ответа
	newToken, err := utils.GenerateJWT(ctx, userDb.Email)
	if err != nil {
		errDto := models.NewExceptionDto(
			http.StatusInternalServerError,
			err.Error(),
		)
		if errors.Is(err, db.ErrNotFound) {
			errDto.StatusCode = http.StatusNotFound
		}
		errDto.WriteException(w)
		return
	}
	w.Header().Set("Authorization", newToken)
	w.Header().Del("Content-Type")
	w.WriteHeader(http.StatusNoContent)
}

// Обновление пароля пользователя.
func UpdatePassword(w http.ResponseWriter, req *http.Request) {
	// Принимает значение нового пароля и обновляет старый.
	//
	// - Если метод запроса отличается от PUT – MethodNotAllowed
	// - Если некорректное тело запроса – BadRequest
	// - Если неверный пароль или некорректный jwt токен – Unauthorized
	// - Если id пользователя не найден – NotFound
	// Успех – NoContent

	// проверка метода запроса
	if req.Method != http.MethodPut {
		errDto := models.NewExceptionDto(
			http.StatusMethodNotAllowed,
			"допустимым методом является только PUT",
		)
		errDto.WriteException(w)
		return
	}

	ctx := req.Context()
	var user models.UserUpdatePasswordRequest

	// проверка тела запроса
	if err := models.JsonToStruct(&user, req.Body); err != nil {
		errDto := models.NewExceptionDto(
			http.StatusBadRequest,
			fmt.Sprintf("%s, допустимые поля: %s", models.ErrBadBody.Error(), "new_password"),
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

	// проверка jwt токена
	tokenString := req.Header.Get("Authorization")
	claims, err := utils.ParseJWT(tokenString)
	if err != nil {
		errDto := models.NewExceptionDto(
			http.StatusUnauthorized,
			err.Error(),
		)
		errDto.WriteException(w)
		return
	}

	// сохранение нового пароля в бд и отправление ответа
	if err := db.UpdatePassword(ctx, claims.Id, user.NewPassword); err != nil {
		errDto := models.NewExceptionDto(
			http.StatusInternalServerError,
			err.Error(),
		)
		if errors.Is(err, db.ErrNotFound) {
			errDto.StatusCode = http.StatusNotFound
		}
		errDto.WriteException(w)
		return
	}
	w.Header().Del("Content-Type")
	w.WriteHeader(http.StatusNoContent)
}

// Удаление пользователя.
func DeleteUser(w http.ResponseWriter, req *http.Request) {
	// Позволяет удалить свой аккаунт, если пользователь не является главным
	// администратором — в противном случае ему необходимо передать свои права.
	// Сделано это для того, чтобы гарантированно сохранялся один MAIN_ADMIN.
	//
	// Если jwt токен истек – Unauthorized.
	// Если пользователь имеет роль MAIN_ADMIN – Forbidden.
	// Если id пользователя не существует – NotFound.
	// Успех – NoContent.

	tokenString := req.Header.Get("Authorization")
	claims, err := utils.ParseJWT(tokenString)
	if err != nil {
		errDto := models.NewExceptionDto(
			http.StatusUnauthorized,
			err.Error(),
		)
		errDto.WriteException(w)
		return
	}

	ctx := req.Context()
	if claims.Role == models.Role_MAIN_ADMIN.String() {
		errDto := models.ExceptionDto{
			StatusCode:   http.StatusForbidden,
			ErrorMessage: "MAIN_ADMIN должен сначала передать свои права другому пользователю",
		}
		errDto.WriteException(w)
		return
	}

	id := claims.Id
	if err := db.DeleteUser(ctx, id); err != nil {
		errDto := models.NewExceptionDto(
			http.StatusInternalServerError,
			err.Error(),
		)
		if errors.Is(err, db.ErrNotFound) {
			errDto.StatusCode = http.StatusNotFound
		}
		errDto.WriteException(w)
		return
	}
	w.Header().Del("Content-Type")
	w.WriteHeader(http.StatusNoContent)
}
