# go-common

[![CI](https://github.com/soderasen-au/go-common/workflows/CI/badge.svg)](https://github.com/soderasen-au/go-common/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/soderasen-au/go-common)](https://goreportcard.com/report/github.com/soderasen-au/go-common)
[![codecov](https://codecov.io/gh/soderasen-au/go-common/branch/main/graph/badge.svg)](https://codecov.io/gh/soderasen-au/go-common)
[![Go Reference](https://pkg.go.dev/badge/github.com/soderasen-au/go-common.svg)](https://pkg.go.dev/github.com/soderasen-au/go-common)
[![License](https://img.shields.io/github/license/soderasen-au/go-common)](LICENSE)

A collection of reusable Go packages providing common utilities for security, authentication, logging, async execution, and more.

## Features

- **🔐 Crypto** - RSA encryption, TLS certificates, and key management
- **⚡ Exec** - Asynchronous task execution framework with state management
- **📊 FX** - Expression evaluation engine for dynamic calculations
- **📝 Loggers** - Structured logging built on zerolog with CSV support and rotation
- **🔑 OAuth** - OAuth 1.0 and 2.0 authentication handling
- **🛡️ SAML** - SAML 2.0 SSO integration
- **🛠️ Util** - Foundation utilities (Result-based error handling, pointers, files, diff, templates)

## Requirements

- Go 1.25 or higher

## Installation

```bash
go get github.com/soderasen-au/go-common
```

## Quick Start

### Result-Based Error Handling

```go
import "github.com/soderasen-au/go-common/util"

func YourFunction(ctx string) *util.Result {
    if err := something(); err != nil {
        return util.Error(ctx, err)
    }
    return util.OK(ctx)
}
```

### Expression Evaluation

```go
import "github.com/soderasen-au/go-common/fx"

// Create an expression: 5 < 10
exp, _ := fx.NewBinOpExp(fx.LessOp, fx.Number(5), fx.Number(10))
result := exp.Calc() // Returns fx.True()
```

### Logging

```go
import "github.com/soderasen-au/go-common/loggers"

logger := loggers.GetLogger("app.log")
logger.Info().Msg("Application started")
```

### Cryptography

```go
import "github.com/soderasen-au/go-common/crypto"

// Use internal RSA key pair
cipher, _ := crypto.InternalEncypt([]byte("secret"))
text, _ := crypto.InternalDecrypt(cipher)
```

## Packages

### crypto

Cryptographic utilities for RSA encryption, TLS certificates, and key management.

**Key Features:**
- RSA public/private key operations
- X509 certificate handling (PEM, PKCS#12)
- TLS configuration generation
- Multiple certificate sources (files, inline PEM, PKCS#12)

**Coverage:** 75.9%

### exec

Asynchronous task execution framework with state management.

**Key Features:**
- Request/Response pattern
- Async execution orchestration
- Status tracking (Ready → Running → Ok/Failed)
- In-memory request keeper

**Coverage:** 0.0% (needs tests)

### fx

Expression evaluation engine for dynamic calculations.

**Key Features:**
- Type-safe value system (numeric, text, boolean, error)
- Comparison operators (<, >, ==, <=, >=, !=)
- Text operations (contains ~=)
- Logical operations (AND, OR, NOT)
- Extensible operator registry

**Coverage:** 99.3%

### loggers

Structured logging built on zerolog with CSV support and log rotation.

**Key Features:**
- File logging with rotation (lumberjack)
- CSV output for audit trails
- Global configuration
- Console and JSON output modes

**Coverage:** 0.0% (needs tests)

### oauth

OAuth 1.0 and 2.0 authentication handling.

**Key Features:**
- OAuth 2.0 Authorization Code flow
- OAuth 2.0 Client Credentials flow
- OAuth 1.0 with RSA signatures
- Token exchange and refresh

**Coverage:** 0.0% (needs tests)

### sso/saml

SAML 2.0 SSO integration with IdP metadata loading.

**Key Features:**
- IdP metadata from file, URL, or inline
- Service provider configuration
- Certificate integration
- SAML response validation

**Coverage:** 0.0% (needs tests)

### util

Foundation utilities used across all packages.

**Key Features:**
- Result-based error handling pattern
- Pointer helpers (Ptr, MaybeNil, MaybeDefault)
- File operations (Exists, ListFiles, FilterFiles, MoveFile)
- Object comparison (Diff, JsonDiff)
- Template rendering with custom functions
- JSON marshaling utilities

**Coverage:** 75.6%

## Development

### Prerequisites

```bash
# Install Go 1.25+
# Install golangci-lint (optional, for linting)
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
```

### Building

```bash
make build
```

### Testing

```bash
# Run all tests
make test

# Run tests with coverage
make coverage

# Generate HTML coverage report
make coverage-html

# Run tests for specific package
make test-util
make test-fx
make test-crypto
```

### Code Quality

```bash
# Format code
make fmt

# Check formatting
make fmt-check

# Run go vet
make vet

# Run linters
make lint

# Run all checks
make check
```

### CI Pipeline

```bash
# Run full CI pipeline locally
make ci
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests (`make test`)
5. Run checks (`make check`)
6. Commit your changes (`git commit -m 'Add amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

### Coding Standards

- All functions should return `*util.Result` for consistent error handling
- Write tests for new functionality (aim for >75% coverage)
- Run `make fmt` before committing
- Ensure `make check` passes
- Follow existing code patterns

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Project Structure

```
.
├── crypto/         # Cryptographic utilities
├── exec/           # Async task execution
├── fx/             # Expression evaluation
├── loggers/        # Structured logging
├── oauth/          # OAuth authentication
├── sso/saml/       # SAML SSO integration
├── util/           # Foundation utilities
├── Makefile        # Build automation
└── .github/
    └── workflows/
        └── ci.yml  # GitHub Actions CI
```

## Makefile Targets

Run `make help` to see all available targets:

```
Development:
  test         - Run all tests
  test-v       - Run all tests with verbose output
  test-short   - Run short tests only
  coverage     - Run tests with coverage report
  coverage-html - Generate HTML coverage report
  build        - Build all packages
  clean        - Clean build artifacts and cache

Code Quality:
  fmt          - Format code with go fmt
  fmt-check    - Check if code is formatted
  vet          - Run go vet
  lint         - Run linters (requires golangci-lint)
  check        - Run fmt, vet, and tests

Dependencies:
  deps         - Download dependencies
  tidy         - Tidy and verify dependencies
  verify       - Verify dependencies

CI/CD:
  ci           - Run full CI pipeline (deps, check, coverage)
```
