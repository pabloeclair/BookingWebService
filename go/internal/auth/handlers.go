package auth

import (
	"cu_coworking_book/go/internal/pb"
	"sync"
)

type AuthServer struct {
	pb.UnimplementedAuthenticationServer
	mu sync.RWMutex
}
