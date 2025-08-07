#!/bin/bash

# AI Git Auto - One-Click Installation Script
# This script installs ai-git-auto globally with all prerequisites

set -e

BINARY_NAME="ai-git-auto"
REPO_URL="https://github.com/TheRealMasterK/Ai-Git-Comments-Auto"
INSTALL_DIR="/usr/local/bin"
TEMP_DIR="/tmp/ai-git-auto-install"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
NC='\033[0m' # No Color

# ASCII Art Banner
print_banner() {
    echo -e "${BLUE}"
    echo "┌─────────────────────────────────────────────────────────┐"
    echo "│                                                         │"
    echo "│      🚀 AI Git Auto - One-Click Installation 🚀        │"
    echo "│                                                         │"
    echo "│     Automate your Git workflow with AI-powered         │"
    echo "│          commit messages using local Ollama            │"
    echo "│                                                         │"
    echo "└─────────────────────────────────────────────────────────┘"
    echo -e "${NC}\n"
}

# Logging functions
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

log_step() {
    echo -e "\n${PURPLE}🔄 $1${NC}"
}

# Check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check if running on macOS
is_macos() {
    [[ "$(uname)" == "Darwin" ]]
}

# Check if running on Linux
is_linux() {
    [[ "$(uname)" == "Linux" ]]
}

# Install Homebrew (macOS)
install_homebrew() {
    if ! command_exists brew; then
        log_step "Installing Homebrew..."
        /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
        log_success "Homebrew installed successfully!"
    else
        log_info "Homebrew is already installed"
    fi
}

# Install Go
install_go() {
    if ! command_exists go; then
        log_step "Installing Go..."
        if is_macos; then
            install_homebrew
            brew install go
        elif is_linux; then
            if command_exists apt-get; then
                sudo apt-get update
                sudo apt-get install -y golang-go
            elif command_exists yum; then
                sudo yum install -y golang
            elif command_exists pacman; then
                sudo pacman -S go
            else
                log_error "Unable to install Go automatically. Please install Go manually from https://golang.org/dl/"
                exit 1
            fi
        else
            log_error "Unsupported operating system. Please install Go manually from https://golang.org/dl/"
            exit 1
        fi
        log_success "Go installed successfully!"
    else
        log_info "Go is already installed ($(go version))"
    fi
}

# Install Ollama
install_ollama() {
    if ! command_exists ollama; then
        log_step "Installing Ollama..."
        curl -fsSL https://ollama.ai/install.sh | sh
        log_success "Ollama installed successfully!"

        log_step "Starting Ollama service..."
        if is_macos; then
            # On macOS, Ollama should start automatically
            sleep 3
        elif is_linux; then
            # On Linux, we might need to start the service
            if command_exists systemctl; then
                sudo systemctl enable ollama
                sudo systemctl start ollama
            else
                # Start Ollama in background
                nohup ollama serve > /dev/null 2>&1 &
            fi
            sleep 3
        fi

        log_step "Installing recommended AI model (llama3.2:3b - fast and efficient)..."
        ollama pull llama3.2:3b
        log_success "llama3.2:3b model installed!"

        log_info "You can install additional models later with: ollama pull <model-name>"
        log_info "Popular options: codellama:7b, mistral:7b, gemma2:2b"
    else
        log_info "Ollama is already installed"
        log_step "Checking Ollama service..."
        if curl -s http://localhost:11434/api/tags >/dev/null 2>&1; then
            log_success "Ollama service is running"
        else
            log_warning "Ollama is installed but service might not be running"
            log_info "Starting Ollama service..."
            if is_macos; then
                # Try to start Ollama
                ollama serve > /dev/null 2>&1 &
                sleep 3
            elif is_linux; then
                if command_exists systemctl; then
                    sudo systemctl start ollama
                else
                    nohup ollama serve > /dev/null 2>&1 &
                fi
                sleep 3
            fi
        fi
    fi
}

# Check prerequisites
check_prerequisites() {
    log_step "Checking prerequisites..."

    local missing_deps=()

    if ! command_exists git; then
        missing_deps+=("git")
    fi

    if ! command_exists curl; then
        missing_deps+=("curl")
    fi

    if [[ ${#missing_deps[@]} -gt 0 ]]; then
        log_error "Missing required dependencies: ${missing_deps[*]}"
        log_info "Please install these dependencies and run the script again"
        exit 1
    fi

    log_success "All basic prerequisites are available"
}

# Install AI Git Auto
install_ai_git_auto() {
    log_step "Installing AI Git Auto..."

    # Clean up any previous installation attempts
    rm -rf "$TEMP_DIR"
    mkdir -p "$TEMP_DIR"

    # Clone the repository
    log_info "Cloning repository..."
    git clone "$REPO_URL" "$TEMP_DIR"

    # Build and install
    log_info "Building application..."
    cd "$TEMP_DIR"

    # Build the binary
    go mod tidy
    go build -o "$BINARY_NAME" ./cmd/ai-git-auto

    # Install to system
    log_info "Installing to $INSTALL_DIR..."
    sudo mv "$BINARY_NAME" "$INSTALL_DIR/"
    sudo chmod +x "$INSTALL_DIR/$BINARY_NAME"

    # Clean up
    cd /
    rm -rf "$TEMP_DIR"

    log_success "AI Git Auto installed successfully!"
}

# Verify installation
verify_installation() {
    log_step "Verifying installation..."

    if command_exists "$BINARY_NAME"; then
        log_success "✨ AI Git Auto is ready to use!"
        echo -e "\n${GREEN}🎉 Installation completed successfully! 🎉${NC}\n"
    else
        log_error "Installation verification failed"
        exit 1
    fi
}

# Show usage instructions
show_usage() {
    echo -e "${BLUE}"
    echo "┌─────────────────────────────────────────────────────────┐"
    echo "│                    🎯 Quick Start                      │"
    echo "└─────────────────────────────────────────────────────────┘"
    echo -e "${NC}"

    echo -e "${YELLOW}Usage:${NC}"
    echo "  1. Navigate to any Git repository:"
    echo -e "     ${BLUE}cd /path/to/your/project${NC}"
    echo ""
    echo "  2. Run AI Git Auto:"
    echo -e "     ${BLUE}ai-git-auto${NC}"
    echo ""
    echo -e "${YELLOW}What it does:${NC}"
    echo "  • Automatically runs 'git add .'"
    echo "  • Scans your staged changes"
    echo "  • Uses AI to generate a descriptive commit message"
    echo "  • Commits your changes with the AI-generated message"
    echo "  • Pushes to your remote repository"
    echo ""
    echo -e "${YELLOW}Features:${NC}"
    echo "  • Interactive AI model selection"
    echo "  • Detailed logging of each step"
    echo "  • Conventional commit message format"
    echo "  • Analyzes actual code changes for context"
    echo ""
    echo -e "${GREEN}Example output:${NC}"
    echo -e "  ${PURPLE}feat(auth): implement JWT token validation with expiry check${NC}"
    echo -e "  ${PURPLE}fix(api): resolve null pointer exception in user service${NC}"
    echo -e "  ${PURPLE}docs(readme): add installation and usage instructions${NC}"
    echo ""
    echo -e "${BLUE}Need help? Run: ai-git-auto --help${NC}"
    echo ""
}

# Main installation function
main() {
    print_banner

    log_info "Starting AI Git Auto installation..."
    log_info "This script will install Go, Ollama, and AI Git Auto automatically"
    echo ""

    # Ask for confirmation
    echo -e "${YELLOW}Continue with installation? (y/N)${NC}"
    read -r response
    if [[ ! "$response" =~ ^[Yy]$ ]]; then
        log_info "Installation cancelled"
        exit 0
    fi

    echo ""

    # Run installation steps
    check_prerequisites
    install_go
    install_ollama
    install_ai_git_auto
    verify_installation
    show_usage

    echo -e "${GREEN}🚀 Ready to automate your Git workflow with AI! 🚀${NC}\n"
}

# Handle script interruption
trap 'echo -e "\n${RED}Installation interrupted${NC}"; exit 1' INT TERM

# Run main function
main "$@"
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
