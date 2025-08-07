#!/bin/bash

# AI Git Auto - Installation Script
# This script installs the ai-git-auto CLI tool globally

set -e

BINARY_NAME="ai-git-auto"
REPO_URL="https://github.com/TheRealMasterK/Ai-Git-Comments-Auto"
INSTALL_DIR="/usr/local/bin"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Helper functions
log_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

log_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

log_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

log_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Check if Go is installed
check_go() {
    if ! command -v go &> /dev/null; then
        log_error "Go is not installed. Please install Go 1.21 or higher."
        log_info "Visit: https://golang.org/doc/install"
        exit 1
    fi

    GO_VERSION=$(go version | cut -d' ' -f3 | sed 's/go//')
    log_success "Go $GO_VERSION found"
}

# Check if Ollama is installed
check_ollama() {
    if ! command -v ollama &> /dev/null; then
        log_warning "Ollama is not installed. The tool requires Ollama to work."
        log_info "To install Ollama:"
        echo "  • macOS: brew install ollama"
        echo "  • Linux: curl -fsSL https://ollama.ai/install.sh | sh"
        echo "  • Or visit: https://ollama.ai"
        echo ""
        read -p "Continue installation anyway? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            exit 1
        fi
    else
        log_success "Ollama found"
    fi
}

# Check if directory is writable
check_permissions() {
    if [ ! -w "$INSTALL_DIR" ] && [ "$EUID" -ne 0 ]; then
        log_error "Cannot write to $INSTALL_DIR. Please run with sudo or choose a different directory."
        log_info "Trying to install to ~/bin instead..."
        INSTALL_DIR="$HOME/bin"
        mkdir -p "$INSTALL_DIR"

        # Add to PATH if not already there
        if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
            log_warning "$INSTALL_DIR is not in your PATH."
            log_info "Add this line to your ~/.bashrc, ~/.zshrc, or ~/.profile:"
            echo "export PATH=\"\$HOME/bin:\$PATH\""
        fi
    fi
}

# Install from source
install_from_source() {
    log_info "Installing from source..."

    # Create temporary directory
    TEMP_DIR=$(mktemp -d)
    cd "$TEMP_DIR"

    # Clone repository
    log_info "Cloning repository..."
    git clone "$REPO_URL" .

    # Build binary
    log_info "Building binary..."
    go build -o "$BINARY_NAME" "./cmd/$BINARY_NAME"

    # Install binary
    log_info "Installing to $INSTALL_DIR..."
    cp "$BINARY_NAME" "$INSTALL_DIR/"
    chmod +x "$INSTALL_DIR/$BINARY_NAME"

    # Cleanup
    cd -
    rm -rf "$TEMP_DIR"

    log_success "Installation completed!"
}

# Install via go install (if available)
install_via_go_install() {
    log_info "Installing via go install..."
    go install "github.com/TheRealMasterK/Ai-Git-Comments-Auto/cmd/$BINARY_NAME@latest"
    log_success "Installation completed via go install!"
}

# Main installation function
main() {
    echo "🚀 AI Git Auto - Installation Script"
    echo "===================================="

    # Check prerequisites
    check_go
    check_ollama
    check_permissions

    # Try go install first, fallback to source
    if command -v go &> /dev/null && go list -m "github.com/TheRealMasterK/Ai-Git-Comments-Auto" &> /dev/null; then
        install_via_go_install
    else
        install_from_source
    fi

    # Verify installation
    if command -v "$BINARY_NAME" &> /dev/null; then
        VERSION=$($BINARY_NAME --version 2>/dev/null || echo "unknown")
        log_success "$BINARY_NAME installed successfully!"
        log_info "Version: $VERSION"
        log_info "Location: $(which $BINARY_NAME)"

        echo ""
        echo "🎉 Setup Complete!"
        echo "=================="
        echo "Next steps:"
        echo "1. Make sure Ollama is running: ollama serve"
        echo "2. Pull an AI model: ollama pull llama2"
        echo "3. Navigate to a Git repository"
        echo "4. Run: $BINARY_NAME"
        echo ""
        echo "For help: $BINARY_NAME --help"

    else
        log_error "Installation failed. Binary not found in PATH."
        exit 1
    fi
}

# Run main function
main "$@"
