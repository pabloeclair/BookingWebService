package handlers

import (
	"cu_coworking_book/go/internal/models"
	"log"
	"net/http"
	"strings"
)

func LoggingMiddleware(handler http.Handler) http.Handler {
	// Данный мидлвейр фиксирует все статусы и сообщения после выполнения работы
	// хэндлеров и логирует их на терминал.
	//
	// Кроме того, он также проверяет каждый отправленный запрос на тип тела - если
	// тип не указан или он не типа json, то возвращает Bad Request Error или, в непредвиденных
	// случаях, Internal Server Error

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if !strings.Contains(r.URL.Path, "login") {
			w.Header().Set("Content-Type", "application/json")
		}

		// возвращение ошибки Bad Request
		if r.Header.Get("Content-Type") == "" || r.Header.Get("Content-Type") != "application/json" {
			errDto := models.ExceptionDto{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: models.ErrBadContentType.Error(),
			}
			errDto.WriteException(w)
			return
		}

		// выполнение хэндлера и логирование
		lrw := models.NewLoggingResponseWriter(w)
		handler.ServeHTTP(lrw, r)
		if lrw.StatusMessage == "" {
			lrw.StatusMessage = "Success"
		}
		log.Printf("%s %s: %d - %s", r.Method, r.URL.Path, lrw.StatusCode, lrw.StatusMessage)
	})
}
