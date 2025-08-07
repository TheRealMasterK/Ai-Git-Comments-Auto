#!/bin/bash

# AI Git Auto - Publication Script
# This script helps publish to both Homebrew and npm

set -e

VERSION="1.0.0"
PACKAGE_NAME="ai-git-auto"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
NC='\033[0m' # No Color

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

# ASCII Banner
print_banner() {
    echo -e "${BLUE}"
    echo "┌─────────────────────────────────────────────────────────┐"
    echo "│                                                         │"
    echo "│    📦 AI Git Auto - Publication Assistant 🍺           │"
    echo "│                                                         │"
    echo "│         Publishing to Homebrew and npm                 │"
    echo "│                                                         │"
    echo "└─────────────────────────────────────────────────────────┘"
    echo -e "${NC}\n"
}

# Check prerequisites
check_prerequisites() {
    log_step "Checking prerequisites..."

    # Check if we're in the right directory
    if [[ ! -f "package.json" ]] || [[ ! -f "Formula/ai-git-auto.rb" ]]; then
        log_error "Please run this script from the AI Git Auto project root directory"
        exit 1
    fi

    # Check git status
    if [[ -n $(git status --porcelain) ]]; then
        log_warning "You have uncommitted changes. Consider committing them first."
        echo -e "${YELLOW}Continue anyway? (y/N)${NC}"
        read -r response
        if [[ ! "$response" =~ ^[Yy]$ ]]; then
            log_info "Publication cancelled"
            exit 0
        fi
    fi

    log_success "Prerequisites check passed"
}

# Prepare packages
prepare_packages() {
    log_step "Preparing packages..."

    # Clean and build
    make clean 2>/dev/null || true
    make npm-prepare
    make brew-prepare

    log_success "Packages prepared"
}

# Publish to npm
publish_npm() {
    log_step "Publishing to npm..."

    # Check if logged in
    if ! npm whoami >/dev/null 2>&1; then
        log_error "Not logged into npm. Please run: npm login"
        return 1
    fi

    # Check package name availability
    if npm info "$PACKAGE_NAME" >/dev/null 2>&1; then
        log_error "Package name '$PACKAGE_NAME' is already taken on npm"
        return 1
    fi

    log_info "Publishing to npm registry..."
    npm publish

    log_success "Published to npm! Users can install with: npm install -g $PACKAGE_NAME"
}

# Publish to Homebrew
publish_homebrew() {
    log_step "Setting up Homebrew publication..."

    log_info "To publish to Homebrew, you need to:"
    echo ""
    echo "1. Create a GitHub repository: TheRealMasterK/homebrew-tap"
    echo "2. Copy the formula to that repository:"
    echo "   ${BLUE}cp Formula/ai-git-auto.rb /path/to/homebrew-tap/Formula/${NC}"
    echo "3. Create a GitHub release with tag v$VERSION"
    echo "4. Push the formula to your tap repository"
    echo ""
    log_info "See HOMEBREW-PUBLISH-GUIDE.md for detailed instructions"
}

# Test installations
test_installations() {
    log_step "Testing installations..."

    # Test npm package
    log_info "Testing npm package locally..."
    npm pack

    local tarball="${PACKAGE_NAME}-${VERSION}.tgz"
    if [[ -f "$tarball" ]]; then
        log_success "npm package created: $tarball"

        # Optional: Test install locally
        echo -e "${YELLOW}Test npm installation locally? (y/N)${NC}"
        read -r response
        if [[ "$response" =~ ^[Yy]$ ]]; then
            npm install -g "./$tarball"
            if command -v "$PACKAGE_NAME" >/dev/null; then
                log_success "Local npm installation test passed"
                "$PACKAGE_NAME" --version
            else
                log_error "Local npm installation test failed"
            fi
        fi
    fi
}

# Show publication status
show_status() {
    log_step "Publication Status"

    echo ""
    echo -e "${GREEN}🎉 Package prepared for publication! 🎉${NC}"
    echo ""
    echo -e "${YELLOW}Next steps:${NC}"
    echo ""
    echo "📦 For npm:"
    echo "   1. Make sure you're logged in: ${BLUE}npm login${NC}"
    echo "   2. Publish: ${BLUE}npm publish${NC}"
    echo "   3. Users install with: ${BLUE}npm install -g $PACKAGE_NAME${NC}"
    echo ""
    echo "🍺 For Homebrew:"
    echo "   1. Create homebrew-tap repository on GitHub"
    echo "   2. Copy formula and create release"
    echo "   3. Users install with: ${BLUE}brew install TheRealMasterK/tap/$PACKAGE_NAME${NC}"
    echo ""
    echo -e "${BLUE}See the detailed guides:${NC}"
    echo "   - NPM-PUBLISH-GUIDE.md"
    echo "   - HOMEBREW-PUBLISH-GUIDE.md"
    echo ""
}

# Main function
main() {
    print_banner

    log_info "Starting publication preparation for AI Git Auto v$VERSION"
    echo ""

    check_prerequisites
    prepare_packages
    test_installations

    echo -e "${YELLOW}Proceed with publication? (y/N)${NC}"
    read -r response
    if [[ "$response" =~ ^[Yy]$ ]]; then
        # Try npm publication
        if publish_npm; then
            log_success "npm publication completed!"
        else
            log_warning "npm publication skipped or failed"
        fi

        publish_homebrew
    fi

    show_status

    echo -e "${GREEN}🚀 Ready to share AI Git Auto with the world! 🚀${NC}"
}

# Handle script interruption
trap 'echo -e "\n${RED}Publication interrupted${NC}"; exit 1' INT TERM

# Run main function
main "$@"
