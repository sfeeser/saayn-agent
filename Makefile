# SAAYN:CHUNK_START:makefile-v1-m1a2k3e4
# BUSINESS_PURPOSE: Build automation for cross-platform distribution and developer workflow.
# SPEC_LINK: SpecBook v1.7 Chapter 0 & 5

BINARY_NAME=saayn
VERSION=1.0.0
BUILD_DIR=bin

.PHONY: all build clean test release help

all: build

## build: Build the binary for the current architecture
build:
	@echo "🏗️  Building $(BINARY_NAME)..."
	@go build -ldflags="-s -w" -o $(BINARY_NAME) main.go
	@echo "✅ Build complete: ./$(BINARY_NAME)"

## test: Run level 1 and level 2 syntax validation tests
test:
	@echo "🧪 Running internal tests..."
	@go test ./internal/... -v

## clean: Remove build artifacts and temporary transaction files
clean:
	@echo "🧹 Cleaning up..."
	@rm -f $(BINARY_NAME)
	@rm -rf $(BUILD_DIR)
	@rm -f .saayn.lock
	@echo "✨ Workspace clean."

## release: Cross-compile for Linux, macOS, and Windows (The Masses Target)
release: clean
	@echo "🌍 Compiling release binaries for v$(VERSION)..."
	mkdir -p $(BUILD_DIR)
	# Linux 64-bit
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 main.go
	# macOS Intel
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 main.go
	# macOS Apple Silicon (M1/M2/M3)
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 main.go
	# Windows 64-bit
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe main.go
	@echo "📦 Release binaries created in ./$(BUILD_DIR)"

## help: Show this help message
help:
	@echo "SAAYN Agent Build System"
	@echo "Usage: make [target]"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' |  sed -e 's/^/ /'

# SAAYN:CHUNK_END:makefile-v1-m1a2k3e4
