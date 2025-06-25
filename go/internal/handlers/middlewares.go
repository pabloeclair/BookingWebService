package handlers

import (
	"cu_coworking_book/go/internal/models"
	"cu_coworking_book/go/internal/utils"
	"log"
	"net/http"
	"slices"
	"strings"
)

// Мидлвейр для логирования результатов обработки запросов.
func LoggingMiddleware(handler http.Handler) http.Handler {
	// Данный мидлвейр фиксирует все статусы и сообщения после выполнения работы
	// хэндлеров и логирует их на терминал.
	//
	// Кроме того, он также проверяет каждый отправленный запрос на тип тела - если
	// тип не указан или он не типа json, то возвращает Bad Request Error или, в непредвиденных
	// случаях, Internal Server Error

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		// проверка типа тела запроса
		if (r.URL.Path != "/api/v1/user" || r.Method != http.MethodGet) && (r.Header.Get("Content-Type") == "" || r.Header.Get("Content-Type") != "application/json; charset=utf-8") {
			errDto := models.ExceptionDto{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: models.ErrBadContentType.Error(),
			}
			errDto.WriteException(w)
			log.Printf(
				"%s %s: %d - Bad Request: %s",
				r.Method,
				r.URL.Path,
				errDto.StatusCode,
				errDto.ErrorMessage,
			)
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

// Мидлвейр для проверки корректности jwt и наличия доступа пользователя.
func AuthMiddleware(handler http.Handler) http.Handler {
	// Данный мидлверй проверяет jwt токен каждого запроса, отправленный
	// по пути /api/v1/user или /api/v1/admin. Кроме того, во втором случае
	// проверяется также и роль пользователя.

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "/api/v1/user") {
			handler.ServeHTTP(w, r)
			return
		}

		// получение jwt токена
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			log.Printf("%s %s: %d - Unauthorized: отсутствует JWT-токен", r.Method, r.URL.Path, http.StatusForbidden)
			errDto := models.ExceptionDto{
				StatusCode:   http.StatusUnauthorized,
				ErrorMessage: "отсутствует JWT-токен",
			}
			errDto.WriteException(w)
			return
		}

		// получение данных о пользователе из jwt токена
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

		// проверка корректности роли
		if !slices.Contains(models.Role_array, claims.Role) {
			log.Printf("%s %s: %d - Forbidden: отказано в доступе", r.Method, r.URL.Path, http.StatusForbidden)
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

// Мидлвейр для возможности обращения к серверу со стороны фронтенда
func СorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Expose-Headers", "Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
