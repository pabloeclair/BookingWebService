package main

import (
	"context"
	"cu_coworking_book/go/internal/auth"
	"cu_coworking_book/go/internal/pb"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Error: you must enter the server address to start the server.")
		return
	}

	address := os.Args[1]

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go func() {
		lis, err := net.Listen("tcp", address)
		if err != nil {
			fmt.Println("Error: the invalid server address.")
			return
		}

		s := grpc.NewServer()
		pb.RegisterAuthenticationServer(s, &auth.AuthServer{})
		log.Printf("The gRPC server of authentication starts on address %s.", lis.Addr().String())
		if err = s.Serve(lis); err != nil {
			log.Fatalf("gRPC error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("The gRPC server closed.")
}
