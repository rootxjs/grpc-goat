# Lab 003: Insecure TLS

## Vulnerability
Billing Service using self-signed certificates that can't be verified by clients.

## Objective
Process a payment to capture the flag (despite TLS warnings).

## Run the Lab
```bash
# Build and run
docker build -t grpc-003 .
docker run -p 8003:8003 grpc-003
```

## Exploit
```bash
# Install grpcurl
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

# Connect with insecure TLS (ignore certificate errors)
grpcurl -insecure -d '{
  "customer_id": "CUST001",
  "card_number": "4111111111111111",
  "expiry_date": "12/25",
  "cvv": "123",
  "amount": 99.99,
  "currency": "USD"
}' localhost:8003 billing.BillingService/ProcessPayment

# Alternative: Use openssl to see certificate details
openssl s_client -connect localhost:8003 -servername localhost

# The -insecure flag bypasses certificate validation
# Self-signed cert means no trusted CA validation
```

## Impact
Self-signed certificates allow man-in-the-middle attacks since clients can't verify server identity.

## Flag
Successfully process a payment to get: `GRPC_GOAT{insecure_tls_allows_mitm_attacks}`

## Mitigation
Use proper CA-signed certificates and strong TLS configuration.