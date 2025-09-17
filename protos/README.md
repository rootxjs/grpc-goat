# gRPC Goat Proto Files

This directory contains all the Protocol Buffer (.proto) files for the gRPC Goat labs. These files are needed to interact with the gRPC services using tools like grpcurl, Postman, or custom clients.

## Usage

### With grpcurl
```bash
# Example: Lab 002 - Auth Service
grpcurl -plaintext -proto protos/lab-002-auth.proto \
  -d '{"username": "admin", "password": "password"}' \
  localhost:8002 auth.AuthService/Login
```

### With Postman
1. Import the proto file in Postman's gRPC request
2. Select the service and method
3. Fill in the request data
4. Send the request

### With Custom Clients
Use these proto files to generate client code in your preferred language:
```bash
# Generate Go client
protoc --go_out=. --go-grpc_out=. protos/lab-002-auth.proto

# Generate Python client
python -m grpc_tools.protoc -I protos --python_out=. --grpc_python_out=. protos/lab-002-auth.proto
```

## Lab Proto Files

| Lab | Service | Proto File | Description |
|-----|---------|------------|-------------|
| 001 | Service Discovery | *Uses reflection* | No proto file needed |
| 002 | Auth Service | `lab-002-auth.proto` | User authentication |
| 003 | Billing Service | `lab-003-billing.proto` | Payment processing |
| 004 | Partner API | `lab-004-partner.proto` | Partner integrations |
| 005 | Partner API v2 | `lab-005-partner-v2.proto` | Enhanced partner API |
| 006 | Admin Service | `lab-006-admin.proto` | System administration |
| 007 | User Directory | `lab-007-user-directory.proto` | Employee profiles |
| 008 | File Processor | `lab-008-file-processor.proto` | File processing |
| 009 | Image Preview | `lab-009-image-preview.proto` | Image fetching |

## Notes

- Lab 001 uses gRPC reflection, so no proto file is needed
- All other labs require the corresponding proto file for client interaction
- Proto files are copied from each lab's `server/proto/` directory
- These files are kept in sync with the actual service implementations
