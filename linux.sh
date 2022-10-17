#!/bin/bash
# Install packages
apt update && apt install build-essential checkinstall zlib1g-dev openssl libssl-dev unzip -y
apt install libssl1.0-dev
apt --fix-broken install
apt install libssl1.0-dev
cd ~
## Download wkhtmltox
wget https://github.com/wkhtmltopdf/wkhtmltopdf/releases/download/0.12.3/wkhtmltox-0.12.3_linux-generic-amd64.tar.xz
tar vxf wkhtmltox-0.12.3_linux-generic-amd64.tar.xz
cp wkhtmltox/bin/wk* /usr/local/bin/
## Download ravro_dcrpt
wget https://github.com/ravro-ir/ravro_dcrpt/releases/download/v1.0.1/linux_x64_ravro_dcrpt.zip
unzip xvf linux_x64_ravro_dcrpt.zip
