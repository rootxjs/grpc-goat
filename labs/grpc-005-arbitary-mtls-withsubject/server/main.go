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

	pb "grpc-goat/labs/grpc-005-arbitary-mtls-withsubject/server/proto"
)

// unaryInterceptor logs client IP addresses and certificate info
func unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	peer, ok := peer.FromContext(ctx)
	if ok {
		log.Printf("Request from %s to %s", peer.Addr, info.FullMethod)

		// Log client certificate info if available
		if tlsInfo, ok := peer.AuthInfo.(credentials.TLSInfo); ok {
			if len(tlsInfo.State.PeerCertificates) > 0 {
				cert := tlsInfo.State.PeerCertificates[0]
				log.Printf("Client cert: Subject=%s, Issuer=%s", cert.Subject, cert.Issuer)
			}
		}
	}
	return handler(ctx, req)
}

type partnerServer struct {
	pb.UnimplementedPartnerAPIServer
}

func newPartnerServer() *partnerServer {
	return &partnerServer{}
}

func (s *partnerServer) GetPartnerData(ctx context.Context, req *pb.PartnerDataRequest) (*pb.PartnerDataResponse, error) {
	log.Printf("Partner data request - Partner: %s, Type: %s", req.PartnerId, req.DataType)

	// VULNERABILITY: Validate client certificate subject name
	peer, ok := peer.FromContext(ctx)
	if !ok {
		return &pb.PartnerDataResponse{
			Success: false,
			Message: "No peer information available",
		}, nil
	}

	tlsInfo, ok := peer.AuthInfo.(credentials.TLSInfo)
	if !ok || len(tlsInfo.State.PeerCertificates) == 0 {
		return &pb.PartnerDataResponse{
			Success: false,
			Message: "No client certificate provided",
		}, nil
	}

	cert := tlsInfo.State.PeerCertificates[0]
	expectedSubject := "goatpartner.local"

	// INSECURE: Only checking subject CN, not validating against trusted CA
	if cert.Subject.CommonName != expectedSubject {
		log.Printf("Certificate validation failed: expected CN=%s, got CN=%s", expectedSubject, cert.Subject.CommonName)
		return &pb.PartnerDataResponse{
			Success: false,
			Message: "Invalid client certificate subject expected CN=goatpartner.local",
		}, nil
	}

	log.Printf("Certificate validation passed: CN=%s", cert.Subject.CommonName)

	// Return sensitive partner information (this should be restricted!)
	partners := []*pb.PartnerInfo{
		{
			PartnerId:   "PARTNER_001",
			Name:        "TechCorp Solutions",
			ApiKey:      "tc_live_sk_1234567890abcdef",
			Secret:      "tc_secret_xyz789",
			Permissions: []string{"read", "write", "admin"},
		},
		{
			PartnerId:   "PARTNER_002",
			Name:        "DataFlow Inc",
			ApiKey:      "df_live_pk_fedcba0987654321",
			Secret:      "df_secret_abc123",
			Permissions: []string{"read", "analytics"},
		},
		{
			PartnerId:   "PARTNER_003",
			Name:        "SecureBank API",
			ApiKey:      "sb_live_key_999888777666",
			Secret:      "sb_secret_banking_2024",
			Permissions: []string{"financial", "transactions", "admin"},
		},
	}

	return &pb.PartnerDataResponse{
		Success:  true,
		Message:  "Partner data retrieved successfully",
		Partners: partners,
		Flag:     "GRPC_GOAT{subject_validation_insufficient_for_mtls}",
	}, nil
}

// generateSelfSignedCert creates a self-signed certificate
func generateSelfSignedCert() (tls.Certificate, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return tls.Certificate{}, err
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization:  []string{"Partner API Corp"},
			Country:       []string{"US"},
			Province:      []string{""},
			Locality:      []string{"San Francisco"},
			StreetAddress: []string{""},
			PostalCode:    []string{""},
		},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IPAddresses: []net.IP{net.IPv4(127, 0, 0, 1)},
		DNSNames:    []string{"localhost", "partner-api.local"},
	}

	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return tls.Certificate{}, err
	}

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})

	cert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		return tls.Certificate{}, err
	}

	return cert, nil
}

func main() {
	lis, err := net.Listen("tcp", ":8005")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Generate server certificate
	cert, err := generateSelfSignedCert()
	if err != nil {
		log.Fatalf("Failed to generate certificate: %v", err)
	}

	// VULNERABILITY: mTLS configuration that accepts ANY client certificate
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		// INSECURE: Require client certificates but don't verify them properly
		ClientAuth: tls.RequireAnyClientCert,
		// INSECURE: No client certificate verification
		InsecureSkipVerify: false,
		// INSECURE: Accept any client certificate without validation
	}

	creds := credentials.NewTLS(tlsConfig)
	s := grpc.NewServer(
		grpc.Creds(creds),
		grpc.UnaryInterceptor(unaryInterceptor),
	)

	pb.RegisterPartnerAPIServer(s, newPartnerServer())

	log.Println("Partner API gRPC server starting on port 8005...")
	log.Println("WARNING: mTLS validates subject name but accepts self-signed certificates!")
	log.Println("Attackers can generate certificates with subject 'goatpartner.local' to access partner data")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
