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
brew install unzip openssl wget pkg-config

echo "[+++] Downloading Ravro Decrypt Tools..."
wget -q https://github.com/ravro-ir/ravro_dcrpt/releases/download/v1.0.4/macos_x64_ravro_dcrpt.zip

echo "[+++] Extracting Ravro Decrypt Tools..."
unzip -q -o macos_x64_ravro_dcrpt.zip

echo "[+++] Creating directories..."
mkdir encrypt decrypt key

echo "[+++] Cleanup..."
rm macos_x64_ravro_dcrpt.zip

echo "[+++] Installation complete!"