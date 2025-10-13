#!/bin/bash

# Ravro Decryption Tool - Kali Linux Installation Script
# Optimized installation for Kali Linux penetration testing environment

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
NC='\033[0m' # No Color

echo -e "${PURPLE}╔═══════════════════════════════════════════════════════╗${NC}"
echo -e "${PURPLE}║   Ravro Decryption Tool - Kali Linux Installer       ║${NC}"
echo -e "${PURPLE}║              🐉 Optimized for Kali 🐉                ║${NC}"
echo -e "${PURPLE}╚═══════════════════════════════════════════════════════╝${NC}"
echo ""

# Check if running as root
if [ "$EUID" -eq 0 ]; then 
    echo -e "${RED}❌ Please do not run this script as root${NC}"
    echo -e "${YELLOW}   Run without sudo, the script will ask for permission when needed${NC}"
    exit 1
fi

# Verify we're on Kali Linux
if [ ! -f /etc/os-release ] || ! grep -q "kali" /etc/os-release; then
    echo -e "${RED}❌ This script is designed for Kali Linux${NC}"
    echo -e "${YELLOW}   For other distributions, use install-linux.sh${NC}"
    exit 1
fi

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

print_info() {
    echo -e "${PURPLE}ℹ $1${NC}"
}

# Detect Kali version
KALI_VERSION=$(grep VERSION_ID /etc/os-release | cut -d'"' -f2)
print_success "Detected Kali Linux $KALI_VERSION"
echo ""

# Update package repositories
print_status "Updating Kali repositories..."
sudo apt-get update -qq
print_success "Repositories updated"

# Install essential development tools
print_status "Installing development tools..."
sudo apt-get install -y \
    build-essential \
    gcc \
    g++ \
    make \
    pkg-config \
    git \
    curl \
    wget
print_success "Development tools installed"

# Install GUI dependencies
print_status "Installing GUI libraries..."
sudo apt-get install -y \
    libgl1-mesa-dev \
    libgl1-mesa-glx \
    xorg-dev \
    libx11-dev \
    libxcursor-dev \
    libxrandr-dev \
    libxinerama-dev \
    libxi-dev \
    libxxf86vm-dev \
    libxss1 \
    libgconf-2-4 \
    libxcomposite1 \
    libasound2-dev
print_success "GUI libraries installed"

# Install OpenSSL
print_status "Installing OpenSSL..."
sudo apt-get install -y \
    libssl-dev \
    libssl3 \
    openssl
print_success "OpenSSL installed"

# Install PDF generation tools
print_status "Installing PDF tools..."
sudo apt-get install -y wkhtmltopdf
print_success "PDF tools installed"

# Install additional Kali-specific tools that might be useful
print_status "Installing additional forensics tools..."
sudo apt-get install -y \
    hexedit \
    binwalk \
    file \
    strings \
    xxd
print_success "Forensics tools installed"

# Check if Go is installed
if ! command -v go &> /dev/null; then
    print_warning "Go is not installed. Installing Go..."
    
    # Download and install Go
    GO_VERSION="1.23.0"
    print_status "Downloading Go $GO_VERSION..."
    wget -q "https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz" -O /tmp/go.tar.gz
    
    print_status "Installing Go..."
    sudo rm -rf /usr/local/go
    sudo tar -C /usr/local -xzf /tmp/go.tar.gz
    
    # Add Go to PATH
    if ! grep -q "/usr/local/go/bin" ~/.bashrc; then
        echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
        echo 'export GOPATH=$HOME/go' >> ~/.bashrc
        echo 'export PATH=$PATH:$GOPATH/bin' >> ~/.bashrc
    fi
    
    # Add to current session
    export PATH=$PATH:/usr/local/go/bin
    export GOPATH=$HOME/go
    export PATH=$PATH:$GOPATH/bin
    
    rm /tmp/go.tar.gz
    print_success "Go $GO_VERSION installed"
else
    GO_VERSION=$(go version | awk '{print $3}')
    print_success "Go $GO_VERSION already installed"
fi

# Create directories for Kali tools
print_status "Setting up directories..."
mkdir -p ~/kali-tools/ravro_dcrpt
mkdir -p ~/.local/bin
print_success "Directories created"

# Download latest release
print_status "Downloading latest Ravro Decryption Tool..."
LATEST_URL="https://api.github.com/repos/ravro-ir/ravro_dcrpt/releases/latest"
DOWNLOAD_URL=$(curl -s $LATEST_URL | grep "browser_download_url.*kali-linux-amd64.tar.gz" | cut -d '"' -f 4)

if [ -z "$DOWNLOAD_URL" ]; then
    print_warning "Kali-specific build not found, downloading regular Linux build..."
    DOWNLOAD_URL=$(curl -s $LATEST_URL | grep "browser_download_url.*linux-amd64.tar.gz" | head -1 | cut -d '"' -f 4)
fi

if [ -n "$DOWNLOAD_URL" ]; then
    wget -q "$DOWNLOAD_URL" -O /tmp/ravro_dcrpt.tar.gz
    cd ~/kali-tools/ravro_dcrpt
    tar -xzf /tmp/ravro_dcrpt.tar.gz
    chmod +x ravro_dcrpt*
    
    # Create symlinks in ~/.local/bin
    ln -sf ~/kali-tools/ravro_dcrpt/ravro_dcrpt ~/.local/bin/ravro_dcrpt
    if [ -f ~/kali-tools/ravro_dcrpt/ravro_dcrpt_gui ]; then
        ln -sf ~/kali-tools/ravro_dcrpt/ravro_dcrpt_gui ~/.local/bin/ravro_dcrpt_gui
    fi
    
    rm /tmp/ravro_dcrpt.tar.gz
    print_success "Ravro Decryption Tool installed"
else
    print_error "Could not download release. Please check internet connection."
    exit 1
fi

# Add ~/.local/bin to PATH if not already there
if ! grep -q "$HOME/.local/bin" ~/.bashrc; then
    echo 'export PATH=$HOME/.local/bin:$PATH' >> ~/.bashrc
fi

# Test installation
print_status "Testing installation..."
cd ~/kali-tools/ravro_dcrpt
if ./ravro_dcrpt --help > /dev/null 2>&1; then
    print_success "CLI tool working correctly"
else
    print_warning "CLI tool may have issues"
fi

if [ -f ./ravro_dcrpt_gui ]; then
    print_success "GUI tool available"
else
    print_warning "GUI tool not found"
fi

echo ""
echo -e "${GREEN}╔═══════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║     Kali Linux installation completed! 🎉🐉          ║${NC}"
echo -e "${GREEN}╚═══════════════════════════════════════════════════════╝${NC}"
echo ""

print_info "Installation location: ~/kali-tools/ravro_dcrpt/"
print_info "Binaries linked to: ~/.local/bin/"
echo ""

echo -e "${BLUE}Usage Examples:${NC}"
echo "  # CLI usage"
echo "  ravro_dcrpt decrypt encrypted.ravro output.pdf --key private.key"
echo "  ravro_dcrpt info encrypted.ravro"
echo ""
echo "  # GUI usage (if X11/Wayland available)"
echo "  ravro_dcrpt_gui"
echo ""
echo "  # Direct execution"
echo "  cd ~/kali-tools/ravro_dcrpt"
echo "  ./ravro_dcrpt --help"
echo ""

echo -e "${YELLOW}Kali-specific notes:${NC}"
echo "  • For headless environments, use CLI mode"
echo "  • For GUI in SSH, use X11 forwarding: ssh -X user@host"
echo "  • For VNC access, ensure X11 server is running"
echo "  • Tools are in ~/kali-tools/ravro_dcrpt/ for easy access"
echo ""

print_info "Restart your terminal or run: source ~/.bashrc"
print_success "Ready for penetration testing! 🔍"
