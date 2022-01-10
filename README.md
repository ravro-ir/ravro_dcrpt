### ravro_dcrpt - Decrypt secret report files revor

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

# Building from source

Install a [Go compiler](https://golang.org/dl).

Run the following commands in the checked-out repository:
```
go run main.go
Or
go build -o main
```
(Add the appropriate `.exe` extension on Windows systems, of course.)

# License
GNU General Public License, version 3
# Author
Ramin Farajpour Cami <<ramin.blackhat@gmail.com>>



