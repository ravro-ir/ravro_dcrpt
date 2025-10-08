# ravro_dcrpt

A versatile Go-based tool for decrypting and converting Ravro platform bug bounty reports to PDF.

## 🚀 Introduction

`ravro_dcrpt` is a cross-platform utility designed to decrypt and process reports submitted by hunters on the Ravro platform, embodying the "write once, run anywhere" philosophy of Go.

## ⚡ Quick Start

1. **Install prerequisites** using the autoinstall scripts (see Installation section)
2. **Download** the latest release for your platform from [Releases](https://github.com/ravro-ir/ravro_dcrpt/releases)
3. **For GUI users**: Double-click `ravro_dcrpt_gui` and use the visual interface
4. **For CLI users**: Run `./ravro_dcrpt -init` to set up directories, then `./ravro_dcrpt` for interactive mode

## ✨ Features

### Core Features
- 🔐 **PKCS7 Decryption** using OpenSSL (CGO-based for native performance)
- 📄 **PDF Generation** with wkhtmltopdf (beautiful, styled reports)
- 🗝️ **Multi-key support** for different organizations
- 🔄 **Batch processing** - multiple reports in one go
- 📦 **ZIP extraction** - automatic handling of compressed reports
- 🖥️ **Cross-platform** - Windows, Linux, macOS

### User Interface
- 🎨 **GUI Application** - User-friendly Fyne-based interface
  - Large file browser dialogs (1000×700) for easy navigation
  - Real-time processing logs and status
  - Directory initialization wizard
  - Key validation before processing
- 💻 **CLI Application** - Perfect for automation and scripting
  - Interactive and command-line modes
  - JSON export support
  - Colored output and progress indicators

### PDF Features
- 📅 **Persian Date Conversion** - Automatic Gregorian to Shamsi (Jalali) conversion
- 💰 **Formatted Amounts** - Thousand separators for rewards (e.g., 10,500,000 ریال)
- 📎 **Attachment Tables** - Beautiful tables showing attachment files with types
- 🎨 **Modern Design** - Clean, professional styling with Vazirmatn font
- 🌐 **RTL Support** - Full right-to-left layout for Persian content
- ⚖️ **Conditional Sections** - Judge information shown only when available

## 🛠️ Installation

### Prerequisites Installation

Before using ravro_dcrpt, you need to install the required dependencies on your system.

#### Linux (Ubuntu/Debian)
```bash
wget https://raw.githubusercontent.com/ravro-ir/ravro_dcrpt/refs/heads/main/autoinstall/linux.sh -O - | sh
```

This script will install:
- Build tools (build-essential, checkinstall, zlib1g-dev)
- OpenSSL development libraries
- wkhtmltopdf for PDF generation
- Additional dependencies (unzip, wget, xfonts-75dpi)

#### macOS
```bash
wget https://raw.githubusercontent.com/ravro-ir/ravro_dcrpt/refs/heads/main/autoinstall/darwin.sh -O - | sh
```

This script will install:
- Homebrew dependencies (if Homebrew is available)
- OpenSSL, unzip, wget, pkg-config
- Required build tools

**Note:** Make sure you have Homebrew installed first. Visit [https://brew.sh](https://brew.sh) for installation instructions.

### Download Application

After installing prerequisites, download the latest release for your platform from [Releases](https://github.com/ravro-ir/ravro_dcrpt/releases)

## 📂 Project Structure

### Clean Architecture Layout

```
ravro_dcrpt/
├── cmd/                          # Application entry points
│   ├── cli/main.go              # CLI application
│   └── gui/main.go              # GUI application
├── internal/                     # Internal packages
│   ├── adapters/                # Interface implementations
│   │   ├── crypto/              # PKCS7 decryption (OpenSSL CGO)
│   │   ├── pdfgen/              # PDF generation (wkhtmltopdf)
│   │   └── storage/             # File system operations
│   ├── core/                    # Business logic
│   │   ├── decrypt/             # Decryption service
│   │   └── report/              # Report processing service
│   └── ports/                   # Interface definitions
├── pkg/                         # Public packages
│   └── models/                  # Data models
├── lib/                         # Native libraries (DLLs)
├── BUILD.md                     # Comprehensive build guide
└── README.md                    # This file
```

### Working Directory Structure

```
.
├── decrypt/                     # Output directory
│   └── ir2025-02-27-0055/
│       ├── company__id__hunter.pdf    # Generated PDF report
│       ├── report/              # Decrypted report data
│       └── amendment-1/         # Decrypted attachments
├── encrypt/                     # Input directory
│   ├── ir2025-02-27-0055.zip    # Compressed report
│   └── ir2025-02-27-0055/       # Or extracted report folder
│       ├── report/
│       │   └── data.ravro       # Encrypted report
│       ├── judgment-1/
│       │   └── data.ravro       # Encrypted judgment
│       └── amendment-1/
│           └── screenshot.ravro # Encrypted attachment
└── key/                         # Private keys directory
    └── COMPANY-PRIVATEKEY.txt   # Private key file
```

## 💻 Usage

### GUI Application (Recommended for Desktop Users)

Simply double-click `ravro_dcrpt_gui` to launch the graphical interface.

**Features:**
- 📁 **Large file browsers** (1000×700) - Browse and select directories/files easily
- 🔑 **Key validation** - Validate private key before processing
- 📊 **Live logs** - Real-time processing status and progress
- ✅ **Directory initialization** - One-click setup of required folders
- 🎯 **Visual feedback** - Clear status messages and error handling

**Steps:**
1. Click "📁 Initialize Directories" to create `encrypt/`, `decrypt/`, and `key/` folders (first time only)
2. Click "Browse" next to each field to select:
   - **Input Directory**: Your encrypted reports folder or specific report
   - **Output Directory**: Where PDFs will be saved
   - **Private Key**: Your `.pem` or `.txt` key file
3. (Optional) Click "🔍 Validate Key" to verify your key
4. Click "🚀 Start Processing" to decrypt and generate PDFs
5. Check the log area for detailed progress and any errors

### CLI Application (For Automation & Scripting)

#### Interactive Mode
```bash
$ ./ravro_dcrpt -init        # Initialize directories
$ ./ravro_dcrpt              # Run in interactive mode
```

#### Command-Line Mode
```bash
# Process a single report
$ ./ravro_dcrpt -in=encrypt/report.zip -out=decrypt -key=key/private.pem

# Process multiple reports in a directory
$ ./ravro_dcrpt -in=encrypt -out=decrypt -key=key

# Process with JSON export
$ ./ravro_dcrpt -in=encrypt -out=decrypt -key=key -json
```

#### Additional Commands
```bash
# Initialize directories (create encrypt/, decrypt/, key/)
$ ./ravro_dcrpt -init

# View version
$ ./ravro_dcrpt -version

# Get help
$ ./ravro_dcrpt -help
```

### 📄 Generated PDF Features

The generated PDF reports include:

**✅ Report Information:**
- Report ID and submission date (Persian calendar)
- Hunter username and target company
- Activity date range (Persian dates)
- Current status and target details
- IP addresses and URLs

**✅ CVSS Scoring:**
- Hunter's CVSS vector and score
- Judge's CVSS evaluation (if available)
- Severity ratings with color coding

**✅ Vulnerability Details:**
- Scenario description (with Markdown support)
- Proof of Concept with full details
- Technical description

**✅ Attachments:**
- Organized in a table format
- File numbering, names, and types
- Automatic file type detection (Image/PDF/File)

**✅ Judge Information (when available):**
- Reward amount with thousand separators (e.g., 10,500,000 ریال)
- Judge's comments and recommendations
- Vulnerability definition and fix suggestions
- Review date (Persian calendar)

**✅ Styling:**
- Clean, professional design
- Persian (Farsi) Vazirmatn font
- Right-to-left (RTL) layout
- Color-coded severity badges
- Responsive table layouts

## 🛠️ Technologies & Dependencies

This project is built with:

**Core:**
- **Go 1.21+** - Primary programming language
- **CGO** - C bindings for native library integration

**Cryptography:**
- **OpenSSL** - PKCS7 encryption/decryption via CGO
- **go.mozilla.org/pkcs7** - Pure Go PKCS7 support (fallback)

**PDF Generation:**
- **wkhtmltopdf** - HTML to PDF conversion with WebKit rendering
- **html/template** - Go's built-in template engine

**GUI Framework:**
- **Fyne v2** - Cross-platform GUI toolkit (Pure Go)

**CLI Framework:**
- **Cobra** - Command-line interface framework

**Utilities:**
- **go-persian-calendar** - Gregorian to Shamsi (Jalali) date conversion
- **archive/zip** - ZIP file extraction

**Fonts:**
- **Vazirmatn** - Modern Persian font from Google Fonts

## 🔨 Building from Source

For comprehensive build instructions including cross-compilation and troubleshooting, see **[BUILD.md](BUILD.md)**.

### Quick Build

**Manually**
```bash
cd /home/raminfp/GolandProjects/ravro_dcrpt && openssl smime -decrypt -inform DER -in encrypt/r2025-02-27-xxxx/report/data.ravro -inkey "key/PRIVATETKEY-20250225-1.txt" -out /tmp/decrypted.json
```

**Linux/macOS:**
```bash
git clone https://github.com/ravro-ir/ravro_dcrpt.git
cd ravro_dcrpt
make build
```


**📘 See [BUILD.md](BUILD.md) for:**
- Prerequisites and dependencies
- Platform-specific instructions
- Cross-compilation guide
- Troubleshooting common issues

### 🍎 macOS Special Notes

Due to recent changes in Homebrew, `wkhtmltopdf` has been deprecated. Use our special installation script:

```bash
# Install wkhtmltopdf for macOS
./install_wkhtmltopdf_macos.sh

# Then build the project
./build_macos.sh
```

## 🐧 Arch Linux Installation
```bash
git clone https://aur.archlinux.org/ravro_dcrpt-git.git
cd ravro_dcrpt-git
makepkg -sri
```

## 📋 Changelog

### v2.0.0 - Major Rewrite (2025-02-02)

#### 🏗️ **Architecture Overhaul**
- ✨ Complete rewrite using **Clean Architecture** pattern
- 📦 Modular design with clear separation of concerns
- 🔌 Port & Adapter pattern for better testability
- 🧹 Removed legacy code (23 old files cleaned up)
- 📁 New project structure: `cmd/`, `internal/`, `pkg/`

#### 🎨 **GUI Improvements**
- 🖥️ Added full-featured GUI application with Fyne framework
- 📂 Large file browser dialogs (1000×700) for easier navigation
- ✅ Real-time validation and status updates
- 📊 Live processing logs and progress tracking
- 🎯 Window size optimization (800×600 main window)

#### 📄 **PDF Generation Enhancements**
- 📅 **Automatic date conversion** - Gregorian to Persian (Shamsi/Jalali)
- 💰 **Formatted amounts** with thousand separators (e.g., 10,500,000 ریال)
- 📎 **Attachment tables** - Files displayed in organized tables with types
- 🎨 **Modern styling** - Professional design with Vazirmatn font
- 🌐 **Google Fonts integration** - Beautiful Persian typography
- ⚖️ **Conditional sections** - Judge info only shown when available
- 🔧 **Improved wkhtmltopdf integration** with external resource loading

#### 🔐 **Crypto & Security**
- 🔧 Fixed OpenSSL Windows CGO integration
- 📝 Added proper include headers for SSL initialization
- 🔑 Improved key validation and error handling
- 🛡️ PKCS7 decryption with native performance

#### 📚 **Documentation**
- 📝 Comprehensive [BUILD.md](BUILD.md) documentation
- 🔨 Platform-specific build instructions
- 🌍 Cross-compilation guide
- 🐛 Troubleshooting section
- 🏗️ Improved Makefile with clear targets

#### 🐛 **Bug Fixes**
- ✅ Fixed attachment decryption (all `.ravro` files now processed)
- ✅ Fixed JSON field mapping (camelCase vs snake_case)
- ✅ Fixed report directory detection (direct folder processing)
- ✅ Removed unnecessary debug HTML files
- ✅ Corrected PDF generation with proper data population

#### 🛠️ **Developer Experience**
- 🚀 Faster build times with optimized dependencies
- 📦 `go mod tidy` for clean dependency management
- 🎯 Better error messages and logging
- 🔍 Improved code organization and readability

### v1.0.4
- Use CGO for OpenSSL and wkhtmltopdf

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

