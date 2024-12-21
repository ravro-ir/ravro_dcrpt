@echo off

echo [+++] Creating directories...
mkdir encrypt decrypt key

echo [+++] Downloading OpenSSL Libraries...
curl -L -o libcrypto-3-x64.dll https://github.com/ravro-ir/ravro_dcrpt/raw/main/lib/libcrypto-3-x64.dll
curl -L -o libssl-3-x64.dll https://github.com/ravro-ir/ravro_dcrpt/raw/main/lib/libssl-3-x64.dll

echo [+++] Downloading wkhtmltox Library...
curl -L -o wkhtmltox.dll https://github.com/ravro-ir/ravro_dcrpt/raw/main/lib/wkhtmltox.dll

echo [+++] Downloading Ravro Decrypt Tools...
curl -L -o ravro_dcrpt.zip https://github.com/ravro-ir/ravro_dcrpt/releases/download/v1.0.4/win_x64_ravro_dcrpt.zip

echo [+++] Extracting Ravro Decrypt Tools...
powershell -command "Expand-Archive -Force 'ravro_dcrpt.zip' '.'"

echo [+++] Cleanup...
del ravro_dcrpt.zip

echo [+++] Installation complete!