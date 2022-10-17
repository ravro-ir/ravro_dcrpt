# ravro_dcrpt - Decrypt report files of ravro to pdf

(Written in Go because, you know, "write once, run anywhere.")

# Introduction
This is a tool to decrypt reports submitted by a hunter from the Ravro platform bug bounty.

# Manual Install Tools 

1 - Install openssl <br />
  * Windows : https://slproweb.com/products/Win32OpenSSL.html <br /> 
  * Linux (Ubuntu) : `apt update && apt install build-essential checkinstall zlib1g-dev openssl libssl-dev -y`
  * Mac OS : `brew install openssl`<br />

2 - Install wkhtmltopdf 
  * Windows : https://wkhtmltopdf.org/downloads.html and add environment variable `C:\ProgramFiles\wkhtmltopdf`
  * Linux : <br />
        
        $ apt install libssl1.0-dev
        $ apt --fix-broken install
        $ apt install libssl1.0-dev
        $ cd ~
        $ wget https://github.com/wkhtmltopdf/wkhtmltopdf/releases/download/0.12.3/wkhtmltox-0.12.3_linux-generic-amd64.tar.xz
        $ tar vxf wkhtmltox-0.12.3_linux-generic-amd64.tar.xz
        $ cp wkhtmltox/bin/wk* /usr/local/bin/
    
    And you can confirm with:
    
        $ wkhtmltopdf --version
        wkhtmltopdf 0.12.3 (with patched qt)

  * Mac OS : `brew install wkhtmltopdf` <br />

3 - Rename your private `key` name to `key.private` and copy to `key` folder <br />
4 - Download `.zip` file report, Copy `zip` file in the `encrypt` folder and extract it <br />
5 - Run `ravro_dcrpt.exe` /  `ravro_dcrpt` <br />



# Automation Install Tools 

#### Linux 
```bash
root# chmod +x linux.sh
root# ./linux.sh
```
#### Windows
```bash
C:\Users\ravro> win64.bat
```

#### MacOS
```bash
root# ./darwin.sh
```

### Schema 

```bash

.
├── decrypt
│   └── ir2020-07-16-0002
│       └── test__ir2020-07-16-0002__user3.pdf
├── encrypt
│   └── report-ir2020-07-16-0002
│       ├── judgment
│       │   └── data.ravro
│       └── report
│           └── data.ravro
├── key
│   └── key.private

```

# Usage :
Use without command line :
```bash
$ ./ravro_dcrpt -init=init
$ ./ravro_dcrpt
>> Version : ravro_dcrpt/1.0.0
>> HomePage: https://github.com/ravro-ir/ravro_dcrpt
>> Issue: https://github.com/ravro-ir/ravro_dcrpt/issues
>> Author : Ramin Farajpour Cami
>> Help : ravro_dcrpt --help


[++++] Starting for decrypting Report . . .
[++++] Starting for decrypting Judgment . . . 
[++++] Starting for decrypting Amendment . . . 
[++++] decrypted successfully 
[++++] Starting report to pdf . . . 
[++++] pdf generated successfully

```

Use with command line :
```bash
$ ./ravro_dcrpt -init=init
$ ./ravro_dcrpt -in=<INPUT PATH DIR> -out=<OUTPUT PATH DIR> -key=<KEY PATH DIR>\key.private
>> Version : ravro_dcrpt/1.0.0
>> HomePage: https://github.com/ravro-ir/ravro_dcrpt
>> Issue: https://github.com/ravro-ir/ravro_dcrpt/issues
>> Author : Ramin Farajpour Cami
>> Help : ravro_dcrpt --help


[++++] Starting for decrypting Report . . .
[++++] Starting for decrypting Judgment . . . 
[++++] Starting for decrypting Amendment . . . 
[++++] decrypted successfully 
[++++] Starting report to pdf . . . 
[++++] pdf generated successfully

```

# Building from source

Install a [Go compiler](https://golang.org/dl).

Run the following commands in the checked-out repository:
```bash
$ git clone https://github.com/ravro-ir/ravro_dcrpt.git
$ cd ravro_dcrpt
$ go build ravro_dcrpt
$ go run ravro_dcrpt
```
(Add the appropriate `.exe` extension on Windows systems)

## Install on Arch Linux

You can use this package which compiles and installs from latest commit of main branch:

https://aur.archlinux.org/packages/ravro_dcrpt-git/

```
git clone https://aur.archlinux.org/ravro_dcrpt-git.git
cd ravro_dcrpt-git
makepkg -sri
```

# Bugs
Please use github issues to [report](https://github.com/ravro-ir/ravro_dcrpt/issues) bugs.

# License
GNU General Public License, version 3

# Author
Ramin Farajpour Cami <<ramin.blackhat@gmail.com>>, <<farajpour@ravro.ir>>



