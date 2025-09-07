package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"

	pb "grpc-goat/labs/grpc-009-ssrf/server/proto"
)

// unaryInterceptor logs client connections
func unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	peer, ok := peer.FromContext(ctx)
	if ok {
		log.Printf("Request from %s to %s", peer.Addr, info.FullMethod)
	}
	return handler(ctx, req)
}

type imagePreviewServer struct {
	pb.UnimplementedImagePreviewServer
}

func newImagePreviewServer() *imagePreviewServer {
	return &imagePreviewServer{}
}

// startFlagServer starts a local HTTP server with the flag
func startFlagServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/flag" {
			fmt.Fprint(w, "GRPC_GOAT{ssrf_internal_service_access}")
		} else {
			fmt.Fprint(w, "Internal service - try /flag endpoint")
		}
	})

	log.Println("Starting internal flag server on localhost:8080")
	go func() {
		if err := http.ListenAndServe("127.0.0.1:8080", nil); err != nil {
			log.Printf("Flag server error: %v", err)
		}
	}()
}

func (s *imagePreviewServer) FetchImage(ctx context.Context, req *pb.FetchImageRequest) (*pb.FetchImageResponse, error) {
	log.Printf("Fetching image from URL: %s", req.Url)

	// VULNERABILITY: SSRF - directly fetching user-provided URLs without validation
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(req.Url)
	if err != nil {
		return &pb.FetchImageResponse{
			Success: false,
			Content: "Failed to fetch URL: " + err.Error(),
		}, nil
	}
	defer resp.Body.Close()

	// Read response content
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &pb.FetchImageResponse{
			Success: false,
			Content: "Failed to read response: " + err.Error(),
		}, nil
	}

	content := string(body)

	return &pb.FetchImageResponse{
		Success: true,
		Content: content,
	}, nil
}

func main() {
	// Start the internal flag server
	startFlagServer()
	time.Sleep(1 * time.Second) // Give flag server time to start

	lis, err := net.Listen("tcp", ":8009")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer(grpc.UnaryInterceptor(unaryInterceptor))
	pb.RegisterImagePreviewServer(s, newImagePreviewServer())

	log.Println("Image Preview gRPC server starting on port 8009...")
	log.Println("Internal flag server running on localhost:8080")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
