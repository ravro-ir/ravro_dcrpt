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

# Check if wkhtmltopdf is installed
if ! brew list wkhtmltopdf &> /dev/null; then
    echo "📦 Installing wkhtmltopdf via Homebrew..."
    brew install wkhtmltopdf
else
    echo "✅ wkhtmltopdf is already installed"
fi

# Set OpenSSL paths for both Intel and Apple Silicon Macs
if [[ $(uname -m) == "arm64" ]]; then
    # Apple Silicon Mac
    OPENSSL_PREFIX="/opt/homebrew"
else
    # Intel Mac
    OPENSSL_PREFIX="/usr/local"
fi

echo "🔧 Using OpenSSL from: $OPENSSL_PREFIX"

# Set environment variables
export PKG_CONFIG_PATH="$OPENSSL_PREFIX/lib/pkgconfig:$PKG_CONFIG_PATH"
export CGO_CFLAGS="-I$OPENSSL_PREFIX/include"
export CGO_LDFLAGS="-L$OPENSSL_PREFIX/lib -lssl -lcrypto"

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
