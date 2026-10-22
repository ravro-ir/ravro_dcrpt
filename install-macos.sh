#!/bin/bash

# Ravro Decryption Tool - macOS Installation Script
# This script installs all required dependencies for running the GUI

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}╔═══════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║   Ravro Decryption Tool - Dependency Installation    ║${NC}"
echo -e "${BLUE}║                    macOS                              ║${NC}"
echo -e "${BLUE}╚═══════════════════════════════════════════════════════╝${NC}"
echo ""

# Function to print status
print_status() {
    echo -e "${BLUE}→ $1${NC}"
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

# Check macOS version
MACOS_VERSION=$(sw_vers -productVersion)
ARCH=$(uname -m)

echo -e "${GREEN}✓ macOS: $MACOS_VERSION${NC}"
echo -e "${GREEN}✓ Architecture: $ARCH${NC}"
echo ""

# Check if Homebrew is installed
if ! command -v brew &> /dev/null; then
    print_warning "Homebrew is not installed"
    print_status "Installing Homebrew..."
    /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
    
    # Add Homebrew to PATH for the current session
    if [[ "$ARCH" == "arm64" ]]; then
        eval "$(/opt/homebrew/bin/brew shellenv)"
    else
        eval "$(/usr/local/bin/brew shellenv)"
    fi
    
    print_success "Homebrew installed"
else
    print_success "Homebrew is already installed"
fi

echo ""

# Install OpenSSL
print_status "Installing OpenSSL..."
if ! brew list openssl@3 &>/dev/null; then
    brew install openssl@3
    print_success "OpenSSL@3 installed"
else
    print_success "OpenSSL@3 is already installed"
fi

echo ""
echo -e "${GREEN}╔═══════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║        Dependencies installed successfully! 🎉        ║${NC}"
echo -e "${GREEN}╚═══════════════════════════════════════════════════════╝${NC}"
echo ""
echo -e "${BLUE}Next steps:${NC}"
echo "  1. Download the latest release from GitHub"
echo "  2. For Intel Mac: Download ravro_dcrpt-darwin-amd64.tar.gz"
echo "     For Apple Silicon: Download ravro_dcrpt-darwin-arm64.tar.gz"
echo "  3. Extract the archive:"
echo "     tar -xzf ravro_dcrpt-darwin-*.tar.gz"
echo "  4. Run the .app:"
echo "     open 'Ravro Decryption Tool.app'"
echo ""
echo -e "${YELLOW}Note: First time running, macOS may ask for security permission${NC}"
echo -e "${YELLOW}      Go to System Preferences > Security & Privacy and allow the app${NC}"

