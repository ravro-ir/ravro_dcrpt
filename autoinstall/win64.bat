// Download OpenSSL
echo [+++] Downloading OpenSSL Lib ....
powershell -command "Invoke-WebRequest -OutFile openssl.exe -Uri https://github.com/ravro-ir/ravro_dcrpt/blob/develop/lib/libcrypto-3-x64.dll"
powershell -command "Invoke-WebRequest -OutFile openssl.exe -Uri https://github.com/ravro-ir/ravro_dcrpt/blob/develop/lib/libssl-3-x64.dll"

// Download wkhtmltox
echo [+++] Downloading wkhtmltox Lib....
powershell -command "Invoke-WebRequest -OutFile wkhtmltox.exe -Uri https://github.com/ravro-ir/ravro_dcrpt/blob/develop/lib/wkhtmltox.dll"

// Download ravro_dcrpt
echo [+++] Ravro Decrypt Tools ....
powershell -command "Invoke-WebRequest -OutFile ravro_dcrpt.zip -Uri https://github.com/ravro-ir/ravro_dcrpt/releases/download/v1.0.3/win_x64_ravro_dcrpt.zip"

// Extract zip file
powershell -command "Expand-Archive -Force 'ravro_dcrpt.zip' 'ravro_dcrpt'"
