package models

import (
	"errors"
	"net/http"
)

var (
	ErrBadContentType   error = errors.New("недопустимый тип тела запроса: принимается только 'application/json; charset=utf-8'")
	ErrBadBody          error = errors.New("недопустимое поле")
	ErrPermissionDenied error = errors.New("доступ запрещен")
)

// Модель ошибки
type ExceptionDto struct {
	// HTTP код ошибки
	StatusCode int `json:"status_code"`

	// Отображаемое сообщение об ошибке
	ErrorMessage string `json:"error_message"`
}

// Конструктор ExceptionDto
func NewExceptionDto(statusCode int, errorMessage string) ExceptionDto {
	return ExceptionDto{
		StatusCode:   statusCode,
		ErrorMessage: errorMessage,
	}
}

// Возвращение ошибки в виде ответа на запрос и установка HTTP статуса
func (errDto *ExceptionDto) WriteException(w http.ResponseWriter) {
	res, err := StructToJson(errDto)
	if err != nil {
		http.Error(w, "произошла серверная ошибка: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(errDto.StatusCode)
	w.Write(res)
}
