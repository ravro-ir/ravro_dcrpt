# Building ravro_dcrpt from Source

This guide provides comprehensive instructions for building `ravro_dcrpt` from source on different platforms.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Building on Windows](#building-on-windows)
- [Building on Linux](#building-on-linux)
- [Building on macOS](#building-on-macos)
- [Cross-Compilation](#cross-compilation)
- [Using Makefile](#using-makefile)
- [Troubleshooting](#troubleshooting)

---

## Prerequisites

### All Platforms

- **Go 1.21.1 or later**: [Download Go](https://golang.org/dl)
- **Git**: For cloning the repository

### Windows

- **OpenSSL for Windows**: [Download](https://slproweb.com/products/Win32OpenSSL.html)
  - Install to `C:\OpenSSL-Win64`
- **wkhtmltopdf**: [Download](https://wkhtmltopdf.org/downloads.html)
  - Install to `C:\wkhtmltox`

### Linux

```bash
# Ubuntu/Debian
sudo apt-get update
sudo apt-get install -y libssl-dev wkhtmltopdf

# Fedora/RHEL
sudo dnf install -y openssl-devel wkhtmltopdf

# Arch Linux
sudo pacman -S openssl wkhtmltopdf
```

### macOS

```bash
# Using Homebrew
brew install openssl wkhtmltopdf
```

---

## Building on Windows

### 1. Clone Repository

```powershell
git clone https://github.com/ravro-ir/ravro_dcrpt.git
cd ravro_dcrpt
```

### 2. Set Environment Variables

```powershell
$env:PATH="C:/OpenSSL-Win64/bin;C:/wkhtmltox/bin;$env:PATH"
$env:CGO_CFLAGS="-IC:/OpenSSL-Win64/include -IC:/wkhtmltox/include"
$env:CGO_LDFLAGS="-LC:/OpenSSL-Win64/lib/VC/x64/MD -LC:/wkhtmltox/lib -L/C:/wkhtmltox/bin -lssl -lcrypto -lws2_32 -lcrypt32 -lwkhtmltox"
```

### 3. Build

**Build CLI:**
```powershell
go build -ldflags="-s -w" -o ravro_dcrpt.exe ./cmd/cli
```

**Build GUI (without console window):**
```powershell
go build -ldflags="-s -w -H windowsgui" -o ravro_dcrpt_gui.exe ./cmd/gui
```

### Output Files
- `ravro_dcrpt.exe` - CLI application
- `ravro_dcrpt_gui.exe` - GUI application (no console window)

**Note:** The `-H windowsgui` flag hides the console window for GUI applications on Windows.

---

## Building on Linux

### 1. Clone Repository

```bash
git clone https://github.com/ravro-ir/ravro_dcrpt.git
cd ravro_dcrpt
```

### 2. Install Dependencies

```bash
# Ubuntu/Debian
sudo apt-get update
sudo apt-get install -y libssl-dev wkhtmltopdf

# Arch Linux
sudo pacman -S openssl wkhtmltopdf
```

### 3. Build

**Build CLI:**
```bash
go build -ldflags="-s -w" -o ravro_dcrpt ./cmd/cli
```

**Build GUI:**
```bash
go build -ldflags="-s -w" -o ravro_dcrpt_gui ./cmd/gui
```

### Output Files
- `ravro_dcrpt` - CLI application
- `ravro_dcrpt_gui` - GUI application

---

## Building on macOS

### 1. Clone Repository

```bash
git clone https://github.com/ravro-ir/ravro_dcrpt.git
cd ravro_dcrpt
```

### 2. Install Dependencies

```bash
brew install openssl wkhtmltopdf
```

### 3. Set Environment Variables

```bash
export PKG_CONFIG_PATH=$(brew --prefix openssl)/lib/pkgconfig
export CGO_CFLAGS="-I$(brew --prefix openssl)/include"
export CGO_LDFLAGS="-L$(brew --prefix openssl)/lib"
```

### 4. Build

**Build CLI:**
```bash
go build -ldflags="-s -w" -o ravro_dcrpt ./cmd/cli
```

**Build GUI:**
```bash
go build -ldflags="-s -w" -o ravro_dcrpt_gui ./cmd/gui
```

### Output Files
- `ravro_dcrpt` - CLI application
- `ravro_dcrpt_gui` - GUI application

---

## Cross-Compilation

### Building Windows Binaries from Linux

**Prerequisites:**
```bash
sudo apt-get install -y gcc-mingw-w64-x86-64
```

**Using the provided script:**
```bash
./build_windows.sh
```

This will create:
- `build/ravro_dcrpt.exe` - CLI
- `build/ravro_dcrpt_gui.exe` - GUI (console hidden)

**Manual cross-compilation:**
```bash
# CLI
CGO_ENABLED=1 \
GOOS=windows \
GOARCH=amd64 \
CC=x86_64-w64-mingw32-gcc \
CXX=x86_64-w64-mingw32-g++ \
go build -ldflags="-s -w" -o build/ravro_dcrpt.exe ./cmd/cli

# GUI
CGO_ENABLED=1 \
GOOS=windows \
GOARCH=amd64 \
CC=x86_64-w64-mingw32-gcc \
CXX=x86_64-w64-mingw32-g++ \
go build -ldflags="-s -w -H windowsgui" -o build/ravro_dcrpt_gui.exe ./cmd/gui
```

---

## Using Makefile

The project includes a Makefile for convenient building:

### Available Commands

```bash
# Build both CLI and GUI
make build

# Build only CLI
make build-cli

# Build only GUI
make build-gui

# Build and run CLI
make run-cli

# Build and run GUI
make run-gui

# Clean build artifacts
make clean

# Install dependencies
make deps

# Show all available commands
make help
```

### Windows Note

On Windows, if you have `make` installed (via MinGW, Cygwin, or WSL), the Makefile will automatically detect Windows and apply the `-H windowsgui` flag for GUI builds.

---

## Troubleshooting

### Windows

**Problem:** `cannot find -lssl` or `cannot find -lcrypto`

**Solution:** Ensure OpenSSL is installed correctly:
```powershell
# Check if OpenSSL is installed
Test-Path C:\OpenSSL-Win64\lib\VC\x64\MD\libssl.lib
Test-Path C:\OpenSSL-Win64\lib\VC\x64\MD\libcrypto.lib
```

**Problem:** Console window appears with GUI application

**Solution:** Make sure you use the `-H windowsgui` flag:
```powershell
go build -ldflags="-s -w -H windowsgui" -o ravro_dcrpt_gui.exe ./cmd/gui
```

**Problem:** `wkhtmltox.dll not found` at runtime

**Solution:** Either:
1. Add `C:\wkhtmltox\bin` to your PATH
2. Copy `wkhtmltox.dll` to the same directory as the executable

---

### Linux

**Problem:** `undefined reference to OpenSSL functions`

**Solution:** Install OpenSSL development files:
```bash
# Ubuntu/Debian
sudo apt-get install libssl-dev

# Fedora/RHEL
sudo dnf install openssl-devel
```

**Problem:** `wkhtmltopdf: command not found`

**Solution:** Install wkhtmltopdf:
```bash
sudo apt-get install wkhtmltopdf
```

---

### macOS

**Problem:** `ld: library not found for -lssl`

**Solution:** Set OpenSSL paths correctly:
```bash
export PKG_CONFIG_PATH=$(brew --prefix openssl)/lib/pkgconfig
export CGO_CFLAGS="-I$(brew --prefix openssl)/include"
export CGO_LDFLAGS="-L$(brew --prefix openssl)/lib"
```

**Problem:** wkhtmltopdf not working on macOS

**Solution:** Make sure wkhtmltopdf is installed and in PATH:
```bash
brew install wkhtmltopdf
which wkhtmltopdf  # Should show the path
```

---

## Build Flags Explained

- `-ldflags="-s -w"`: Reduces binary size by stripping debug information
  - `-s`: Omit symbol table
  - `-w`: Omit DWARF symbol table
- `-H windowsgui`: (Windows only) Creates a GUI subsystem binary without console window
- `-o <filename>`: Specifies output filename

---

## Distribution

After building, distribute these files together:

**Windows:**
- `ravro_dcrpt.exe` or `ravro_dcrpt_gui.exe`
- `libcrypto-3-x64.dll` (from `C:\OpenSSL-Win64\bin`)
- `libssl-3-x64.dll` (from `C:\OpenSSL-Win64\bin`)
- `wkhtmltox.dll` (from `C:\wkhtmltox\bin`)

**Linux:**
- `ravro_dcrpt` or `ravro_dcrpt_gui`
- Users need to install: `libssl` and `wkhtmltopdf`

**macOS:**
- `ravro_dcrpt` or `ravro_dcrpt_gui`
- Users need to install: `brew install openssl wkhtmltopdf`

---

## Additional Resources

- [Go Build Documentation](https://golang.org/cmd/go/#hdr-Compile_packages_and_dependencies)
- [CGO Documentation](https://golang.org/cmd/cgo/)
- [Fyne Cross-Compilation](https://developer.fyne.io/started/cross-compiling)

---

**Need help?** Open an issue at [GitHub Issues](https://github.com/ravro-ir/ravro_dcrpt/issues)

