# ravro_dcrpt

A versatile Go-based tool for decrypting and converting Ravro platform bug bounty reports to PDF.

## 🚀 Introduction

`ravro_dcrpt` is a cross-platform utility designed to decrypt and process reports submitted by hunters on the Ravro platform, embodying the "write once, run anywhere" philosophy of Go.

## ✨ Features

- 🔐 Decrypt encrypted Ravro report files
- 📄 Convert decrypted reports to PDF
- 🗝️ Multi-key support
- 🔄 Multi-zip file decryption
- 📋 JSON conversion option
- 🖥️ Cross-platform compatibility (Windows, Linux, macOS)
- 🆕 Built-in update mechanism
- 🐞 Comprehensive error logging

## 🛠️ Installation


### Automated Installation

#### Linux
```bash
wget https://raw.githubusercontent.com/ravro-ir/ravro_dcrpt/refs/heads/main/autoinstall/linux.sh -O - | sh
```

#### Windows
```bash
Ps> curl -L -o win64.bat https://raw.githubusercontent.com/ravro-ir/ravro_dcrpt/refs/heads/main/autoinstall/win64.bat

Ps> .\win64.bat
[+++] Downloading OpenSSL Libraries...
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
  0     0    0     0    0     0      0      0 --:--:-- --:--:-- --:--:--     0
100 6445k  100 6445k    0     0  2571k      0  0:00:02  0:00:02 --:--:-- 8572k
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
  0     0    0     0    0     0      0      0 --:--:-- --:--:-- --:--:--     0
100 1224k  100 1224k    0     0  46645      0  0:00:26  0:00:26 --:--:-- 51891
[+++] Downloading wkhtmltox Library...
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
  0     0    0     0    0     0      0      0 --:--:-- --:--:-- --:--:--     0
100 28.7M  100 28.7M    0     0  2247k      0  0:00:13  0:00:13 --:--:-- 1839k
[+++] Downloading Ravro Decrypt Tools...
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
  0     0    0     0    0     0      0      0 --:--:-- --:--:-- --:--:--     0
100 6100k  100 6100k    0     0  33418      0  0:03:06  0:03:06 --:--:-- 29295
[+++] Extracting Ravro Decrypt Tools...
[+++] Cleanup...
[+++] Installation complete!

```

#### macOS (Soon)
```bash
wget https://raw.githubusercontent.com/ravro-ir/ravro_dcrpt/refs/heads/main/autoinstall/darwin.sh -O - | sh
```

## 📂 Project Structure

```
.
├── decrypt
│   └── ir2020-07-16-0002
│       └── test__ir2020-07-16-0002__user3.pdf
├── encrypt
│   └── report-ir2020-07-16-0002
│       ├── judgment
│       │   └── data.ravro
│       └── report
│           └── data.ravro
└── key
    └── key.private
```

## 💻 Usage

### Interactive Mode
```bash
$ ./ravro_dcrpt -init
$ ./ravro_dcrpt
```

### Command-Line Mode
```bash
$ ./ravro_dcrpt -init
$ ./ravro_dcrpt -in=<input_path> -out=<output_path> -key=<key_path>
```

### Additional Commands
- Update to latest version:
  ```bash
  $ ./ravro_dcrpt -update
  ```
- View error logs:
  ```bash
  $ ./ravro_dcrpt -log
  ```
- Convert report to JSON:
  ```bash
  $ ./ravro_dcrpt -json
  ```

## 🔨 Building from Source

### Prerequisites
- [Go compiler](https://golang.org/dl)

### Standard Build (Linux)
```bash
$ git clone https://github.com/ravro-ir/ravro_dcrpt.git
$ cd ravro_dcrpt
$ go build ravro_dcrpt
$ go run ravro_dcrpt
```

### Cross-Platform Builds (Developing)

#### Build for Windows with OpenSSL
```powershell
Ps> $env:CGO_CFLAGS="-IC:/OpenSSL-Win64/include"
Ps> $env:CGO_LDFLAGS="-LC:/OpenSSL-Win64/lib/VC/x64/MD -lssl -lcrypto -lws2_32 -lcrypt32"
Ps> go build
```

#### Build for Windows with OpenSSL and wkhtmltopdf
```powershell
Ps> $env:PATH="C:/OpenSSL-Win64/bin;C:/wkhtmltox/bin;$env:PATH"
Ps> $env:CGO_CFLAGS="-IC:/OpenSSL-Win64/include -IC:/wkhtmltox/include"
Ps> $env:CGO_LDFLAGS="-LC:/OpenSSL-Win64/lib/VC/x64/MD -LC:/wkhtmltox/lib -L/C:/wkhtmltox/bin -lssl -lcrypto -lws2_32 -lcrypt32 -lwkhtmltox"
Ps> go build
```

## 🐧 Arch Linux Installation
```bash
git clone https://aur.archlinux.org/ravro_dcrpt-git.git
cd ravro_dcrpt-git
makepkg -sri
```

## 📋 Changelog

### v1.0.3
- Multi-zip file decryption support
- Improved key selection process
- Enhanced path handling
- Comprehensive error handling
- Code refactoring

### v1.0.2
- Added logging capabilities
- Implemented loading spinner
- Added update functionality
- Improved PDF generation performance
- JSON conversion support
- Project packaging
- Bug fixes

## 📄 License

GNU General Public License, version 3

## 👥 Author

Ramin Farajpour Cami
- Email: ramin.blackhat@gmail.com
- Alternate Email: farajpour@ravro.ir

