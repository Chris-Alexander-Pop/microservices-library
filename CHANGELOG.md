# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- `pkg/README.md` navigation index for package discovery
- `pkg/compute/doc.go` package documentation
- `pkg/data/doc.go` package documentation
- Cross-reference documentation in ratelimit and AI packages

### Changed
- Reorganized package structure for better discoverability:
  - `pkg/crypto` → `pkg/security/crypto`
  - `pkg/iam` → `pkg/security/iam`
  - `pkg/bigdata` → `pkg/data/bigdata`
  - `pkg/client` → `pkg/api/client`
  - `pkg/secrets` merged into `pkg/security/secrets`
- Reduced top-level packages from 39 to 34

### Removed
- `pkg/secrets` (consolidated into `pkg/security/secrets`)

## [1.0.0] - 2026-01-27

### Added
- Initial release with 34 core packages
- 130+ reference microservices
- Comprehensive package standards documentation
- CI/CD with GitHub Actions
- Full test coverage for core packages

### Packages Included
- **Core**: errors, logger, config, validator, concurrency, resilience, events, telemetry, test
- **Data**: cache, database, storage, data, streaming
- **Communication**: messaging, communication, api
- **Security**: auth, security (iam, crypto, secrets, fraud, captcha, waf, scanning)
- **Infrastructure**: network, compute, cloud, servicemesh
- **Domain**: commerce, workflow, analytics, audit, metering, enterprise
- **AI/ML**: ai, iot, web3
- **Algorithms**: datastructures, algorithms

[Unreleased]: https://github.com/chris-alexander-pop/go-hyperforge/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/chris-alexander-pop/go-hyperforge/releases/tag/v1.0.0
