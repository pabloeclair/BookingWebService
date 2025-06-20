package models

import (
	"errors"
	"net/http"
)

var (
	ErrBadContentType   error = errors.New("недопустимый тип тела запроса: принимается только application/json")
	ErrBadBody          error = errors.New("недопустимое тело запроса")
	ErrPermissionDenied error = errors.New("доступ запрещен")
)

type ExceptionDto struct {
	StatusCode   int    `json:"status_code"`
	ErrorMessage string `json:"error_message"`
}

func NewExceptionDto(statusCode int, errorMessage string) ExceptionDto {
	return ExceptionDto{
		StatusCode:   statusCode,
		ErrorMessage: errorMessage,
	}
}

func (errDto *ExceptionDto) WriteException(w http.ResponseWriter) {
	res, err := StructToJson(errDto)
	if err != nil {
		http.Error(w, "произошла серверная ошибка: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(errDto.StatusCode)
	w.Write(res)
}
