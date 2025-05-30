package admin

import (
	"context"
	"cu_coworking_book/go/internal/pb"
	"log"
	"sync"

	"google.golang.org/grpc"
)

type AdminService struct {
	pb.UnimplementedAuthenticationServer
	mu sync.RWMutex
}

func LogInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Printf("Server call: %s", info.FullMethod)

	resp, err := handler(ctx, req)

	if err != nil {
		log.Printf("RPC failed with error: %v", err)
	}

	return resp, err
}
