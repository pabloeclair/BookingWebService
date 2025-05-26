package main

import (
	"context"
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
	if len(os.Args) != 2 {
		fmt.Println("Error: you must enter the server address to start the server")
		return
	}

	address := os.Args[1]

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go func() {
		lis, err := net.Listen("tcp", address)
		if err != nil {
			fmt.Println("Error: the invalid server address")
			return
		}

		s := grpc.NewServer(grpc.UnaryInterceptor(auth.LogInterceptor))
		pb.RegisterAuthenticationServer(s, &auth.AuthServer{})
		log.Printf("The gRPC server of authentication starts on address %s.", lis.Addr().String())
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
	log.Println("The gRPC server closed.")
}
