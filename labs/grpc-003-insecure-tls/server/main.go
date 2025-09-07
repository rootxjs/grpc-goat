package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"log"
	"math/big"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"

	pb "grpc-goat/labs/grpc-003-insecure-tls/server/proto"
)

// unaryInterceptor logs client IP addresses for all requests
func unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	peer, ok := peer.FromContext(ctx)
	if ok {
		log.Printf("Request from %s to %s", peer.Addr, info.FullMethod)
	}
	return handler(ctx, req)
}

type billingServer struct {
	pb.UnimplementedBillingServiceServer
}

func newBillingServer() *billingServer {
	return &billingServer{}
}

func (s *billingServer) ProcessPayment(ctx context.Context, req *pb.PaymentRequest) (*pb.PaymentResponse, error) {
	log.Printf("Payment processing - Customer: %s, Card: %s, CVV: %s, Amount: %.2f %s",
		req.CustomerId, req.CardNumber, req.Cvv, req.Amount, req.Currency)

	// Simple validation
	if len(req.CardNumber) < 13 || len(req.Cvv) != 3 {
		return &pb.PaymentResponse{
			Success: false,
			Message: "Invalid card details",
		}, nil
	}

	// Generate transaction ID
	transactionID := "TXN_" + req.CustomerId + "_" + time.Now().Format("20060102150405")

	return &pb.PaymentResponse{
		Success:       true,
		Message:       "Payment processed successfully",
		TransactionId: transactionID,
		Flag:          "GRPC_GOAT{insecure_tls_allows_mitm_attacks}",
	}, nil
}

// generateSelfSignedCert creates a self-signed certificate
func generateSelfSignedCert() (tls.Certificate, error) {
	// Generate private key
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return tls.Certificate{}, err
	}

	// Create certificate template
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization:  []string{"Insecure Corp"},
			Country:       []string{"US"},
			Province:      []string{""},
			Locality:      []string{"San Francisco"},
			StreetAddress: []string{""},
			PostalCode:    []string{""},
		},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().Add(365 * 24 * time.Hour), // Valid for 1 year
		KeyUsage:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IPAddresses: []net.IP{net.IPv4(127, 0, 0, 1)},
		DNSNames:    []string{"localhost", "insecure-billing.local"},
	}

	// Create certificate
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return tls.Certificate{}, err
	}

	// Encode certificate and key to PEM format
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})

	// Create TLS certificate
	cert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		return tls.Certificate{}, err
	}

	return cert, nil
}

func main() {
	lis, err := net.Listen("tcp", ":8003")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// VULNERABILITY: Generate self-signed certificate
	cert, err := generateSelfSignedCert()
	if err != nil {
		log.Fatalf("Failed to generate certificate: %v", err)
	}

	// VULNERABILITY: Create TLS config with insecure settings
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		// INSECURE: Self-signed certificate (no CA validation)
		// INSECURE: Allow any client to connect without proper verification
		ClientAuth: tls.NoClientCert,
		// INSECURE: Allow older TLS versions
		MinVersion: tls.VersionTLS12,
		MaxVersion: tls.VersionTLS13,
	}

	creds := credentials.NewTLS(tlsConfig)
	s := grpc.NewServer(
		grpc.Creds(creds),
		grpc.UnaryInterceptor(unaryInterceptor),
	)

	pb.RegisterBillingServiceServer(s, newBillingServer())

	log.Println("Billing gRPC server starting on port 8003...")
	log.Println("WARNING: Using self-signed certificate and weak TLS configuration!")
	log.Println("Test payment: customer_id=CUST001, card_number=4111111111111111, cvv=123")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
