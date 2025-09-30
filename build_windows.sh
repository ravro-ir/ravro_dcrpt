#!/bin/bash

echo "🔨 Building CLI for Windows..."

# Install mingw-w64 if not installed
if ! command -v x86_64-w64-mingw32-gcc &> /dev/null; then
    echo "📦 Installing mingw-w64-gcc..."
    sudo apt-get update
    sudo apt-get install -y gcc-mingw-w64-x86-64
fi

# Build CLI for Windows
echo "🏗️  Building ravro_dcrpt.exe..."
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
    echo "❌ Build failed!"
    exit 1
fi
