# Makefile for go-common project
# A reusable Go common library

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet

# Project parameters
MODULE=github.com/soderasen-au/go-common
PACKAGES=$(shell go list ./...)
COVERAGE_FILE=coverage.out
COVERAGE_HTML=coverage.html

# Colors for output
COLOR_RESET=\033[0m
COLOR_BOLD=\033[1m
COLOR_GREEN=\033[32m
COLOR_YELLOW=\033[33m
COLOR_BLUE=\033[34m

.PHONY: all help clean test coverage build fmt vet lint deps tidy verify check ci

# Default target
all: check build

## help: Display this help message
help:
	@echo "$(COLOR_BOLD)Available targets:$(COLOR_RESET)"
	@echo ""
	@echo "  $(COLOR_GREEN)all$(COLOR_RESET)          - Run check and build (default)"
	@echo "  $(COLOR_GREEN)help$(COLOR_RESET)         - Display this help message"
	@echo ""
	@echo "$(COLOR_BOLD)Development:$(COLOR_RESET)"
	@echo "  $(COLOR_GREEN)test$(COLOR_RESET)         - Run all tests"
	@echo "  $(COLOR_GREEN)test-v$(COLOR_RESET)       - Run all tests with verbose output"
	@echo "  $(COLOR_GREEN)test-short$(COLOR_RESET)   - Run short tests only"
	@echo "  $(COLOR_GREEN)coverage$(COLOR_RESET)     - Run tests with coverage report"
	@echo "  $(COLOR_GREEN)coverage-html$(COLOR_RESET) - Generate HTML coverage report"
	@echo "  $(COLOR_GREEN)build$(COLOR_RESET)        - Build all packages"
	@echo "  $(COLOR_GREEN)clean$(COLOR_RESET)        - Clean build artifacts and cache"
	@echo ""
	@echo "$(COLOR_BOLD)Code Quality:$(COLOR_RESET)"
	@echo "  $(COLOR_GREEN)fmt$(COLOR_RESET)          - Format code with go fmt"
	@echo "  $(COLOR_GREEN)fmt-check$(COLOR_RESET)    - Check if code is formatted"
	@echo "  $(COLOR_GREEN)vet$(COLOR_RESET)          - Run go vet"
	@echo "  $(COLOR_GREEN)lint$(COLOR_RESET)         - Run linters (requires golangci-lint)"
	@echo "  $(COLOR_GREEN)check$(COLOR_RESET)        - Run fmt, vet, and tests"
	@echo ""
	@echo "$(COLOR_BOLD)Dependencies:$(COLOR_RESET)"
	@echo "  $(COLOR_GREEN)deps$(COLOR_RESET)         - Download dependencies"
	@echo "  $(COLOR_GREEN)tidy$(COLOR_RESET)         - Tidy and verify dependencies"
	@echo "  $(COLOR_GREEN)verify$(COLOR_RESET)       - Verify dependencies"
	@echo ""
	@echo "$(COLOR_BOLD)CI/CD:$(COLOR_RESET)"
	@echo "  $(COLOR_GREEN)ci$(COLOR_RESET)           - Run full CI pipeline (deps, check, coverage)"
	@echo ""

## clean: Remove build artifacts and cache
clean:
	@echo "$(COLOR_BLUE)Cleaning build artifacts...$(COLOR_RESET)"
	$(GOCLEAN)
	rm -f $(COVERAGE_FILE) $(COVERAGE_HTML)
	@echo "$(COLOR_GREEN)✓ Clean complete$(COLOR_RESET)"

## deps: Download dependencies
deps:
	@echo "$(COLOR_BLUE)Downloading dependencies...$(COLOR_RESET)"
	$(GOMOD) download
	@echo "$(COLOR_GREEN)✓ Dependencies downloaded$(COLOR_RESET)"

## tidy: Tidy and verify dependencies
tidy:
	@echo "$(COLOR_BLUE)Tidying dependencies...$(COLOR_RESET)"
	$(GOMOD) tidy
	@echo "$(COLOR_GREEN)✓ Dependencies tidied$(COLOR_RESET)"

## verify: Verify dependencies
verify:
	@echo "$(COLOR_BLUE)Verifying dependencies...$(COLOR_RESET)"
	$(GOMOD) verify
	@echo "$(COLOR_GREEN)✓ Dependencies verified$(COLOR_RESET)"

## build: Build all packages
build:
	@echo "$(COLOR_BLUE)Building all packages...$(COLOR_RESET)"
	$(GOBUILD) ./...
	@echo "$(COLOR_GREEN)✓ Build complete$(COLOR_RESET)"

## fmt: Format code with go fmt
fmt:
	@echo "$(COLOR_BLUE)Formatting code...$(COLOR_RESET)"
	$(GOFMT) ./...
	@echo "$(COLOR_GREEN)✓ Code formatted$(COLOR_RESET)"

## fmt-check: Check if code is formatted
fmt-check:
	@echo "$(COLOR_BLUE)Checking code formatting...$(COLOR_RESET)"
	@unformatted=$$(gofmt -l .); \
	if [ -n "$$unformatted" ]; then \
		echo "$(COLOR_YELLOW)⚠ The following files are not formatted:$(COLOR_RESET)"; \
		echo "$$unformatted"; \
		exit 1; \
	fi
	@echo "$(COLOR_GREEN)✓ All files are formatted$(COLOR_RESET)"

## vet: Run go vet
vet:
	@echo "$(COLOR_BLUE)Running go vet...$(COLOR_RESET)"
	$(GOVET) ./...
	@echo "$(COLOR_GREEN)✓ Vet complete$(COLOR_RESET)"

## lint: Run golangci-lint
lint:
	@echo "$(COLOR_BLUE)Running linters...$(COLOR_RESET)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
		echo "$(COLOR_GREEN)✓ Linting complete$(COLOR_RESET)"; \
	else \
		echo "$(COLOR_YELLOW)⚠ golangci-lint not installed, skipping...$(COLOR_RESET)"; \
		echo "  Install with: curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin"; \
	fi

## test: Run all tests
test:
	@echo "$(COLOR_BLUE)Running tests...$(COLOR_RESET)"
	$(GOTEST) -race ./...
	@echo "$(COLOR_GREEN)✓ Tests complete$(COLOR_RESET)"

## test-v: Run all tests with verbose output
test-v:
	@echo "$(COLOR_BLUE)Running tests (verbose)...$(COLOR_RESET)"
	$(GOTEST) -v -race ./...

## test-short: Run short tests only
test-short:
	@echo "$(COLOR_BLUE)Running short tests...$(COLOR_RESET)"
	$(GOTEST) -short ./...
	@echo "$(COLOR_GREEN)✓ Short tests complete$(COLOR_RESET)"

## coverage: Run tests with coverage report
coverage:
	@echo "$(COLOR_BLUE)Running tests with coverage...$(COLOR_RESET)"
	$(GOTEST) -race -coverprofile=$(COVERAGE_FILE) -covermode=atomic ./...
	@echo ""
	@echo "$(COLOR_BOLD)Coverage Summary:$(COLOR_RESET)"
	@$(GOCMD) tool cover -func=$(COVERAGE_FILE) | grep total | awk '{print "  Total Coverage: " $$3}'
	@echo ""
	@echo "$(COLOR_GREEN)✓ Coverage report generated: $(COVERAGE_FILE)$(COLOR_RESET)"

## coverage-html: Generate HTML coverage report
coverage-html: coverage
	@echo "$(COLOR_BLUE)Generating HTML coverage report...$(COLOR_RESET)"
	$(GOCMD) tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@echo "$(COLOR_GREEN)✓ HTML coverage report: $(COVERAGE_HTML)$(COLOR_RESET)"
	@if command -v xdg-open >/dev/null 2>&1; then \
		xdg-open $(COVERAGE_HTML); \
	elif command -v open >/dev/null 2>&1; then \
		open $(COVERAGE_HTML); \
	else \
		echo "  Open $(COVERAGE_HTML) in your browser"; \
	fi

## check: Run fmt, vet, and tests
check: fmt-check vet test
	@echo "$(COLOR_GREEN)✓ All checks passed$(COLOR_RESET)"

## ci: Run full CI pipeline
ci: deps verify fmt-check vet coverage
	@echo ""
	@echo "$(COLOR_BOLD)$(COLOR_GREEN)========================================$(COLOR_RESET)"
	@echo "$(COLOR_BOLD)$(COLOR_GREEN)  ✓ CI Pipeline Complete$(COLOR_RESET)"
	@echo "$(COLOR_BOLD)$(COLOR_GREEN)========================================$(COLOR_RESET)"
	@echo ""

# Package-specific test targets
.PHONY: test-util test-fx test-crypto test-exec test-loggers test-oauth test-saml

## test-util: Run tests for util package
test-util:
	@echo "$(COLOR_BLUE)Running util tests...$(COLOR_RESET)"
	$(GOTEST) -v -race -cover ./util

## test-fx: Run tests for fx package
test-fx:
	@echo "$(COLOR_BLUE)Running fx tests...$(COLOR_RESET)"
	$(GOTEST) -v -race -cover ./fx

## test-crypto: Run tests for crypto package
test-crypto:
	@echo "$(COLOR_BLUE)Running crypto tests...$(COLOR_RESET)"
	$(GOTEST) -v -race -cover ./crypto

## test-exec: Run tests for exec package
test-exec:
	@echo "$(COLOR_BLUE)Running exec tests...$(COLOR_RESET)"
	$(GOTEST) -v -race -cover ./exec

## test-loggers: Run tests for loggers package
test-loggers:
	@echo "$(COLOR_BLUE)Running loggers tests...$(COLOR_RESET)"
	$(GOTEST) -v -race -cover ./loggers

## test-oauth: Run tests for oauth package
test-oauth:
	@echo "$(COLOR_BLUE)Running oauth tests...$(COLOR_RESET)"
	$(GOTEST) -v -race -cover ./oauth

## test-saml: Run tests for saml package
test-saml:
	@echo "$(COLOR_BLUE)Running saml tests...$(COLOR_RESET)"
	$(GOTEST) -v -race -cover ./sso/saml
