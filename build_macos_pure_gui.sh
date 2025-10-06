#!/bin/bash

# Pure Go build for GUI only (fallback when OpenSSL linking fails)

set -e

echo "🍺 Building Ravro GUI for macOS (Pure Go - No OpenSSL)..."

# Check if wkhtmltopdf is installed
if ! command -v wkhtmltopdf &> /dev/null; then
    echo "⚠️  wkhtmltopdf not found. Please install it first:"
    echo "   brew install --cask wkhtmltopdf"
    exit 1
fi

# Disable CGO to use pure Go PKCS7 implementation
export CGO_ENABLED=0

# Create build directory
mkdir -p build

echo "🔨 Building GUI (Pure Go - No OpenSSL)..."
go build -ldflags="-s -w" -o build/ravro_dcrpt_gui ./cmd/gui

echo "✅ Pure Go GUI build completed successfully!"
echo ""
echo "📁 Built file:"
echo "   - build/ravro_dcrpt_gui (GUI - Pure Go)"
echo ""
echo "ℹ️  Note: This build uses pure Go PKCS7 implementation"
echo "🚀 To run:"
echo "   ./build/ravro_dcrpt_gui"
