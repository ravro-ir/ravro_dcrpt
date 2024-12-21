#!/bin/bash

set -e

echo "[+++] Updating package lists and installing dependencies..."
sudo apt update -y
sudo apt install -y build-essential checkinstall zlib1g-dev openssl libssl-dev unzip wget xfonts-75dpi

# Attempt to install wkhtmltopdf
echo "Installing wkhtmltopdf..."
sudo apt-get install -y wkhtmltopdf

# Check for unmet dependencies and fix them if needed
if [ $? -ne 0 ]; then
    echo "Resolving unmet dependencies..."
    sudo apt --fix-broken install -y
    echo "Retrying installation of wkhtmltopdf..."
    # Define the URL for the wkhtmltox package
    PACKAGE_URL="https://github.com/wkhtmltopdf/packaging/releases/download/0.12.6-1/wkhtmltox_0.12.6-1.focal_amd64.deb"
    PACKAGE_NAME="wkhtmltox_0.12.6-1.focal_amd64.deb"

    # Download the package
    echo "Downloading $PACKAGE_NAME..."
    wget -O $PACKAGE_NAME $PACKAGE_URL

    # Install the package
    echo "Installing $PACKAGE_NAME..."
    sudo dpkg -i $PACKAGE_NAME

    # Fix any broken dependencies
    echo "Fixing dependencies..."
    sudo apt-get install -f -y
else
    echo "[+++] wkhtmltopdf installed successfully."

    # Define the URL for the wkhtmltox package (same as in the failure block)
    PACKAGE_URL="https://github.com/wkhtmltopdf/packaging/releases/download/0.12.6-1/wkhtmltox_0.12.6-1.focal_amd64.deb"
    PACKAGE_NAME="wkhtmltox_0.12.6-1.focal_amd64.deb"

    # Download the package
    echo "Downloading $PACKAGE_NAME..."
    wget -O $PACKAGE_NAME $PACKAGE_URL

    # Install the package
    echo "Installing $PACKAGE_NAME..."
    sudo dpkg -i $PACKAGE_NAME

    # Fix any broken dependencies
    echo "Fixing dependencies..."
    sudo apt-get install -f -y
fi

echo "[+++]"
echo "[+++] wkhtmltopdf installation process completed."
echo "[+++]"

echo "[+++] Downloading Ravro Decrypt Tools..."
wget -q https://github.com/ravro-ir/ravro_dcrpt/releases/download/v1.0.4/linux_x64_ravro_dcrpt.zip

echo "[+++] Extracting Ravro Decrypt Tools..."
unzip -q -o linux_x64_ravro_dcrpt.zip
echo "[+++] Creating directories..."
mkdir encrypt decrypt key
echo "[+++] Cleanup..."
rm -rf linux_x64_ravro_dcrpt.zip wkhtmltox_0.12.6-1.focal_amd64.deb

echo "[+++] Installation complete!"