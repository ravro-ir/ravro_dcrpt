#!/bin/bash

apt update && apt install build-essential checkinstall zlib1g-dev openssl libssl-dev -y
apt install libssl1.0-dev
apt --fix-broken install
apt install libssl1.0-dev
cd ~
wget https://github.com/wkhtmltopdf/wkhtmltopdf/releases/download/0.12.3/wkhtmltox-0.12.3_linux-generic-amd64.tar.xz
tar vxf wkhtmltox-0.12.3_linux-generic-amd64.tar.xz
cp wkhtmltox/bin/wk* /usr/local/bin/

