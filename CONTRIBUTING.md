# Contributing to go-hyperforge

Thank you for your interest in contributing! This document provides guidelines for contributing to the project.

## Code of Conduct

Be respectful and constructive in all interactions.

## Getting Started

```bash
# Clone the repository
git clone https://github.com/chris-alexander-pop/go-hyperforge.git
cd go-hyperforge

# Install development tools
make setup

# Start infrastructure (Postgres, Redis, NATS)
make up

# Run quality checks
make check
```

## Development Workflow

### 1. Create a Branch

```bash
git checkout -b feature/your-feature-name
```

### 2. Make Changes

Follow the [Package Standards](pkg/PACKAGE_STANDARDS.md):

- **Interface-first design** - Define interfaces before implementations
- **Decorator pattern** - Use `instrumented.go` for observability
- **Memory adapters** - Always include `adapters/memory/` for testing
- **Context-first** - All I/O methods take `context.Context` first
- **Use `pkg/errors`** - Never use `fmt.Errorf` for errors

### 3. Write Tests

```bash
# Run tests
make test

# Run with coverage
make test-cover
```

Aim for >80% coverage on new code.

### 4. Quality Checks

```bash
# Run all quality gates (required before pushing)
make check
```

This runs:
- `gofmt` formatting check
- `go vet` static analysis
- `staticcheck` linting
- `go build` compilation
- `go test` unit tests

### 5. Submit a Pull Request

- Write a clear PR description
- Reference any related issues
- Ensure CI passes

## Package Structure

When adding new packages, follow this structure:

```
pkg/{capability}/
├── {capability}.go      # Interface, Config, types
├── errors.go            # Sentinel errors (if applicable)
├── instrumented.go      # Observability wrapper (REQUIRED)
├── doc.go               # Package documentation
├── adapters/
│   ├── memory/          # In-memory adapter (REQUIRED)
│   │   └── memory.go
│   └── {driver}/        # Real implementations
└── tests/
    └── {capability}_test.go
```

## Commit Messages

Use conventional commits:

```
feat: add new cache adapter for Memcached
fix: resolve race condition in worker pool
docs: update README with new examples
test: add benchmarks for message bus
refactor: simplify error handling in database package
```

## Questions?

Open an issue for discussion before making large changes.
