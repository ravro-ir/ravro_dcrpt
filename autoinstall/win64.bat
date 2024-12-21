@echo off

echo [+++] Downloading OpenSSL Libraries...
powershell -command "& {
    $ProgressPreference = 'SilentlyContinue'
    Invoke-WebRequest -OutFile libcrypto-3-x64.dll -Uri https://github.com/ravro-ir/ravro_dcrpt/raw/main/lib/libcrypto-3-x64.dll
    Invoke-WebRequest -OutFile libssl-3-x64.dll -Uri https://github.com/ravro-ir/ravro_dcrpt/raw/main/lib/libssl-3-x64.dll
}"

echo [+++] Downloading wkhtmltox Library...
powershell -command "& {
    $ProgressPreference = 'SilentlyContinue'
    Invoke-WebRequest -OutFile wkhtmltox.dll -Uri https://github.com/ravro-ir/ravro_dcrpt/raw/main/lib/wkhtmltox.dll
}"

echo [+++] Downloading Ravro Decrypt Tools...
powershell -command "& {
    $ProgressPreference = 'SilentlyContinue'
    Invoke-WebRequest -OutFile ravro_dcrpt.zip -Uri https://github.com/ravro-ir/ravro_dcrpt/releases/download/v1.0.4/ravro_dcrpt.zip
}"

echo [+++] Extracting Ravro Decrypt Tools...
powershell -command "Expand-Archive -Force 'ravro_dcrpt.zip' '.'"
echo [+++] Cleanup...
del ravro_dcrpt.zip

echo [+++] Installation complete!