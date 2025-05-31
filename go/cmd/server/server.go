package main

import (
	"context"
	"cu_coworking_book/go/internal/admin"
	"cu_coworking_book/go/internal/auth"
	"cu_coworking_book/go/internal/db"
	"cu_coworking_book/go/internal/pb"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Error: необходимо ввести адрес сервера и его тип")
		return
	}

	if os.Args[2] != "admin" && os.Args[2] != "auth" {
		fmt.Println("Error: адрес может быть только двух типов - admin и auth")
		return
	}

	address := os.Args[1]
	type_ := os.Args[2]

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	lis, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("Error: некорректный адрес сервера")
		return
	}

	var s *grpc.Server
	if type_ == "auth" {
		s = grpc.NewServer(grpc.UnaryInterceptor(auth.LogInterceptor))
		pb.RegisterAuthenticationServer(s, &auth.AuthServer{})
	} else {
		s = grpc.NewServer(grpc.UnaryInterceptor(auth.LogInterceptor))
		pb.RegisterAdminServiceServer(s, &admin.AdminService{})
	}
	log.Printf("GRPC %s сервер запущен по адресу %s", type_, lis.Addr().String())

	go func() {
		if err = s.Serve(lis); err != nil {
			log.Fatalf("gRPC error: %v", err)
		}
	}()

	go func() {
		<-time.After(time.Second * 7)
		if err := db.CreateTable(); err != nil {
			if errors.Is(err, db.ErrConDB) {
				log.Fatal(err)
			} else {
				log.Println(err)
			}
		}
	}()

	<-ctx.Done()
	s.GracefulStop()
	log.Println("GRPC сервер закрыт")
}
