package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os/exec"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"

	pb "grpc-goat/labs/grpc-008-grpc-command-injection/server/proto"
)

// unaryInterceptor logs client connections
func unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	peer, ok := peer.FromContext(ctx)
	if ok {
		log.Printf("Request from %s to %s", peer.Addr, info.FullMethod)
	}
	return handler(ctx, req)
}

type fileProcessorServer struct {
	pb.UnimplementedFileProcessorServer
}

func newFileProcessorServer() *fileProcessorServer {
	return &fileProcessorServer{}
}

func (s *fileProcessorServer) ListFiles(ctx context.Context, req *pb.ListFilesRequest) (*pb.ListFilesResponse, error) {
	log.Printf("Listing files in directory: %s", req.Directory)

	// VULNERABILITY: Command injection - directly using user input
	command := fmt.Sprintf("ls -la %s", req.Directory)

	log.Printf("Executing command: %s", command)

	output, err := exec.Command("sh", "-c", command).Output()
	if err != nil {
		return &pb.ListFilesResponse{
			Success: false,
			Output:  fmt.Sprintf("Command execution failed: %v", err),
		}, nil
	}
	outputStr := strings.TrimSpace(string(output))
	flag := ""

	// Check if flag was read through command injection
	if strings.Contains(outputStr, "GRPC_GOAT{command_injection_file_listing}") {
		flag = "GRPC_GOAT{command_injection_file_listing}"
	}

	return &pb.ListFilesResponse{
		Success: true,
		Output:  outputStr,
		Flag:    flag,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":8008")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer(grpc.UnaryInterceptor(unaryInterceptor))
	pb.RegisterFileProcessorServer(s, newFileProcessorServer())

	log.Println("File Listing gRPC server starting on port 8008...")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
