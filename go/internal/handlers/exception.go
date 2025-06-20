package handlers

import (
	"cu_coworking_book/go/internal/models"
	"net/http"
)

func NotFoundException(w http.ResponseWriter, req *http.Request) {
	errDto := models.NewExceptionDto(
		http.StatusNotFound,
		"адресный путь не найден",
	)
	errDto.WriteException(w)
}

func MethodNotAllowedException(w http.ResponseWriter, req *http.Request) {
	errDto := models.NewExceptionDto(
		http.StatusMethodNotAllowed,
		"недопустимый метод",
	)
	errDto.WriteException(w)
}
