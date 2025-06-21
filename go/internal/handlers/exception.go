package handlers

import (
	"cu_coworking_book/go/internal/models"
	"net/http"
)

// Сообщает о несуществуеющей адресной строке
func NotFoundException(w http.ResponseWriter, req *http.Request) {
	errDto := models.NewExceptionDto(
		http.StatusNotFound,
		"адресная строка не найдена",
	)
	errDto.WriteException(w)
}

// Сообщает о неверном методе запроса
func MethodNotAllowedException(w http.ResponseWriter, req *http.Request) {
	errDto := models.NewExceptionDto(
		http.StatusMethodNotAllowed,
		"недопустимый метод",
	)
	errDto.WriteException(w)
}
