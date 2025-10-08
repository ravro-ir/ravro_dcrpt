#!/bin/bash

echo "🔨 Building Windows applications..."

# Install mingw-w64 if not installed
if ! command -v x86_64-w64-mingw32-gcc &> /dev/null; then
    echo "📦 Installing mingw-w64-gcc..."
    sudo apt-get update
    sudo apt-get install -y gcc-mingw-w64-x86-64
fi

# Create build directory
mkdir -p build

# Build CLI for Windows
echo "🏗️  Building ravro_dcrpt.exe (CLI)..."
CGO_ENABLED=1 \
GOOS=windows \
GOARCH=amd64 \
CC=x86_64-w64-mingw32-gcc \
CXX=x86_64-w64-mingw32-g++ \
go build -ldflags="-s -w" \
  -o build/ravro_dcrpt.exe \
  ./cmd/cli

if [ $? -eq 0 ]; then
    echo "✅ CLI built successfully: build/ravro_dcrpt.exe"
    ls -lh build/ravro_dcrpt.exe
else
    echo "❌ CLI build failed!"
    exit 1
fi

# Build GUI for Windows (with hidden console window)
echo "🏗️  Building ravro_dcrpt_gui.exe (GUI)..."
CGO_ENABLED=1 \
GOOS=windows \
GOARCH=amd64 \
CC=x86_64-w64-mingw32-gcc \
CXX=x86_64-w64-mingw32-g++ \
go build -ldflags="-s -w -H windowsgui" \
  -o build/ravro_dcrpt_gui.exe \
  ./cmd/gui

if [ $? -eq 0 ]; then
    echo "✅ GUI built successfully: build/ravro_dcrpt_gui.exe"
    ls -lh build/ravro_dcrpt_gui.exe
    echo ""
    echo "📋 Build Summary:"
    echo "  - CLI: build/ravro_dcrpt.exe"
    echo "  - GUI: build/ravro_dcrpt_gui.exe (console hidden)"
else
    echo "❌ GUI build failed!"
    exit 1
fi
