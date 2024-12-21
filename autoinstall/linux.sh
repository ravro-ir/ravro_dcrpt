#!/bin/bash

set -e

echo "[+++] Updating package lists and installing dependencies..."
sudo apt update -y
sudo apt install -y build-essential checkinstall zlib1g-dev openssl libssl-dev unzip wget

# Attempt to install wkhtmltopdf
echo "Installing wkhtmltopdf..."
sudo apt-get install -y wkhtmltopdf

# Check for unmet dependencies and fix them if needed
if [ $? -ne 0 ]; then
    echo "Resolving unmet dependencies..."
    sudo apt --fix-broken install -y
    echo "Retrying installation of wkhtmltopdf..."
    sudo apt-get install -y wkhtmltopdf
fi

# Verify installation
if command -v wkhtmltopdf &> /dev/null; then
    echo "wkhtmltopdf installed successfully."
else
    echo "Failed to install wkhtmltopdf. Please check the errors above."
fi

#echo "[+++] Downloading Ravro Decrypt Tools..."
#wget -q https://github.com/ravro-ir/ravro_dcrpt/releases/download/v1.0.3/linux_x64_ravro_dcrpt.zip

#echo "[+++] Extracting Ravro Decrypt Tools..."
#unzip -q linux_x64_ravro_dcrpt.zip -d ravro_dcrpt

#echo "[+++] Cleanup..."
#rm wkhtmltox_0.12.6.1-2.bullseye_amd64.deb linux_x64_ravro_dcrpt.zip

echo "[+++] Installation complete!"