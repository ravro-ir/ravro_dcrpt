#!/bin/bash
# Install packages
apt update && apt install build-essential checkinstall zlib1g-dev openssl libssl-dev unzip -y
## Download wkhtmltox
wget https://github.com/wkhtmltopdf/packaging/releases/download/0.12.6.1-2/wkhtmltox_0.12.6.1-2.bullseye_amd64.deb
sudo dpkg -i wkhtmltox_0.12.6.1-2.bullseye_amd64.deb
sudo ldconfig
## Download ravro_dcrpt
wget https://github.com/ravro-ir/ravro_dcrpt/releases/download/v1.0.3/linux_x64_ravro_dcrpt.zip
unzip linux_x64_ravro_dcrpt.zip
