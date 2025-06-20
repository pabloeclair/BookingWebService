package handlers

import (
	"cu_coworking_book/go/internal/models"
	"cu_coworking_book/go/internal/utils"
	"log"
	"net/http"
	"slices"
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
		w.Header().Set("Content-Type", "application/json")

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
		log.Printf(
			"%s %s: %d - %s",
			r.Method,
			r.URL.Path,
			lrw.StatusCode,
			lrw.StatusMessage,
		)
	})
}

func AuthMiddleware(handler http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "/api/v1/user") {
			handler.ServeHTTP(w, r)
			return
		}

		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			log.Printf("%s %s: %d - отсутствует JWT-токен", r.Method, r.URL.Path, http.StatusForbidden)
			errDto := models.ExceptionDto{
				StatusCode:   http.StatusUnauthorized,
				ErrorMessage: "отсутствует JWT-токен",
			}
			errDto.WriteException(w)
			return
		}

		claims, err := utils.ParseJWT(tokenString)
		if err != nil {
			log.Printf("%s %s: %d - %s", r.Method, r.URL.Path, http.StatusInternalServerError, err.Error())
			errDto := models.ExceptionDto{
				StatusCode:   http.StatusUnauthorized,
				ErrorMessage: err.Error(),
			}
			errDto.WriteException(w)
			return
		}

		roles := []string{
			models.Role_USER.String(),
			models.Role_ADMIN.String(),
			models.Role_MAIN_ADMIN.String(),
		}
		if !slices.Contains(roles, claims.Role) {
			log.Printf("%s %s: %d - отказано в доступе", r.Method, r.URL.Path, http.StatusForbidden)
			errDto := models.ExceptionDto{
				StatusCode:   http.StatusForbidden,
				ErrorMessage: "отказано в доступе",
			}
			errDto.WriteException(w)
			return
		}

		handler.ServeHTTP(w, r)
	})
}
