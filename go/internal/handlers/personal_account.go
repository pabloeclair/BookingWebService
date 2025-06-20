package handlers

import (
	"cu_coworking_book/go/internal/db"
	"cu_coworking_book/go/internal/models"
	"errors"
	"fmt"
	"net/http"
)

func UpdateUser(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	var user models.UserUpdateRequest
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

	tokenString := req.Header.Get("Authorization")
	claims, err := models.ParseJWT(tokenString)
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
		FirstName:  claims.FirstName,
		SecondName: claims.SecondName,
		Patronymic: user.Patronymic,
	}
	if user.Email != "" {
		userDb.Email = user.Email
	}
	if user.FirstName != "" {
		userDb.FirstName = user.FirstName
	}
	if user.SecondName != "" {
		userDb.SecondName = user.SecondName
	}

	res := models.UserEmailPassword{
		Email:    userDb.Email,
		Password: user.Password,
	}
	if _, err := res.ComparePassword(ctx); err != nil {
		errDto := models.ExceptionDto{ErrorMessage: err.Error()}
		if errors.Is(err, db.ErrNotFound) {
			errDto.StatusCode = http.StatusNotFound
		} else {
			errDto.StatusCode = http.StatusUnauthorized
		}
		errDto.WriteException(w)
		return
	}

	if err := db.UpdateUser(ctx, userDb); err != nil {
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

	newToken, err := res.GenerateJWT(ctx)
	if err != nil {
		errDto := models.NewExceptionDto(
			http.StatusInternalServerError,
			err.Error(),
		)
		errDto.WriteException(w)
		return
	}
	w.Header().Set("Authorization", newToken)
	w.Header().Del("Content-Type")
	w.WriteHeader(http.StatusNoContent)
}
