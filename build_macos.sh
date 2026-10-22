#!/bin/bash

# Build script for macOS
# This script handles OpenSSL linking issues on macOS

set -e

echo "🍺 Building Ravro Decryption Tool for macOS..."

# Check if Homebrew is installed
if ! command -v brew &> /dev/null; then
    echo "❌ Homebrew is not installed. Please install it first:"
    echo "   /bin/bash -c \"\$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)\""
    exit 1
fi

# Check if OpenSSL is installed
if ! brew list openssl &> /dev/null; then
    echo "📦 Installing OpenSSL via Homebrew..."
    brew install openssl
else
    echo "✅ OpenSSL is already installed"
fi

# Set OpenSSL paths for both Intel and Apple Silicon Macs
if [[ $(uname -m) == "arm64" ]]; then
    # Apple Silicon Mac
    OPENSSL_PREFIX="/opt/homebrew"
    OPENSSL_VERSION_PATH="/opt/homebrew/Cellar/openssl@3"
else
    # Intel Mac
    OPENSSL_PREFIX="/usr/local"
    OPENSSL_VERSION_PATH="/usr/local/Cellar/openssl@3"
fi

# Find the exact OpenSSL version path
if [ -d "$OPENSSL_VERSION_PATH" ]; then
    OPENSSL_EXACT_PATH=$(find "$OPENSSL_VERSION_PATH" -maxdepth 1 -type d -name "*.*.*" | head -1)
    if [ -n "$OPENSSL_EXACT_PATH" ]; then
        OPENSSL_PREFIX="$OPENSSL_EXACT_PATH"
    fi
fi

echo "🔧 Using OpenSSL from: $OPENSSL_PREFIX"

# Set environment variables with proper paths
export PKG_CONFIG_PATH="$OPENSSL_PREFIX/lib/pkgconfig:$PKG_CONFIG_PATH"
export CGO_CFLAGS="-I$OPENSSL_PREFIX/include"
export CGO_LDFLAGS="-L$OPENSSL_PREFIX/lib -lssl -lcrypto"

# Set macOS deployment target to match OpenSSL version
export MACOSX_DEPLOYMENT_TARGET="14.0"

# Alternative: Try with current system version
# export MACOSX_DEPLOYMENT_TARGET=$(sw_vers -productVersion | cut -d. -f1-2)

# Disable CGO for pure Go build (fallback option)
# export CGO_ENABLED=0

# Create build directory
mkdir -p build

echo "🔨 Building CLI..."
go build -ldflags="-s -w" -o build/ravro_dcrpt ./cmd/cli

echo "🔨 Building GUI..."
go build -ldflags="-s -w" -o build/ravro_dcrpt_gui ./cmd/gui

echo "✅ Build completed successfully!"
echo ""
echo "📁 Built files:"
echo "   - build/ravro_dcrpt (CLI)"
echo "   - build/ravro_dcrpt_gui (GUI)"
echo ""
echo "🚀 To run:"
echo "   ./build/ravro_dcrpt --help"
echo "   ./build/ravro_dcrpt_gui"
