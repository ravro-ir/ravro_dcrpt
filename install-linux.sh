#!/bin/bash

# Ravro Decryption Tool - Linux Installation Script
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
echo -e "${BLUE}║                    Linux                              ║${NC}"
echo -e "${BLUE}╚═══════════════════════════════════════════════════════╝${NC}"
echo ""

# Check if running as root
if [ "$EUID" -eq 0 ]; then 
    echo -e "${RED}❌ Please do not run this script as root${NC}"
    echo -e "${YELLOW}   Run without sudo, the script will ask for permission when needed${NC}"
    exit 1
fi

# Detect Linux distribution
if [ -f /etc/os-release ]; then
    . /etc/os-release
    OS=$ID
    VER=$VERSION_ID
else
    echo -e "${RED}❌ Cannot detect Linux distribution${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Detected: $PRETTY_NAME${NC}"
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

# Install dependencies based on distribution
case "$OS" in
    ubuntu|debian|linuxmint|pop|kali)
        print_status "Installing dependencies for Ubuntu/Debian/Kali..."
        
        sudo apt-get update
        
        print_status "Installing system libraries..."
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
            libssl-dev \
            libssl3 \
            pkg-config
        
        print_status "Installing wkhtmltopdf..."
        if ! command -v wkhtmltopdf &> /dev/null; then
            sudo apt-get install -y wkhtmltopdf
            print_success "wkhtmltopdf installed"
        else
            print_success "wkhtmltopdf already installed"
        fi
        ;;
        
    fedora|rhel|centos|rocky|almalinux)
        print_status "Installing dependencies for Fedora/RHEL..."
        
        print_status "Installing system libraries..."
        sudo dnf install -y \
            mesa-libGL-devel \
            libX11-devel \
            libXcursor-devel \
            libXrandr-devel \
            libXinerama-devel \
            libXi-devel \
            libXxf86vm-devel \
            openssl-devel \
            pkgconfig
        
        print_status "Installing wkhtmltopdf..."
        if ! command -v wkhtmltopdf &> /dev/null; then
            sudo dnf install -y wkhtmltopdf
            print_success "wkhtmltopdf installed"
        else
            print_success "wkhtmltopdf already installed"
        fi
        ;;
        
    arch|manjaro)
        print_status "Installing dependencies for Arch Linux..."
        
        print_status "Installing system libraries..."
        sudo pacman -S --noconfirm \
            mesa \
            libx11 \
            libxcursor \
            libxrandr \
            libxinerama \
            libxi \
            libxxf86vm \
            openssl \
            pkgconf
        
        print_status "Installing wkhtmltopdf..."
        if ! command -v wkhtmltopdf &> /dev/null; then
            sudo pacman -S --noconfirm wkhtmltopdf
            print_success "wkhtmltopdf installed"
        else
            print_success "wkhtmltopdf already installed"
        fi
        ;;
        
    opensuse*|sles)
        print_status "Installing dependencies for openSUSE..."
        
        print_status "Installing system libraries..."
        sudo zypper install -y \
            Mesa-libGL-devel \
            libX11-devel \
            libXcursor-devel \
            libXrandr-devel \
            libXinerama-devel \
            libXi-devel \
            libXxf86vm-devel \
            libopenssl-devel \
            pkg-config
        
        print_status "Installing wkhtmltopdf..."
        if ! command -v wkhtmltopdf &> /dev/null; then
            sudo zypper install -y wkhtmltopdf
            print_success "wkhtmltopdf installed"
        else
            print_success "wkhtmltopdf already installed"
        fi
        ;;
        
    *)
        print_error "Unsupported distribution: $OS"
        echo ""
        echo "Please install the following packages manually:"
        echo "  - OpenSSL development libraries"
        echo "  - X11 development libraries (libX11, libXcursor, libXrandr, libXinerama, libXi, libXxf86vm)"
        echo "  - Mesa/OpenGL libraries"
        echo "  - wkhtmltopdf"
        exit 1
        ;;
esac

echo ""
echo -e "${GREEN}╔═══════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║        Dependencies installed successfully! 🎉        ║${NC}"
echo -e "${GREEN}╚═══════════════════════════════════════════════════════╝${NC}"
echo ""
echo -e "${BLUE}Next steps:${NC}"
echo "  1. Download the latest release from GitHub"
echo "  2. Extract the archive:"
echo "     tar -xzf ravro_dcrpt-linux-amd64.tar.gz"
echo "  3. Run the GUI:"
echo "     ./ravro_dcrpt_gui"
echo ""
echo -e "${YELLOW}Note: Make sure you have a display server (X11 or Wayland) running${NC}"

