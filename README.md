# ravro_dcrpt - Decrypt report files of ravro to pdf

(Written in Go because, you know, "write once, run anywhere.")

# Introduction
This is a tool to decrypt reports submitted by a hunter from the Ravro platform bug bounty.

# Install Tools

1 - Install openssl <br />
  * Linux (Ubuntu) : `apt update && apt install build-essential checkinstall zlib1g-dev openssl libssl-dev -y`
  * Mac OS : `brew install openssl`<br />

2 - Install wkhtmltopdf <br />
    Download and the latest installation package for your system from https://wkhtmltopdf.org/downloads.html.
  * Linux : sudo dpkg -i wkhtmltox.deb  && sudo ldconfig

  * Mac OS : `brew install wkhtmltopdf` <br />

3 - Copy your private key to `key` directory <br />
4 - Download `.zip` file report, Copy `zip` file in the `encrypt` directory.<br />
5 - Run `ravro_dcrpt.exe` /  `ravro_dcrpt` <br />



# Automation Install Tools 

#### [Linux](https://github.com/ravro-ir/ravro_dcrpt/blob/main/autoinstall/linux.sh)
```bash
root# chmod +x linux.sh
root# ./linux.sh
```
#### [Windows](https://github.com/ravro-ir/ravro_dcrpt/blob/main/autoinstall/win64.bat)
```bash
C:\Users\ravro> win64.bat
```

#### [MacOS](https://github.com/ravro-ir/ravro_dcrpt/blob/main/autoinstall/darwin.sh)
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
$ ./ravro_dcrpt -init
$ ./ravro_dcrpt
>> Current Version : ravro_dcrpt/1.0.2
>> Github : https://github.com/ravro-ir/ravro_dcrp
>> Issue : https://github.com/ravro-ir/ravro_dcrp/issues
>> Author : Ravro Development Team (RDT)
>> Help : ravro_dcrpt --help 


Use the arrow keys to navigate: ↓ ↑ → ← 
? Please choose a key: 
  ▸ ravro_key2
    ravro_key1


[++++] Starting for decrypting Report . . . 
[++++] Starting for decrypting Judgment . . . 
[++++] Starting for decrypting Amendment . . . 
[++++] Decrypted successfully 
[++++] Starting report to pdf . . . 
[++++] PDF generated successfully


```

###### Receive latest version :

```bash
$ ./ravro_dcrpt -update
```

###### Receive log of error
```bash
$ ./ravro_dcrpt -log
```

###### Convert report to json
```bash
$ ./ravro_dcrpt -json
```

Use with command line :
```bash
$ ./ravro_dcrpt -init
$ ./ravro_dcrpt -in=<Inout path, /home/irx0xx-xx-xx-000x> -out=<Output path, Ex : /home/output/> -key=<KEY PATH DIR, Ex: key.private>
$ mkdir /home/output
$ mkdir /home/key
$ ./ravro_dcrpt -in=/home/irx0xx-xx-xx-000x -out=/home/output/ -key=/home/key/key.private
>> Current Version : ravro_dcrpt/1.0.2
>> Github : https://github.com/ravro-ir/ravro_dcrp
>> Issue : https://github.com/ravro-ir/ravro_dcrp/issues
>> Author : Ravro Development Team (RDT)
>> Help : ravro_dcrpt --help 


Use the arrow keys to navigate: ↓ ↑ → ← 
? Please choose a key: 
  ▸ ravro_key2
    ravro_key1


[++++] Starting for decrypting Report . . . 
[++++] Starting for decrypting Judgment . . . 
[++++] Starting for decrypting Amendment . . . 
[++++] Decrypted successfully 
[++++] Starting report to pdf . . . 
[++++] PDF generated successfully


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
Building other platform:


##### Build in windows (OpenSSL)
```cmd
$env:CGO_CFLAGS="-IC:/OpenSSL-Win64/include"
$env:CGO_LDFLAGS="-LC:/OpenSSL-Win64/lib/VC/x64/MD -lssl -lcrypto -lws2_32 -lcrypt32"
go build
```


##### Build in windows (OpenSSL / wkhtml2pdf)
```cmd
$env:PATH="C:/OpenSSL-Win64/bin;C:/wkhtmltox/bin;$env:PATH"
$env:CGO_CFLAGS="-IC:/OpenSSL-Win64/include -IC:/wkhtmltox/include"
$env:CGO_LDFLAGS="-LC:/OpenSSL-Win64/lib/VC/x64/MD -LC:/wkhtmltox/lib -L/C:/wkhtmltox/bin -lssl -lcrypto -lws2_32 -lcrypt32 -lwkhtmltox"
go build
```


```bash
$ GOOS=windows GOARCH=amd64 go build .

$ GOOS=darwin GOARCH=amd64 go build .

$ GOOS=linux GOARCH=amd64 go build .
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

# Chagelog
> v1.0.3
* Added multi zip file for decrypting
* Change read argument of multi zip file for decrypting
* Multi select key and refactor code of paths
* Refactor code
* Error handling
* Fixed bugs

> v1.0.2
* Added feature logger
* Added feature spinner load
* Added update ravro_dcrpt
* Better performance pdf result
* Project packaging
* Added convert to json
* Bug fix


# License
GNU General Public License, version 3

# Author
Ramin Farajpour Cami <<ramin.blackhat@gmail.com>>, <<farajpour@ravro.ir>>


