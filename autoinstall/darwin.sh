#!/bin/bash

set -e

echo "[+++] Checking for Homebrew..."
if ! command -v brew &> /dev/null; then
    echo "Homebrew not found. Please install Homebrew first."
    echo "Visit https://brew.sh for installation instructions."
    exit 1
fi

echo "[+++] Updating Homebrew and installing dependencies..."
brew update
brew install unzip openssl wkhtmltopdf wget

echo "[+++] Downloading Ravro Decrypt Tools..."
wget -q https://github.com/ravro-ir/ravro_dcrpt/releases/download/v1.0.3/macos_x64_ravro_dcrpt.zip

echo "[+++] Extracting Ravro Decrypt Tools..."
unzip -q macos_x64_ravro_dcrpt.zip -d ravro_dcrpt

echo "[+++] Cleanup..."
rm macos_x64_ravro_dcrpt.zip

echo "[+++] Installation complete!"