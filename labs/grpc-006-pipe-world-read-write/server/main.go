package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/user"

	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"

	pb "grpc-goat/labs/grpc-006-pipe-world-read-write/server/proto"
)

// unaryInterceptor logs client connections
func unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	peer, ok := peer.FromContext(ctx)
	if ok {
		log.Printf("Request from %s to %s", peer.Addr, info.FullMethod)
	}
	return handler(ctx, req)
}

type adminServer struct {
	pb.UnimplementedAdminServiceServer
}

func newAdminServer() *adminServer {
	return &adminServer{}
}

func (s *adminServer) GetSystemInfo(ctx context.Context, req *pb.SystemInfoRequest) (*pb.SystemInfoResponse, error) {
	log.Printf("System info requested")

	// Get hostname
	hostname, _ := os.Hostname()

	// Get current user
	currentUser, _ := user.Current()
	username := currentUser.Username

	// Get working directory
	workingDir, _ := os.Getwd()

	// Get some environment variables
	envVars := []string{
		"PATH=" + os.Getenv("PATH"),
		"HOME=" + os.Getenv("HOME"),
		"USER=" + os.Getenv("USER"),
		"SHELL=" + os.Getenv("SHELL"),
	}

	return &pb.SystemInfoResponse{
		Success:          true,
		Hostname:         hostname,
		Username:         username,
		WorkingDirectory: workingDir,
		EnvironmentVars:  envVars,
	}, nil
}

func (s *adminServer) ExecuteCommand(ctx context.Context, req *pb.CommandRequest) (*pb.CommandResponse, error) {
	log.Printf("Command execution requested: %s %v", req.Command, req.Args)

	staticOutput := "GRPC_GOAT{unix_socket_world_writable_privilege_escalation}\nroot\n/root\nuid=0(root) gid=0(root) groups=0(root)"

	return &pb.CommandResponse{
		Success:  true,
		Output:   staticOutput,
		Error:    "",
		ExitCode: 0,
	}, nil
}

func main() {
	socketPath := "/tmp/grpc-admin.sock"

	os.Remove(socketPath)

	lis, err := net.Listen("unix", socketPath)
	if err != nil {
		log.Fatalf("Failed to listen on Unix socket: %v", err)
	}

	err = os.Chmod(socketPath, 0666) // rw-rw-rw-
	if err != nil {
		log.Printf("Warning: Failed to set socket permissions: %v", err)
	}

	if info, err := os.Stat(socketPath); err == nil {
		log.Printf("Socket permissions: %s", info.Mode().Perm())
	}

	s := grpc.NewServer(grpc.UnaryInterceptor(unaryInterceptor))
	pb.RegisterAdminServiceServer(s, newAdminServer())

	log.Printf("Admin gRPC server starting on Unix socket: %s", socketPath)
	log.Printf("WARNING: Socket has world read/write permissions (0666)!")
	log.Printf("Any user on the system can connect and access admin functions")

	// Cleanup socket on exit
	defer func() {
		os.Remove(socketPath)
	}()

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
