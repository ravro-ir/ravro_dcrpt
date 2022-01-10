### ravro_dcrpt - Decrypt secret report files

(Written in Go because, you know, "write once, run anywhere.")

### Install Tools 

1 - Install openssl <br />
* Windows : https://slproweb.com/products/Win32OpenSSL.html <br /> 
* Linux   : `apt update && apt install build-essential checkinstall zlib1g-dev openssl libssl-dev -y`

2 - Install wkhtmltopdf 
* Windows : https://wkhtmltopdf.org/downloads.html and add environment variable `C:\Program Files\wkhtmltopdf`
* Linux : `apt-get install wkhtmltopdf`

2 - Rename your private `key` name to `key.private` and copy to `key` folder <br />
3 - Copy `.ravro` to `dataencrypt` folder <br />
4 - Run `ravro_dcrpt.exe` /  `ravro_dcrpt` <br />

### Usage :
```bash
# ./ravro_dcrpt
[++++] Starting for decrypting . . .
[++++] decrypted successfully
[++++] Starting report to pdf . . .
[++++] pdf generated successfully
```



