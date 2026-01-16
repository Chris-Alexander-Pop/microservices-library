# Security & Compliance Roadmap

## Application Security
- [ ] **Middleware**:
    - CSRF Protection (Double Submit Cookie)
    - CSP (Content Security Policy) header builder
    - HSTS / SSL enforcement
    - CORS (Advanced configuration)
- [ ] **Input Validation**:
    - XSS Sanitization (Bluemonday)
    - SQL Injection prevention checks
- [ ] **Cryptography**:
    - Key Management System (KMS) integration (AWS/GCP/Vault)
    - Data-at-rest encryption helpers (Envelope Encryption)
    - PGP / GPG signing helpers

## Identity & Access (pkg/auth extended)
- [ ] **OIDC / OAuth2**:
    - Client implementation (Log in with Google/Apple)
    - Server implementation (Ory Hydra / Fosite integration)
- [ ] **MFA (Multi-Factor Auth)**:
    - TOTP (Google Authenticator) generation & verification
    - WebAuthn / Passkeys support
- [ ] **Authorization**:
    - **OPA (Open Policy Agent)**: Rego policy evaluation
    - **Casbin**: RBAC/ABAC model implementation
    - **Zanzibar**: Google-style relationship-based access control (SpiceDB)

## Compliance & Auditing
- [ ] **Audit Logging**: Structure logs for SIEM ingestion (Splunk/Datadog)
- [ ] **PII Redaction**: Middleware to automatically mask Sensitive Data (Credit Cards, SSNs) in logs
- [ ] **SBOM**: Software Bill of Materials generation
