# ravro_dcrpt - Decrypt report files of ravro to pdf

(Written in Go because, you know, "write once, run anywhere.")

# Introduction
This is a tool to decrypt reports submitted by a hunter from the Ravro platform bug bounty.

# Install Tools 

1 - Install openssl <br />
* Windows : https://slproweb.com/products/Win32OpenSSL.html <br /> 
* Linux (Ubuntu) : `apt update && apt install build-essential checkinstall zlib1g-dev openssl libssl-dev -y`
* Mac OS : `brew install openssl` <br />
2 - Install wkhtmltopdf 
* Windows : https://wkhtmltopdf.org/downloads.html and add environment variable `C:\Program Files\wkhtmltopdf`
* Linux (Ubuntu) : `apt-get install wkhtmltopdf`
* Mac OS : `brew install wkhtmltopdf` <br />

3 - Rename your private `key` name to `key.private` and copy to `key` folder <br />
4 - Copy `.zip` report and copy `zip` file in in the `encrypt` folder and extract it <br />
5 - Run `ravro_dcrpt.exe` /  `ravro_dcrpt` <br />


# Usage :
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
```bash
$ git clone https://github.com/ravro-ir/ravro_dcrp.git
$ go build ravro_dcrpt
$ go run ravro_dcrpt
```
(Add the appropriate `.exe` extension on Windows systems)

# Bugs
Please use github issues to [report](https://github.com/ravro-ir/ravro_dcrp/issues) bugs.

# License
GNU General Public License, version 3

# Author
Ramin Farajpour Cami <<ramin.blackhat@gmail.com>>, <<farajpour@ravro.ir>>



