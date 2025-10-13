#!/bin/bash

# Ravro Decryption Tool - Build Script
# Pure Go implementation with cross-platform support

set -e

VERSION="2.0.0"
APP_NAME="ravro_dcrpt"
BUILD_DIR="build"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}╔═══════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║   Ravro Decryption Tool - Build Script v${VERSION}      ║${NC}"
echo -e "${BLUE}╚═══════════════════════════════════════════════════════╝${NC}"
echo ""

# Function to print colored messages
print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_info() {
    echo -e "${BLUE}→ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

# Check if Go is installed
if ! command -v go &> /dev/null; then
    print_error "Go is not installed. Please install Go 1.21 or higher."
    exit 1
fi

print_success "Go $(go version | awk '{print $3}') detected"

# Create build directory
mkdir -p "$BUILD_DIR"

# Parse command line arguments
BUILD_CLI=true
BUILD_GUI=true
BUILD_LINUX=false
BUILD_WINDOWS=false
BUILD_DARWIN=false
BUILD_KALI=false
BUILD_ALL=false

while [[ $# -gt 0 ]]; do
    case $1 in
        --cli-only)
            BUILD_GUI=false
            shift
            ;;
        --gui-only)
            BUILD_CLI=false
            shift
            ;;
        --linux)
            BUILD_LINUX=true
            shift
            ;;
        --windows)
            BUILD_WINDOWS=true
            shift
            ;;
        --darwin|--macos)
            BUILD_DARWIN=true
            shift
            ;;
        --kali)
            BUILD_KALI=true
            shift
            ;;
        --all)
            BUILD_ALL=true
            BUILD_LINUX=true
            BUILD_WINDOWS=true
            BUILD_DARWIN=true
            BUILD_KALI=true
            shift
            ;;
        --help)
            echo "Usage: $0 [OPTIONS]"
            echo ""
            echo "Options:"
            echo "  --cli-only      Build only CLI"
            echo "  --gui-only      Build only GUI"
            echo "  --linux         Build for Linux"
            echo "  --windows       Build for Windows"
            echo "  --darwin        Build for macOS"
            echo "  --kali          Build for Kali Linux"
            echo "  --all           Build for all platforms"
            echo "  --help          Show this help message"
            echo ""
            exit 0
            ;;
        *)
            print_error "Unknown option: $1"
            echo "Use --help for usage information"
            exit 1
            ;;
    esac
done

# If no platform specified, build for current platform
if [[ "$BUILD_LINUX" == false && "$BUILD_WINDOWS" == false && "$BUILD_DARWIN" == false && "$BUILD_KALI" == false ]]; then
    CURRENT_OS=$(go env GOOS)
    case $CURRENT_OS in
        linux)
            BUILD_LINUX=true
            ;;
        windows)
            BUILD_WINDOWS=true
            ;;
        darwin)
            BUILD_DARWIN=true
            ;;
    esac
fi

# Install dependencies
print_info "Installing dependencies..."
go mod download
go mod tidy
print_success "Dependencies installed"

# Build for Linux
if [[ "$BUILD_LINUX" == true ]]; then
    print_info "Building for Linux..."
    mkdir -p "$BUILD_DIR/linux"
    
    if [[ "$BUILD_CLI" == true ]]; then
        print_info "Building CLI for Linux..."
        GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o "$BUILD_DIR/linux/$APP_NAME" ./cmd/cli
        print_success "Linux CLI built: $BUILD_DIR/linux/$APP_NAME"
    fi
    
    if [[ "$BUILD_GUI" == true ]]; then
        print_info "Building GUI for Linux..."
        GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o "$BUILD_DIR/linux/$APP_NAME-gui" ./cmd/gui
        print_success "Linux GUI built: $BUILD_DIR/linux/$APP_NAME-gui"
    fi
    
    # Create tarball
    if [[ "$BUILD_ALL" == true ]]; then
        print_info "Creating Linux package..."
        cd "$BUILD_DIR/linux" && tar -czf "../$APP_NAME-$VERSION-linux-amd64.tar.gz" * && cd ../..
        print_success "Package created: $BUILD_DIR/$APP_NAME-$VERSION-linux-amd64.tar.gz"
    fi
fi

# Build for Windows
if [[ "$BUILD_WINDOWS" == true ]]; then
    print_info "Building for Windows..."
    mkdir -p "$BUILD_DIR/windows"
    
    if [[ "$BUILD_CLI" == true ]]; then
        print_info "Building CLI for Windows..."
        GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o "$BUILD_DIR/windows/$APP_NAME.exe" ./cmd/cli
        print_success "Windows CLI built: $BUILD_DIR/windows/$APP_NAME.exe"
    fi
    
    if [[ "$BUILD_GUI" == true ]]; then
        print_info "Building GUI for Windows..."
        GOOS=windows GOARCH=amd64 go build -ldflags="-s -w -H=windowsgui" -o "$BUILD_DIR/windows/$APP_NAME-gui.exe" ./cmd/gui
        print_success "Windows GUI built: $BUILD_DIR/windows/$APP_NAME-gui.exe"
    fi
    
    # Create zip
    if [[ "$BUILD_ALL" == true ]]; then
        if command -v zip &> /dev/null; then
            print_info "Creating Windows package..."
            cd "$BUILD_DIR/windows" && zip -r "../$APP_NAME-$VERSION-windows-amd64.zip" * && cd ../..
            print_success "Package created: $BUILD_DIR/$APP_NAME-$VERSION-windows-amd64.zip"
        else
            print_warning "zip command not found, skipping Windows package creation"
        fi
    fi
fi

# Build for macOS
if [[ "$BUILD_DARWIN" == true ]]; then
    print_info "Building for macOS..."
    mkdir -p "$BUILD_DIR/darwin"
    
    if [[ "$BUILD_CLI" == true ]]; then
        print_info "Building CLI for macOS (amd64)..."
        GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o "$BUILD_DIR/darwin/$APP_NAME-amd64" ./cmd/cli
        print_success "macOS CLI (amd64) built: $BUILD_DIR/darwin/$APP_NAME-amd64"
        
        print_info "Building CLI for macOS (arm64)..."
        GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o "$BUILD_DIR/darwin/$APP_NAME-arm64" ./cmd/cli
        print_success "macOS CLI (arm64) built: $BUILD_DIR/darwin/$APP_NAME-arm64"
    fi
    
    if [[ "$BUILD_GUI" == true ]]; then
        print_info "Building GUI for macOS (amd64)..."
        GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o "$BUILD_DIR/darwin/$APP_NAME-gui-amd64" ./cmd/gui
        print_success "macOS GUI (amd64) built: $BUILD_DIR/darwin/$APP_NAME-gui-amd64"
        
        print_info "Building GUI for macOS (arm64)..."
        GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o "$BUILD_DIR/darwin/$APP_NAME-gui-arm64" ./cmd/gui
        print_success "macOS GUI (arm64) built: $BUILD_DIR/darwin/$APP_NAME-gui-arm64"
    fi
    
    # Create tarballs
    if [[ "$BUILD_ALL" == true ]]; then
        print_info "Creating macOS packages..."
        cd "$BUILD_DIR/darwin"
        tar -czf "../$APP_NAME-$VERSION-darwin-amd64.tar.gz" *-amd64
        tar -czf "../$APP_NAME-$VERSION-darwin-arm64.tar.gz" *-arm64
        cd ../..
        print_success "Packages created"
    fi
fi

# Build for Kali Linux (GUI only)
if [[ "$BUILD_KALI" == true ]]; then
    print_info "Building GUI for Kali Linux..."
    mkdir -p "$BUILD_DIR/kali"
    
    if [[ "$BUILD_GUI" == true ]]; then
        print_info "Building GUI for Kali Linux..."
        GOOS=linux GOARCH=amd64 go build -buildvcs=false -ldflags="-s -w" -o "$BUILD_DIR/kali/$APP_NAME-gui" ./cmd/gui
        print_success "Kali GUI built: $BUILD_DIR/kali/$APP_NAME-gui"
    fi
    
    # Create tarball
    if [[ "$BUILD_ALL" == true ]]; then
        print_info "Creating Kali Linux package..."
        cd "$BUILD_DIR/kali" && tar -czf "../$APP_NAME-$VERSION-kali-linux-amd64.tar.gz" * && cd ../..
        print_success "Package created: $BUILD_DIR/$APP_NAME-$VERSION-kali-linux-amd64.tar.gz"
    fi
fi

echo ""
echo -e "${GREEN}╔═══════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║           Build completed successfully! 🎉            ║${NC}"
echo -e "${GREEN}╚═══════════════════════════════════════════════════════╝${NC}"
echo ""
print_info "Build artifacts are in: $BUILD_DIR/"
ls -lh "$BUILD_DIR"

echo ""
echo -e "${BLUE}Next steps:${NC}"
echo "  1. Test the binaries in the $BUILD_DIR directory"
echo "  2. Distribute the binaries or packages as needed"
echo "  3. Run: ./$BUILD_DIR/linux/$APP_NAME --help"

