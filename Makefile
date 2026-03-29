.PHONY: build install clean test

BINARY_NAME=saayn

build:
	@echo "🔨 Building $(BINARY_NAME)..."
	go build -o $(BINARY_NAME) main.go
	@echo "✅ Build complete."

install: build
	@echo "📦 Installing $(BINARY_NAME) to GOPATH..."
	go install
	@echo "✅ Install complete. You can now run '$(BINARY_NAME)' from anywhere."

clean:
	@echo "🧹 Cleaning up..."
	go clean
	rm -f $(BINARY_NAME)
	@echo "✅ Clean complete."

