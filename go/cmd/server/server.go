package main

import (
	"context"
	"cu_coworking_book/go/internal/db"
	"cu_coworking_book/go/internal/handlers"
	"cu_coworking_book/go/internal/utils"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {

	if len(os.Args) != 2 {
		log.Fatal("ошибка запуска: необходимо ввести адрес запуска")
	}
	addrs := os.Args[1]

	// проверка перменных окружения
	godotenv.Load()
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		log.Fatal("ошибка запуска: необходимо указать значение окружения JWT_SECRET_KEY для генерации JWT-токенов")
	}
	durationUser := os.Getenv("JWT_USER_DURATION")
	utils.CheckEnvJWTDuration(durationUser, false)
	durationAdmin := os.Getenv("JWT_ADMIN_DURATION")
	utils.CheckEnvJWTDuration(durationAdmin, true)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	mux := http.NewServeMux()
	// auth
	mux.HandleFunc("/api/v1/signup", handlers.SignupUser)
	mux.HandleFunc("/api/v1/login", handlers.LoginUser)

	// personal account
	mux.HandleFunc("GET /api/v1/user", handlers.GetUserByJWT)
	mux.HandleFunc("PUT /api/v1/user", handlers.UpdateUser)
	mux.HandleFunc("DELETE /api/v1/user", handlers.DeleteUser)
	mux.HandleFunc("/api/v1/user/password", handlers.UpdatePassword)

	// admin
	// todo: admin handlers

	// exseptions
	mux.HandleFunc("/api/v1/user", handlers.MethodNotAllowedException)
	mux.HandleFunc("/", handlers.NotFoundException)

	s := http.Server{
		Addr:    addrs,
		Handler: handlers.AuthMiddleware(handlers.LoggingMiddleware(mux)),
	}
	hasError := false

	go func() {
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("ошибка прослушивания: %s", err)
			hasError = true
		}
	}()

	<-time.After(time.Second * 1) // чтобы не выводил сообщение, если в начале прослушивания произойдет ошибка
	if hasError {
		return
	}

	// создание таблицы, если не существует
	if err := db.CreateTable(); err != nil {
		if err := s.Shutdown(context.Background()); err != nil {
			log.Printf("ошибка закрытия сервера: %s", err)
		}
		log.Fatalf("ошибка подключения к бд: %s", err.Error())
	}

	log.Println("сервер запущен")

	// graseful shutdown
	<-ctx.Done()
	log.Println("сервер закрывается...")
	<-time.After(time.Second * 3)
	if err := s.Shutdown(context.Background()); err != nil {
		log.Printf("ошибка закрытия сервера: %s", err)
	}
	log.Println("сервер закрыт")
}
