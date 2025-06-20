package handlers

import (
	"cu_coworking_book/go/internal/models"
	"errors"
	"log"
	"net/http"
	"os"
	"slices"
	"strings"

	"github.com/golang-jwt/jwt"
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
				StatusCode:   http.StatusForbidden,
				ErrorMessage: "отсутствует JWT-токен",
			}
			errDto.WriteException(w)
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &models.UserClaim{}, func(t *jwt.Token) (interface{}, error) {
			secretKey := []byte(os.Getenv("JWT_SECRET_KEY"))
			if secretKey == nil {
				return nil, errors.New("отсутствует secret key")
			}
			return secretKey, nil
		})

		if err != nil {
			log.Printf("%s %s: %d - %s", r.Method, r.URL.Path, http.StatusInternalServerError, err.Error())
			errDto := models.ExceptionDto{
				StatusCode:   http.StatusInternalServerError,
				ErrorMessage: err.Error(),
			}
			errDto.WriteException(w)
			return
		}

		claims := token.Claims.(*models.UserClaim)
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
