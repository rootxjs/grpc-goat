# gRPC Goat

gRPC Goat is a "Vulnerable by Design" lab created to provide an interactive, hands-on playground for learning and practicing gRPC security.

Each lab contains a **CTF-style flag** that you can capture by successfully exploiting the vulnerability!

## Quick Start

```bash
# Clone the repository
git clone https://github.com/rootxjs/grpc-goat.git
cd grpc-goat

# Start all vulnerable services
docker compose up --build

# Services will be available on localhost:8001-8009
```

## Documentation

For complete documentation, installation guides, and step-by-step walkthroughs, visit:

**📖 [https://rootxjs.github.io/docs/grpc_goat_docs/getting-started/](https://rootxjs.github.io/docs/grpc_goat_docs/getting-started/)**

The documentation includes:
- **gRPC Basics** - Essential concepts and security fundamentals
- **Labs Overview** - All 9 vulnerability scenarios with learning paths
- **Installation Guide** - Detailed setup instructions and troubleshooting
- **Walkthrough** - Step-by-step exploitation guides with code examples

## Labs Overview

| Lab | Vulnerability | Port |
|-----|---------------|------|
| 001 | gRPC Reflection Enabled | 8001 |
| 002 | Plaintext gRPC | 8002 |
| 003 | Insecure TLS | 8003 |
| 004 | Arbitrary mTLS | 8004 |
| 005 | mTLS Subject Validation | 8005 |
| 006 | Unix Socket World Writable | container |
| 007 | SQL Injection | 8007 |
| 008 | Command Injection | 8008 |
| 009 | Server-Side Request Forgery | 8009 |

## Prerequisites

- Docker and Docker Compose
- grpcurl: `go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest`

## Contributing

Contributions are welcome! Please see the documentation website for contribution guidelines.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
