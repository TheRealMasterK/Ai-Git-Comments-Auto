# Makefile for AI Git Comments Auto

.PHONY: build test clean install deps run-example global-install uninstall npm-prepare brew-prepare release

# Variables
MAIN_BINARY=ai-git-auto
OLD_BINARY=git-ai-commit
MAIN_CMD_DIR=./cmd/ai-git-auto
OLD_CMD_DIR=./cmd/git-ai-commit
EXAMPLE_DIR=./examples/basic
INSTALL_DIR=/usr/local/bin
VERSION=1.0.0

# Default target
all: deps test build

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod tidy
	go mod download

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Build both CLI tools
build:
	@echo "Building $(MAIN_BINARY)..."
	go build -ldflags "-X main.version=$(VERSION)" -o $(MAIN_BINARY) $(MAIN_CMD_DIR)
	@echo "Building $(OLD_BINARY)..."
	go build -o $(OLD_BINARY) $(OLD_CMD_DIR)

# Build only the main CLI tool
build-main:
	@echo "Building $(MAIN_BINARY)..."
	go build -ldflags "-X main.version=$(VERSION)" -o $(MAIN_BINARY) $(MAIN_CMD_DIR)

# Build for npm package (places binary in bin/ directory)
npm-prepare: deps
	@echo "Preparing npm package..."
	@mkdir -p bin
	go build -ldflags "-X main.version=$(VERSION)" -o bin/$(MAIN_BINARY) $(MAIN_CMD_DIR)
	@chmod +x bin/$(MAIN_BINARY)
	@echo "npm package prepared in bin/ directory"

# Build for Homebrew (creates tarball)
brew-prepare: deps test
	@echo "Preparing Homebrew package..."
	@git archive --format=tar.gz --prefix=ai-git-auto-$(VERSION)/ HEAD > ai-git-auto-$(VERSION).tar.gz
	@echo "Created ai-git-auto-$(VERSION).tar.gz for Homebrew"
	@echo "SHA256: $$(shasum -a 256 ai-git-auto-$(VERSION).tar.gz | cut -d' ' -f1)"

# Release preparation
release: test npm-prepare brew-prepare
	@echo "Release $(VERSION) prepared!"
	@echo ""
	@echo "Installation methods:"
	@echo "  Homebrew: brew install TheRealMasterK/tap/ai-git-auto"
	@echo "  npm:      npm install -g ai-git-auto"
	@echo "  curl:     curl -fsSL https://raw.githubusercontent.com/TheRealMasterK/Ai-Git-Comments-Auto/main/install.sh | bash"

# Install the main CLI tool globally
global-install: build-main
	@echo "Installing $(MAIN_BINARY) to $(INSTALL_DIR)..."
	@if [ ! -w "$(INSTALL_DIR)" ]; then \
		echo "Need sudo permissions to install to $(INSTALL_DIR)"; \
		sudo cp $(MAIN_BINARY) $(INSTALL_DIR)/; \
		sudo chmod +x $(INSTALL_DIR)/$(MAIN_BINARY); \
	else \
		cp $(MAIN_BINARY) $(INSTALL_DIR)/; \
		chmod +x $(INSTALL_DIR)/$(MAIN_BINARY); \
	fi
	@echo "✅ $(MAIN_BINARY) installed successfully!"
	@echo "Run 'ai-git-auto --help' to get started"

# Uninstall the global CLI tool
uninstall:
	@echo "Removing $(MAIN_BINARY) from $(INSTALL_DIR)..."
	@if [ -f "$(INSTALL_DIR)/$(MAIN_BINARY)" ]; then \
		if [ ! -w "$(INSTALL_DIR)" ]; then \
			sudo rm -f $(INSTALL_DIR)/$(MAIN_BINARY); \
		else \
			rm -f $(INSTALL_DIR)/$(MAIN_BINARY); \
		fi; \
		echo "✅ $(MAIN_BINARY) uninstalled successfully!"; \
	else \
		echo "$(MAIN_BINARY) not found in $(INSTALL_DIR)"; \
	fi

# Install to user's local bin (no sudo required)
install-user: build-main
	@echo "Installing $(MAIN_BINARY) to $$HOME/bin..."
	@mkdir -p $$HOME/bin
	@cp $(MAIN_BINARY) $$HOME/bin/
	@chmod +x $$HOME/bin/$(MAIN_BINARY)
	@echo "✅ $(MAIN_BINARY) installed to $$HOME/bin!"
	@echo "Make sure $$HOME/bin is in your PATH"
	@if echo "$$PATH" | grep -q "$$HOME/bin"; then \
		echo "✅ $$HOME/bin is already in PATH"; \
	else \
		echo "⚠️  Add this to your ~/.bashrc or ~/.zshrc:"; \
		echo "export PATH=\"\$$HOME/bin:\$$PATH\""; \
	fi

# Clean build artifacts
clean:
	@echo "Cleaning up..."
	rm -f $(MAIN_BINARY) $(OLD_BINARY)
	go clean

# Run the basic example
run-example:
	@echo "Running basic example..."
	go run $(EXAMPLE_DIR)/main.go

# Quick test run (build and run with dry-run)
test-run: build-main
	@echo "Testing $(MAIN_BINARY) with dry-run..."
	./$(MAIN_BINARY) --dry-run --force

# Run the main CLI tool
run: build-main
	./$(MAIN_BINARY)

# Run with specific options
run-interactive: build-main
	./$(MAIN_BINARY) --interactive

run-auto: build-main
	./$(MAIN_BINARY) --force

run-dry: build-main
	./$(MAIN_BINARY) --dry-run

# Development helpers
dev-setup:
	@echo "Setting up development environment..."
	@echo "Make sure you have the following installed:"
	@echo "1. Go 1.21 or higher"
	@echo "2. Ollama: brew install ollama (macOS) or curl -fsSL https://ollama.ai/install.sh | sh (Linux)"
	@echo "3. An AI model: ollama pull llama2"
	@echo ""
	@echo "Then run: make deps && make build"

# Format code
fmt:
	go fmt ./...

# Lint code (requires golangci-lint)
lint:
	golangci-lint run

# Run with different models
run-codellama: build-main
	./$(MAIN_BINARY) --model codellama

run-mistral: build-main
	./$(MAIN_BINARY) --model mistral

# List available models
list-models: build-main
	./$(MAIN_BINARY) --list-models

# Create release builds for different platforms
release: clean
	@echo "Building release binaries..."
	@mkdir -p release
	GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.version=$(VERSION)" -o release/$(MAIN_BINARY)-darwin-amd64 $(MAIN_CMD_DIR)
	GOOS=darwin GOARCH=arm64 go build -ldflags "-X main.version=$(VERSION)" -o release/$(MAIN_BINARY)-darwin-arm64 $(MAIN_CMD_DIR)
	GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=$(VERSION)" -o release/$(MAIN_BINARY)-linux-amd64 $(MAIN_CMD_DIR)
	GOOS=linux GOARCH=arm64 go build -ldflags "-X main.version=$(VERSION)" -o release/$(MAIN_BINARY)-linux-arm64 $(MAIN_CMD_DIR)
	GOOS=windows GOARCH=amd64 go build -ldflags "-X main.version=$(VERSION)" -o release/$(MAIN_BINARY)-windows-amd64.exe $(MAIN_CMD_DIR)
	@echo "✅ Release binaries created in release/ directory"

# Help
help:
	@echo "Available targets:"
	@echo "  all              - Install deps, run tests, and build"
	@echo "  deps             - Install Go dependencies"
	@echo "  test             - Run unit tests"
	@echo "  build            - Build both CLI tools"
	@echo "  build-main       - Build main CLI tool only"
	@echo "  global-install   - Install CLI tool globally (requires sudo)"
	@echo "  install-user     - Install CLI tool to ~/bin (no sudo)"
	@echo "  uninstall        - Remove globally installed CLI tool"
	@echo "  clean            - Clean build artifacts"
	@echo "  run-example      - Run the basic example"
	@echo "  test-run         - Test run with dry-run mode"
	@echo "  run              - Run CLI tool interactively"
	@echo "  run-auto         - Run CLI tool automatically (no prompts)"
	@echo "  run-dry          - Run CLI tool in dry-run mode"
	@echo "  list-models      - List available Ollama models"
	@echo "  dev-setup        - Show development setup instructions"
	@echo "  fmt              - Format Go code"
	@echo "  lint             - Lint Go code"
	@echo "  release          - Create release binaries for multiple platforms"
	@echo "  help             - Show this help message"
