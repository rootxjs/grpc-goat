package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	_ "modernc.org/sqlite"

	pb "grpc-goat/labs/grpc-007-sql-injection/server/proto"
)

// unaryInterceptor logs client connections
func unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	peer, ok := peer.FromContext(ctx)
	if ok {
		log.Printf("Request from %s to %s", peer.Addr, info.FullMethod)
	}
	return handler(ctx, req)
}

type userDirectoryServer struct {
	pb.UnimplementedUserDirectoryServer
	db *sql.DB
}

func newUserDirectoryServer() *userDirectoryServer {
	// First create and populate the database
	setupDB, err := sql.Open("sqlite", "/tmp/users.db")
	if err != nil {
		log.Fatalf("Failed to create database: %v", err)
	}

	server := &userDirectoryServer{db: setupDB}
	server.initDatabase()
	setupDB.Close()

	// Now open in read-only mode
	db, err := sql.Open("sqlite", "file:/tmp/users.db?mode=ro")
	if err != nil {
		log.Fatalf("Failed to open read-only database: %v", err)
	}

	server.db = db
	return server
}

func (s *userDirectoryServer) initDatabase() {
	// Create users table
	createTable := `
	CREATE TABLE users (
		username TEXT PRIMARY KEY,
		email TEXT NOT NULL,
		role TEXT NOT NULL
	);`

	_, err := s.db.Exec(createTable)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	// Insert sample data
	users := [][]string{
		{"john", "john@company.com", "user"},
		{"admin", "admin@company.com", "admin"},
		{"flag_user", "flag@company.com", "GRPC_GOAT{sql_injection_data_exfiltration}"},
	}

	for _, user := range users {
		insertSQL := `INSERT INTO users (username, email, role) VALUES (?, ?, ?)`
		_, err := s.db.Exec(insertSQL, user[0], user[1], user[2])
		if err != nil {
			log.Printf("Failed to insert user %s: %v", user[0], err)
		}
	}

	log.Println("Database initialized with sample data")
}

func (s *userDirectoryServer) SearchUsers(ctx context.Context, req *pb.SearchUsersRequest) (*pb.SearchUsersResponse, error) {
	log.Printf("Search users - Username: %s", req.Username)

	// VULNERABILITY: SQL injection - directly concatenating user input
	query := fmt.Sprintf("SELECT username, email, role FROM users WHERE username = '%s'", req.Username)

	log.Printf("Executing SQL: %s", query)

	rows, err := s.db.Query(query)
	if err != nil {
		log.Printf("SQL error: %v", err)
		return &pb.SearchUsersResponse{
			Success: false,
		}, nil
	}
	defer rows.Close()

	var users []*pb.UserInfo
	flag := ""

	for rows.Next() {
		var user pb.UserInfo
		err := rows.Scan(&user.Username, &user.Email, &user.Role)
		if err != nil {
			log.Printf("Scan error: %v", err)
			continue
		}
		users = append(users, &user)

		// Check if we found the flag user or flag in role
		if user.Username == "flag_user" || strings.Contains(user.Role, "GRPC_GOAT") {
			flag = "GRPC_GOAT{sql_injection_data_exfiltration}"
		}
	}

	return &pb.SearchUsersResponse{
		Success: true,
		Users:   users,
		Flag:    flag,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":8007")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer(grpc.UnaryInterceptor(unaryInterceptor))
	pb.RegisterUserDirectoryServer(s, newUserDirectoryServer())

	log.Println("User Directory gRPC server starting on port 8007...")
	log.Println("WARNING: SQL injection vulnerability present!")
	log.Println("Database is read-only - safe for educational purposes")
	log.Println("Users: john, admin, flag_user")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
