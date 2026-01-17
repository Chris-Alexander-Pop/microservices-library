# Security & Compliance Roadmap

## Application Security
- [x] **Middleware**:
    - CSRF Protection (Double Submit Cookie)
    - CSP (Content Security Policy) header builder
    - HSTS / SSL enforcement
    - CORS (Advanced configuration)
- [x] **Input Validation**:
    - XSS Sanitization
    - SQL Injection prevention checks
    - Path Traversal protection
    - Command Injection detection
- [x] **Cryptography**:
    - Key Management System (KMS) integration (Envelope Encryption)
    - Data-at-rest encryption helpers (AES-GCM)
    - Secure password hashing (Argon2id, bcrypt)
    - **Post-Quantum Cryptography**: Hybrid X25519 + Kyber (ML-KEM)

## Identity & Access (pkg/auth extended)
- [x] **OIDC / OAuth2**:
    - Client implementation (existing adapters)
- [x] **MFA (Multi-Factor Auth)**:
    - TOTP (Google Authenticator) generation & verification
    - Recovery codes
- [x] **Authorization**:
    - RBAC (existing in pkg/api/rbac)

## Compliance & Auditing
- [x] **Audit Logging**: Structured logs for SIEM ingestion
- [x] **PII Redaction**: Automatic masking of sensitive data (Credit Cards, SSNs, emails, etc.)
