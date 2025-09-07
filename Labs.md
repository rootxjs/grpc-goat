| # | Service Name | Business Purpose | Vulnerability | Risk / Impact |
|---|--------------|------------------|---------------|---------------|
| 001 | Service Discovery | Internal API registry for developers | Reflection enabled | Attackers can enumerate all gRPC services and hidden methods, discovering sensitive endpoints like admin functions |
| 002 | Auth Service | Handles user login and session tokens | Plaintext gRPC | Credentials sent over plaintext can be intercepted and reused by attackers |
| 003 | Billing Service | Processes customer payments | Insecure TLS | Self-signed TLS allows MITM attacks and interception/manipulation of transactions |
| 004 | Partner API | Exposes partner integrations | Arbitrary mTLS | Accepts any client certificate, letting attackers impersonate trusted partners and access restricted APIs |
| 005 | Partner API v2 | Enhanced partner integrations | mTLS Subject Validation | Validates subject name but accepts self-signed certificates, allowing certificate impersonation |
| 006 | Admin Service | System administration functions | Unix Socket World Writable | Socket with world read/write permissions allows any user to access admin functions |
| 007 | User Directory | Stores employee profiles and permissions | SQL Injection | Unsanitized database queries allow attackers to exfiltrate sensitive data (users, credentials, API keys) |
| 008 | File Processor | Processes uploaded files for reports | Command Injection | Unsanitized input allows attackers to execute arbitrary system commands on the server |
| 009 | Image Preview | Fetches thumbnails from external URLs | SSRF | Attackers can make the server request internal resources, potentially accessing metadata or internal endpoints |