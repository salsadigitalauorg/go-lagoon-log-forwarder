# Go Log Forwarder Makefile
# Provides common development tasks for local development

.PHONY: help test test-race test-cover test-cover-html bench lint fmt vet build clean deps check-deps install-tools all ci

# Default target
.DEFAULT_GOAL := help

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=gofmt
GOVET=$(GOCMD) vet

# Project parameters
BINARY_NAME=go-log-forwarder
COVERAGE_FILE=coverage.out
COVERAGE_HTML=coverage.html

# Colors for output
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[1;33m
BLUE=\033[0;34m
NC=\033[0m # No Color

## help: Show this help message
help:
	@echo "$(BLUE)Go Log Forwarder Development Commands$(NC)"
	@echo ""
	@echo "$(YELLOW)Usage:$(NC)"
	@echo "  make <target>"
	@echo ""
	@echo "$(YELLOW)Available targets:$(NC)"
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*##/ { printf "  $(GREEN)%-15s$(NC) %s\n", $$1, $$2 }' $(MAKEFILE_LIST)
	@echo ""
	@echo "$(YELLOW)Examples:$(NC)"
	@echo "  make test          # Run all tests"
	@echo "  make test-cover    # Run tests with coverage"
	@echo "  make bench         # Run benchmarks"
	@echo "  make lint          # Run all linters"
	@echo "  make ci            # Run full CI pipeline locally"

## test: Run all tests
test:
	@echo "$(BLUE)Running tests...$(NC)"
	$(GOTEST) -v ./...

## test-race: Run tests with race detection
test-race:
	@echo "$(BLUE)Running tests with race detection...$(NC)"
	$(GOTEST) -v -race ./...

## test-short: Run tests in short mode
test-short:
	@echo "$(BLUE)Running tests in short mode...$(NC)"
	$(GOTEST) -v -short ./...

## test-cover: Run tests with coverage
test-cover:
	@echo "$(BLUE)Running tests with coverage...$(NC)"
	$(GOTEST) -v -race -coverprofile=$(COVERAGE_FILE) -covermode=atomic ./...
	@echo "$(GREEN)Coverage report generated: $(COVERAGE_FILE)$(NC)"
	$(GOCMD) tool cover -func=$(COVERAGE_FILE)

## test-cover-html: Run tests with coverage and generate HTML report
test-cover-html: test-cover
	@echo "$(BLUE)Generating HTML coverage report...$(NC)"
	$(GOCMD) tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@echo "$(GREEN)HTML coverage report generated: $(COVERAGE_HTML)$(NC)"
	@echo "$(YELLOW)Open $(COVERAGE_HTML) in your browser to view the report$(NC)"

## bench: Run benchmarks
bench:
	@echo "$(BLUE)Running benchmarks...$(NC)"
	$(GOTEST) -bench=. -benchmem -run=^$$ ./...

## bench-verbose: Run benchmarks with verbose output
bench-verbose:
	@echo "$(BLUE)Running benchmarks with verbose output...$(NC)"
	$(GOTEST) -v -bench=. -benchmem -run=^$$ ./...

## bench-compare: Run benchmarks multiple times for comparison
bench-compare:
	@echo "$(BLUE)Running benchmarks 5 times for comparison...$(NC)"
	$(GOTEST) -bench=. -benchmem -count=5 -run=^$$ ./...

## lint: Run golangci-lint
lint: check-golangci-lint
	@echo "$(BLUE)Running golangci-lint...$(NC)"
	golangci-lint run

## lint-fix: Run golangci-lint with auto-fix
lint-fix: check-golangci-lint
	@echo "$(BLUE)Running golangci-lint with auto-fix...$(NC)"
	golangci-lint run --fix

## fmt: Format Go code
fmt:
	@echo "$(BLUE)Formatting code...$(NC)"
	$(GOFMT) -s -w .
	goimports -w .

## fmt-check: Check if code is formatted correctly
fmt-check:
	@echo "$(BLUE)Checking code formatting...$(NC)"
	@if [ "$$($(GOFMT) -s -l . | wc -l)" -gt 0 ]; then \
		echo "$(RED)The following files are not formatted correctly:$(NC)"; \
		$(GOFMT) -s -l .; \
		echo "$(YELLOW)Run 'make fmt' to fix formatting$(NC)"; \
		exit 1; \
	else \
		echo "$(GREEN)All files are formatted correctly$(NC)"; \
	fi

## vet: Run go vet
vet:
	@echo "$(BLUE)Running go vet...$(NC)"
	$(GOVET) ./...

## security: Run security scan with gosec
security: check-gosec
	@echo "$(BLUE)Running security scan...$(NC)"
	gosec ./...

## build: Build the project
build:
	@echo "$(BLUE)Building project...$(NC)"
	$(GOBUILD) ./...

## clean: Clean build artifacts and test files
clean:
	@echo "$(BLUE)Cleaning up...$(NC)"
	$(GOCLEAN)
	rm -f $(COVERAGE_FILE) $(COVERAGE_HTML)
	rm -rf dist/

## deps: Download and tidy dependencies
deps:
	@echo "$(BLUE)Downloading dependencies...$(NC)"
	$(GOMOD) download
	$(GOMOD) tidy

## deps-update: Update all dependencies
deps-update:
	@echo "$(BLUE)Updating dependencies...$(NC)"
	$(GOGET) -u ./...
	$(GOMOD) tidy

## deps-verify: Verify dependencies
deps-verify:
	@echo "$(BLUE)Verifying dependencies...$(NC)"
	$(GOMOD) verify

## install-tools: Install development tools
install-tools:
	@echo "$(BLUE)Installing development tools...$(NC)"
	@echo "Installing golangci-lint..."
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin; \
	else \
		echo "golangci-lint already installed"; \
	fi
	@echo "Installing gosec..."
	@if ! command -v gosec >/dev/null 2>&1; then \
		go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest; \
	else \
		echo "gosec already installed"; \
	fi
	@echo "Installing goimports..."
	@if ! command -v goimports >/dev/null 2>&1; then \
		go install golang.org/x/tools/cmd/goimports@latest; \
	else \
		echo "goimports already installed"; \
	fi
	@echo "$(GREEN)All tools installed successfully$(NC)"

## check-deps: Check if all dependencies are available
check-deps: check-golangci-lint check-gosec check-goimports

## all: Run fmt, vet, lint, and test
all: fmt vet lint test
	@echo "$(GREEN)All checks passed!$(NC)"

## ci: Run the full CI pipeline locally
ci: deps fmt-check vet lint security test-race bench
	@echo "$(GREEN)CI pipeline completed successfully!$(NC)"

## quick: Quick development check (format, vet, test)
quick: fmt vet test-short
	@echo "$(GREEN)Quick checks completed!$(NC)"

# Helper targets for checking tool availability
check-golangci-lint:
	@command -v golangci-lint >/dev/null 2>&1 || { \
		echo "$(RED)golangci-lint is not installed$(NC)"; \
		echo "$(YELLOW)Run 'make install-tools' to install it$(NC)"; \
		exit 1; \
	}

check-gosec:
	@command -v gosec >/dev/null 2>&1 || { \
		echo "$(RED)gosec is not installed$(NC)"; \
		echo "$(YELLOW)Run 'make install-tools' to install it$(NC)"; \
		exit 1; \
	}

check-goimports:
	@command -v goimports >/dev/null 2>&1 || { \
		echo "$(RED)goimports is not installed$(NC)"; \
		echo "$(YELLOW)Run 'make install-tools' to install it$(NC)"; \
		exit 1; \
	}

# Version information
## version: Show Go version and module information
version:
	@echo "$(BLUE)Go Version:$(NC)"
	@$(GOCMD) version
	@echo ""
	@echo "$(BLUE)Module Information:$(NC)"
	@$(GOMOD) list -m
	@echo ""
	@echo "$(BLUE)Build Information:$(NC)"
	@$(GOCMD) env GOOS GOARCH CGO_ENABLED

## info: Show project information
info:
	@echo "$(BLUE)Project: Go Log Forwarder$(NC)"
	@echo "$(BLUE)Module: $$($(GOMOD) list -m)$(NC)"
	@echo "$(BLUE)Go Version: $$($(GOCMD) version | cut -d' ' -f3)$(NC)"
	@echo "$(BLUE)Files: $$(find . -name '*.go' -not -path './vendor/*' | wc -l) Go files$(NC)"
	@echo "$(BLUE)Lines: $$(find . -name '*.go' -not -path './vendor/*' -exec cat {} + | wc -l) lines of code$(NC)" 