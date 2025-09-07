package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"

	pb "grpc-goat/labs/grpc-002-plaintext-grpc/server/proto"
)

// unaryInterceptor logs client IP addresses for all requests
func unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// Extract client IP from context
	peer, ok := peer.FromContext(ctx)
	if ok {
		log.Printf("Request from %s to %s", peer.Addr, info.FullMethod)
	}

	// Call the actual handler
	return handler(ctx, req)
}

type authServer struct {
	pb.UnimplementedAuthServiceServer
}

func newAuthServer() *authServer {
	return &authServer{}
}

func (s *authServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	log.Printf("Login attempt - Username: %s, Password: %s", req.Username, req.Password)

	// Simple hardcoded credentials check
	validCredentials := map[string]string{
		"admin":     "admin123",
		"user":      "password123",
		"developer": "dev456",
	}

	if password, exists := validCredentials[req.Username]; exists && password == req.Password {
		return &pb.LoginResponse{
			Success:      true,
			Message:      "Login successful",
			SessionToken: "session_" + req.Username + "_12345",
			Flag:         "GRPC_GOAT{plaintext_credentials_intercepted}",
		}, nil
	}

	return &pb.LoginResponse{
		Success: false,
		Message: "Invalid credentials pass username:admin and password:admin123 to login",
	}, nil
}

func main() {
	// Listen on port 8002
	lis, err := net.Listen("tcp", ":8002")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer(grpc.UnaryInterceptor(unaryInterceptor))

	// Register our service
	pb.RegisterAuthServiceServer(s, newAuthServer())

	log.Println("Auth gRPC server starting on port 8002...")
	log.Println("WARNING: Server running in PLAINTEXT mode - all credentials are visible!")
	log.Println("Test users: admin/admin123, user/password123, developer/dev456")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
