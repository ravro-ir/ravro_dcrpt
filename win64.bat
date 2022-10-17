// Download OpenSSL
echo [+++] Downloading OpenSSL ....
powershell -command "Invoke-WebRequest -OutFile openssl.exe -Uri https://slproweb.com/download/Win64OpenSSL_Light-3_0_5.exe"

// Check if exit
if exist "C:\Program Files\OpenSSL-Win64\bin\" (
  echo OpenSSL is installed 
) else (
  openssl.exe
)

// Set Path variable envirment
setx PATH ^%PATH^%;"C:\Program Files\OpenSSL-Win64\bin"


// Download wkhtmltox
echo [+++] Downloading wkhtmltox ....
powershell -command "Invoke-WebRequest -OutFile wkhtmltox.exe -Uri https://github.com/wkhtmltopdf/packaging/releases/download/0.12.6-1/wkhtmltox-0.12.6-1.msvc2015-win64.exe"

// Check if exist
if exist "C:\Program Files\wkhtmltopdf\bin\" (
  echo wkhtmltox is installed 
) else (
  wkhtmltox.exe
)
// Set path 
setx PATH ^%PATH^%;"C:\Program Files\wkhtmltopdf\bin"

// Download ravro_dcrpt
echo [+++] Ravro Decrypt Tools ....
powershell -command "Invoke-WebRequest -OutFile ravro_dcrpt.zip -Uri https://github.com/ravro-ir/ravro_dcrpt/releases/download/v1.0.1/win_x64_ravro_dcrpt.zip"

// Extract zip file
powershell -command "Expand-Archive -Force 'ravro_dcrpt.zip' 'ravro_dcrpt'"
