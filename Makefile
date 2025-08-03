# MAKEFILE
#
# @author      Nicola Asuni <info@tecnick.com>
# @link        https://github.com/advancedlogic/GoOse
#
# Modern Makefile for Go modules project
#
# ------------------------------------------------------------------------------

# List special make targets that are not associated with files
.PHONY: help all build build-cli install test test-verbose test-race format fmt-check vet lint coverage qa deps clean nuke tidy run version

# Use bash as shell (Note: Ubuntu now uses dash which doesn't support PIPESTATUS).
SHELL=/bin/bash

# Project details
PROJECT=GoOse
BINARY_NAME=goose
VERSION=$(shell cat VERSION)
MODULE=github.com/advancedlogic/GoOse

# Build details
BUILD_DIR=bin
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -X main.GitCommit=${GIT_COMMIT}"

# Test settings
COVERAGE_DIR=coverage
COVERAGE_PROFILE=${COVERAGE_DIR}/coverage.out
COVERAGE_HTML=${COVERAGE_DIR}/coverage.html

# --- MAKE TARGETS ---

# Display general help about this command
help:
	@echo ""
	@echo "$(PROJECT) Makefile - Go Modules Project"
	@echo "Version: $(VERSION)"
	@echo ""
	@echo "Available commands:"
	@echo ""
	@echo "  BUILD COMMANDS:"
	@echo "    make build       : Build the CLI binary"
	@echo "    make build-cli   : Build the CLI binary (alias)"
	@echo "    make install     : Install the CLI binary to GOPATH/bin"
	@echo ""
	@echo "  DEVELOPMENT COMMANDS:"
	@echo "    make run         : Run the CLI with arguments (use ARGS=...)"
	@echo "    make version     : Show version information"
	@echo ""
	@echo "  TESTING COMMANDS:"
	@echo "    make test        : Run all tests"
	@echo "    make test-verbose: Run tests with verbose output"
	@echo "    make test-race   : Run tests with race detection"
	@echo "    make coverage    : Generate test coverage report"
	@echo ""
	@echo "  CODE QUALITY COMMANDS:"
	@echo "    make format      : Format source code"
	@echo "    make fmt-check   : Check if source code is formatted"
	@echo "    make vet         : Run go vet"
	@echo "    make lint        : Run golangci-lint"
	@echo "    make qa          : Run all quality assurance checks"
	@echo ""
	@echo "  DEPENDENCY COMMANDS:"
	@echo "    make deps        : Download dependencies"
	@echo "    make tidy        : Tidy go.mod and go.sum"
	@echo ""
	@echo "  CLEANUP COMMANDS:"
	@echo "    make clean       : Remove build artifacts"
	@echo "    make nuke        : Remove all generated files"
	@echo ""
	@echo "Examples:"
	@echo "  make build"
	@echo "  make run ARGS='convert https://example.com'"
	@echo "  make test"
	@echo ""

# Alias for help target
all: help

# Build the CLI binary
build: build-cli

build-cli:
	@echo "Building $(BINARY_NAME) v$(VERSION)..."
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/goose

# Install the CLI binary
install:
	@echo "Installing $(BINARY_NAME)..."
	go install $(LDFLAGS) ./cmd/goose

# Run the CLI (use ARGS="..." to pass arguments)
run:
	@go run ./cmd/goose $(ARGS)

# Show version information
version:
	@echo "Project: $(PROJECT)"
	@echo "Version: $(VERSION)"
	@echo "Module:  $(MODULE)"
	@echo "Git:     $(GIT_COMMIT)"

# Run the unit tests
test:
	@echo "Running tests..."
	go test ./...

# Run tests with verbose output
test-verbose:
	@echo "Running tests with verbose output..."
	go test -v ./...

# Run tests with race detection
test-race:
	@echo "Running tests with race detection..."
	go test -race ./...

# Generate test coverage report
coverage:
	@echo "Generating coverage report..."
	@mkdir -p $(COVERAGE_DIR)
	go test -coverprofile=$(COVERAGE_PROFILE) ./...
	go tool cover -html=$(COVERAGE_PROFILE) -o $(COVERAGE_HTML)
	@echo "Coverage report generated: $(COVERAGE_HTML)"

# Format the source code
format:
	@echo "Formatting source code..."
	go fmt ./...

# Check if the source code has been formatted
fmt-check:
	@echo "Checking code formatting..."
	@if [ "$$(gofmt -s -l . | wc -l)" -gt 0 ]; then \
		echo "The following files need formatting:"; \
		gofmt -s -l .; \
		echo "Please run 'make format' to fix formatting issues."; \
		exit 1; \
	else \
		echo "All files are properly formatted."; \
	fi

# Check for suspicious constructs
vet:
	@echo "Running go vet..."
	go vet ./...

# Run golangci-lint (requires golangci-lint to be installed)
lint:
	@echo "Running golangci-lint..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found. Please install it:"; \
		echo "  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
		echo "Or run: make deps"; \
	fi

# Run all quality assurance checks
qa: fmt-check vet lint test coverage
	@echo "All quality assurance checks completed!"

# --- DEPENDENCIES AND CLEANUP ---

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download
	@echo "Installing development tools..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Tidy go.mod and go.sum
tidy:
	@echo "Tidying go.mod and go.sum..."
	go mod tidy

# Remove build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	rm -rf $(COVERAGE_DIR)
	go clean ./...

# Remove all generated files
nuke: clean
	@echo "Removing all generated files..."
	go clean -cache -testcache -modcache
	rm -rf vendor/
