FROM golang:1.23-alpine AS builder

# Set Go environment
ENV GOTOOLCHAIN=go1.23.4

# Install protoc and plugins
RUN apk add --no-cache protobuf-dev
RUN go version
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0

WORKDIR /app
COPY server/ .

# Generate protobuf code
RUN protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/admin.proto

# Build the application
RUN go mod tidy
RUN go build -o admin-service main.go

# Final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates

# Create a non-root user for demonstration
RUN adduser -D -s /bin/sh normaluser

WORKDIR /root/

COPY --from=builder /app/admin-service .

# Create the socket directory
RUN mkdir -p /tmp

# Expose the socket path as a volume
VOLUME ["/tmp"]

CMD ["./admin-service"]
