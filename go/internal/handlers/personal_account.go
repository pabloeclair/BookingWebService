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

func UpdatePassword(w http.ResponseWriter, req *http.Request) {

	if req.Method != http.MethodPost {
		errDto := models.NewExceptionDto(
			http.StatusMethodNotAllowed,
			"допустимым методом является только PUT",
		)
		errDto.WriteException(w)
		return
	}

	ctx := req.Context()

	var user models.UserUpdatePasswordRequest
	if err := models.JsonToStruct(&user, req.Body); err != nil {
		errDto := models.NewExceptionDto(
			http.StatusBadRequest,
			fmt.Sprintf("%s, допустимые поля: %s", models.ErrBadBody.Error(), "new_password, old_password"),
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
	claims, err := utils.ParseJWT(tokenString)
	if err != nil {
		errDto := models.NewExceptionDto(
			http.StatusUnauthorized,
			err.Error(),
		)
		errDto.WriteException(w)
		return
	}

	emailPassword := models.UserEmailPassword{Email: claims.Email, Password: user.OldPassword}
	if _, err := emailPassword.ComparePassword(ctx); err != nil {
		errDto := models.NewExceptionDto(
			http.StatusInternalServerError,
			err.Error(),
		)
		if errors.Is(err, models.ErrPermissionDenied) {
			errDto.StatusCode = http.StatusUnauthorized
		}
		if errors.Is(err, db.ErrNotFound) {
			errDto.StatusCode = http.StatusNotFound
		}
		errDto.WriteException(w)
		return
	}

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

func DeleteUser(w http.ResponseWriter, req *http.Request) {

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

	}

}
