# gRPC-goat
gRPC Goat is a "Vulnerable by Design" lab created to provide an interactive, hands-on playground for learning and practicing gRPC security.

Each lab contains a **CTF-style flag** that you can capture by successfully exploiting the vulnerability!

## Quick Start

### Run All Labs
```bash
# Start all vulnerable services (use Docker Compose V2)
docker compose up --build

# If you get permission errors, try:
sudo docker compose up --build

# Or add your user to docker group:
sudo usermod -aG docker $USER
# Then logout and login again

# Services will be available on:
# - Lab 001 (Reflection): localhost:8001
# - Lab 002 (Plaintext): localhost:8002
# - Lab 003 (Insecure TLS): localhost:8003
# - Lab 004 (Arbitrary mTLS): localhost:8004
# - Lab 005 (mTLS Subject): localhost:8005
# - Lab 006 (Unix Socket): grpc-006 container
# - Lab 007 (SQL Injection): localhost:8007
# - Lab 008 (Command Injection): localhost:8008
# - Lab 009 (SSRF): localhost:8009
```

### Run Individual Labs
```bash
# Lab 001: gRPC Reflection
cd labs/grpc-001-reflection-enabled
docker build -t grpc-001 .
docker run -p 8001:8001 grpc-001

# Lab 002: Plaintext gRPC
cd labs/grpc-002-plaintext-grpc
docker build -t grpc-002 .
docker run -p 8002:8002 grpc-002

# Lab 003: Insecure TLS
cd labs/grpc-003-insecure-tls
docker build -t grpc-003 .
docker run -p 8003:8003 grpc-003

# Lab 004: Arbitrary mTLS
cd labs/grpc-004-arbitary-mtls
docker build -t grpc-004 .
docker run -p 8004:8004 grpc-004

# Lab 005: mTLS with Subject Validation
cd labs/grpc-005-arbitary-mtls-withsubject
docker build -t grpc-005 .
docker run -p 8005:8005 grpc-005

# Lab 006: Unix Socket World Writable
cd labs/grpc-006-pipe-world-read-write
docker build -t grpc-006 .
docker run -it --rm grpc-006 sh

# Lab 007: SQL Injection
cd labs/grpc-007-sql-injection
docker build -t grpc-007 .
docker run -p 8007:8007 grpc-007
```

## Labs Overview
See [Labs.md](Labs.md) for detailed vulnerability descriptions.

## Prerequisites
- Docker and Docker Compose
- grpcurl for testing: `go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest`
