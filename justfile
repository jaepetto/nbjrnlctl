# nbjrnlctl Justfile
# Date: 2025-09-04
# Description: Automation tasks for the Netbox Journal CLI Tool

# Default recipe - show help
default:
	just --list

# Build the application (statically linked)
build:
	CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o nbjrnlctl ./cmd/nbjrnlctl

# Install the application globally (statically linked)
install:
	CGO_ENABLED=0 go install -a -ldflags '-extldflags "-static"' ./cmd/nbjrnlctl

# Run the application
run:
	go run ./cmd/nbjrnlctl

# Run with arguments (usage: just run-with-args list --limit 5)
run-with-args *args="":
	go run ./cmd/nbjrnlctl {{args}}

# Build for all supported platforms (statically linked)
build-all:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-extldflags "-static"' -o bin/nbjrnlctl-linux-amd64 ./cmd/nbjrnlctl
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -a -ldflags '-extldflags "-static"' -o bin/nbjrnlctl-linux-arm64 ./cmd/nbjrnlctl
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -a -ldflags '-extldflags "-static"' -o bin/nbjrnlctl-darwin-amd64 ./cmd/nbjrnlctl
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -a -ldflags '-extldflags "-static"' -o bin/nbjrnlctl-darwin-arm64 ./cmd/nbjrnlctl
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -a -ldflags '-extldflags "-static"' -o bin/nbjrnlctl-windows-amd64.exe ./cmd/nbjrnlctl

# Clean build artifacts
clean:
	rm -f nbjrnlctl
	rm -rf bin/

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	open coverage.html

# Format code
fmt:
	go fmt ./...

# Vet code for issues
vet:
	go vet ./...

# Tidy go modules
tidy:
	go mod tidy

# Download dependencies
deps:
	go mod download

# Update dependencies
update-deps:
	go get -u ./...

# Check for security vulnerabilities
audit:
	gosec ./...

# Generate documentation
docs:
	go doc -all ./...

# Show version information
version:
	go version

# List all available commands
help:
	@echo "Available commands:"
	@echo "  build          - Build the application (statically linked)"
	@echo "  install        - Install the application globally (statically linked)"
	@echo "  run            - Run the application"
	@echo "  run-with-args  - Run with custom arguments"
	@echo "  build-all      - Build for all platforms (statically linked)"
	@echo "  clean          - Clean build artifacts"
	@echo "  test           - Run tests"
	@echo "  test-cover     - Run tests with coverage"
	@echo "  fmt            - Format code"
	@echo "  vet            - Vet code for issues"
	@echo "  tidy           - Tidy go modules"
	@echo "  deps           - Download dependencies"
	@echo "  update-deps    - Update dependencies"
	@echo "  audit          - Check for security vulnerabilities"
	@echo "  docs           - Generate documentation"
	@echo "  version        - Show version information"
	@echo "  help           - Show this help"
	@echo ""
	@echo "Application usage examples:"
	@echo "  just run list                      # List journal entries"
	@echo "  just run-with-args list --limit 5  # List with arguments"
