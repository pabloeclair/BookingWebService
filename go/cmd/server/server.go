package main

import (
	"context"
	"cu_coworking_book/go/internal/db"
	"cu_coworking_book/go/internal/handlers"
	"cu_coworking_book/go/internal/utils"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if len(os.Args) != 2 {
		log.Fatal("ошибка запуска: необходимо ввести адрес запуска")
	}
	addrs := os.Args[1]

	// проверка перменных окружения
	godotenv.Load()

	var errResult error
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		errResult = errors.Join(
			errResult,
			fmt.Errorf("ошибка запуска: ошибка окружения: %w",
				utils.ErrNotFoundSecretKey,
			))
	}

	if _, err := utils.CheckEnvUserJWTDuration(); err != nil {
		if errors.Is(err, utils.ErrNotFoundJWTDuration) {
			log.Printf("ПРЕДУПРЕЖДЕНИЕ: %v", err)
		} else {
			errResult = errors.Join(errResult, fmt.Errorf("ошибка запуска: %w", err))
		}
	}

	if _, err := utils.CheckEnvAdminJWTDuration(); err != nil {
		if errors.Is(err, utils.ErrNotFoundJWTDuration) {
			log.Printf("ПРЕДУПРЕЖДЕНИЕ: %v", err)
		} else {
			errResult = errors.Join(errResult, fmt.Errorf("ошибка запуска: %w", err))
		}
	}

	psUser := os.Getenv("POSTGRES_USER")
	if psUser == "" {
		errResult = errors.Join(
			errResult,
			fmt.Errorf("ошибка запуска: ошибка окружения: %s",
				"не найден POSTGRES_USER",
			))
	}

	psPsswd := os.Getenv("POSTGRES_PASSWORD")
	if psPsswd == "" {
		errResult = errors.Join(
			errResult,
			fmt.Errorf("ошибка запуска: ошибка окружения: %s",
				"не найден POSTGRES_PASSWORD",
			))
	}

	psDatabase := os.Getenv("POSTGRES_DB")
	if psDatabase == "" {
		errResult = errors.Join(
			errResult,
			fmt.Errorf("ошибка запуска: ошибка окружения: %s",
				"не найден POSTGRES_DB",
			))
	}

	if errResult != nil {
		log.Println("фатальные ошибки:")
		fmt.Println(errResult)
		return
	}

	db.DSN = fmt.Sprintf(
		"postgres://%s:%s@db:5432/%s?sslmode=disable",
		psUser, psPsswd, psDatabase,
	)

	// timeout поднятия бд
	<-time.After(time.Second * 3)

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
		Handler: handlers.СorsMiddleware(handlers.AuthMiddleware(handlers.LoggingMiddleware(mux))),
	}
	hasError := false

	go func() {
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("ошибка прослушивания: %s", err)
			hasError = true
		}
	}()

	<-time.After(time.Second * 1) // чтобы не выводил сообщение, 
								  // если в начале прослушивания произойдет ошибка
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
