@REM @echo off

@REM echo [+++] Creating directories...
@REM mkdir ravro_dcrpt
@REM cd ravro_dcrpt
@REM mkdir encrypt decrypt key

@REM echo [+++] Downloading OpenSSL Libraries...
@REM curl -L -o libcrypto-3-x64.dll https://github.com/ravro-ir/ravro_dcrpt/raw/main/lib/libcrypto-3-x64.dll
@REM curl -L -o libssl-3-x64.dll https://github.com/ravro-ir/ravro_dcrpt/raw/main/lib/libssl-3-x64.dll

@REM echo [+++] Downloading wkhtmltox Library...
@REM curl -L -o wkhtmltox.dll https://github.com/ravro-ir/ravro_dcrpt/raw/main/lib/wkhtmltox.dll

@REM echo [+++] Downloading Ravro Decrypt Tools...
@REM curl -L -o ravro_dcrpt.zip https://github.com/ravro-ir/ravro_dcrpt/releases/download/v1.0.4/win_x64_ravro_dcrpt.zip

@REM echo [+++] Extracting Ravro Decrypt Tools...
@REM powershell -command "Expand-Archive -Force 'ravro_dcrpt.zip' '.'"

@REM echo [+++] Cleanup...
@REM del ravro_dcrpt.zip

@REM echo [+++] Installation complete!
