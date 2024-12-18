#!/bin/bash

set -e

echo "[+++] Updating package lists and installing dependencies..."
sudo apt update && sudo apt install -y build-essential checkinstall zlib1g-dev openssl libssl-dev unzip wget

echo "[+++] Downloading wkhtmltox..."
wget -q https://github.com/wkhtmltopdf/packaging/releases/download/0.12.6.1-2/wkhtmltox_0.12.6.1-2.bullseye_amd64.deb

echo "[+++] Installing wkhtmltox..."
sudo dpkg -i wkhtmltox_0.12.6.1-2.bullseye_amd64.deb
sudo apt-get install -f -y
sudo ldconfig

echo "[+++] Downloading Ravro Decrypt Tools..."
wget -q https://github.com/ravro-ir/ravro_dcrpt/releases/download/v1.0.3/linux_x64_ravro_dcrpt.zip

echo "[+++] Extracting Ravro Decrypt Tools..."
unzip -q linux_x64_ravro_dcrpt.zip -d ravro_dcrpt

echo "[+++] Cleanup..."
rm wkhtmltox_0.12.6.1-2.bullseye_amd64.deb linux_x64_ravro_dcrpt.zip

echo "[+++] Installation complete!"