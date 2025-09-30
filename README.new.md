# 🔐 Ravro Decryption Tool v2.0

A **Pure Go** cross-platform tool for decrypting and converting Ravro platform bug bounty reports to PDF with both CLI and GUI interfaces.

## ✨ Key Features

### 🚀 Pure Go Implementation
- ✅ **No CGO dependencies** - truly cross-platform
- ✅ **No external binaries** required (OpenSSL, wkhtmltopdf)
- ✅ **Easy cross-compilation** from any platform to any platform
- ✅ **Single static binary** - no runtime dependencies

### 🎯 Dual Interface
- 🖥️ **CLI** - Command-line interface with rich output
- 🎨 **GUI** - Modern graphical interface built with Fyne

### 🔒 Security
- 🔐 PKCS7/SMIME decryption using pure Go cryptography
- 🔑 Support for RSA private keys (PEM format)
- ✅ Key validation before processing

### 📄 PDF Generation
- 📝 Beautiful PDF reports with Persian (RTL) support
- 🎨 Clean, professional layout
- 📋 Support for attachments and metadata

### 🌍 Cross-Platform
- 🐧 Linux (amd64)
- 🪟 Windows (amd64)
- 🍎 macOS (amd64, arm64/M1)

## 📦 Installation

### Quick Install (Linux/macOS)

```bash
# Download and extract
curl -L https://github.com/ravro-ir/ravro_dcrpt/releases/latest/download/ravro_dcrpt-2.0.0-linux-amd64.tar.gz | tar xz

# Move to PATH
sudo mv ravro_dcrpt /usr/local/bin/

# Run
ravro_dcrpt --help
```

### Quick Install (Windows)

Download from [Releases](https://github.com/ravro-ir/ravro_dcrpt/releases) and extract to your preferred location.

### Build from Source

#### Prerequisites
- Go 1.21 or higher
- No other dependencies!

#### Build Commands

```bash
# Clone repository
git clone https://github.com/ravro-ir/ravro_dcrpt.git
cd ravro_dcrpt

# Build for current platform
make build

# Or use build script
./build.sh

# Build for all platforms
make build-all
# Or
./build.sh --all

# Build for specific platform
make build-linux
make build-windows
make build-darwin
# Or
./build.sh --linux
./build.sh --windows
./build.sh --darwin

# Build only CLI
./build.sh --cli-only

# Build only GUI
./build.sh --gui-only
```

#### Cross-Compilation Examples

```bash
# From Linux, build for Windows
GOOS=windows GOARCH=amd64 go build -o ravro_dcrpt.exe ./cmd/cli

# From macOS, build for Linux
GOOS=linux GOARCH=amd64 go build -o ravro_dcrpt ./cmd/cli

# From Windows, build for macOS
set GOOS=darwin
set GOARCH=arm64
go build -o ravro_dcrpt ./cmd/cli
```

## 🎮 Usage

### CLI Interface

#### Initialize Directories
```bash
ravro_dcrpt --init
```

This creates:
```
.
├── encrypt/    # Place encrypted reports here
├── decrypt/    # Decrypted PDFs will be saved here
└── key/        # Place your private key here
```

#### Basic Usage
```bash
# Process all reports in default directories
ravro_dcrpt

# Specify custom paths
ravro_dcrpt --in=/path/to/reports --out=/path/to/output --key=/path/to/key.pem

# Export as JSON (in addition to PDF)
ravro_dcrpt --json
```

#### Help
```bash
ravro_dcrpt --help
```

### GUI Interface

Simply run the GUI application:

```bash
# Linux/macOS
./ravro_dcrpt-gui

# Windows
ravro_dcrpt-gui.exe
```

Features:
- 📁 Browse and select directories
- 🔑 Select private key file
- ✅ Validate key before processing
- 📊 Real-time progress and logs
- 🎯 User-friendly interface with Persian support

## 🏗️ Architecture

This project follows **Clean Architecture** principles:

```
ravro_dcrpt/
├── cmd/
│   ├── cli/              # CLI application
│   └── gui/              # GUI application
├── internal/
│   ├── core/             # Business logic
│   │   ├── decrypt/      # Decryption service
│   │   ├── pdf/          # PDF generation service
│   │   └── report/       # Report processing
│   ├── adapters/         # External adapters
│   │   ├── crypto/       # PKCS7 implementation (Pure Go)
│   │   ├── pdfgen/       # PDF generation (Pure Go)
│   │   └── storage/      # File system operations
│   └── ports/            # Interfaces
├── pkg/                  # Public libraries
│   └── models/           # Data models
└── ui/                   # GUI resources
```

### Technology Stack

| Component | Technology | Why? |
|-----------|-----------|------|
| Crypto | [go.mozilla.org/pkcs7](https://github.com/mozilla/pkcs7) | Pure Go PKCS7 (no OpenSSL) |
| PDF Generation | [maroto](https://github.com/johnfercher/maroto) | Pure Go PDF with RTL support |
| CLI Framework | [cobra](https://github.com/spf13/cobra) | Industry standard CLI |
| GUI Framework | [fyne](https://fyne.io/) | Pure Go, cross-platform GUI |
| Persian Calendar | [go-persian-calendar](https://github.com/yaa110/go-persian-calendar) | Shamsi date support |

## 🔄 Migration from v1.x

### Key Differences

| Feature | v1.x | v2.0 |
|---------|------|------|
| OpenSSL | ❌ CGO dependency | ✅ Pure Go |
| wkhtmltopdf | ❌ External binary | ✅ Pure Go |
| Cross-compile | ❌ Complex | ✅ Simple |
| GUI | ❌ No | ✅ Yes |
| Build from Linux for Windows | ❌ Impossible | ✅ Easy |

### Migration Steps

1. **No changes needed** for encrypted report files
2. **Same key format** - use your existing private keys
3. **Same directory structure** - encrypt/decrypt/key folders work the same
4. **Enhanced features** - now with GUI and easier builds!

## 📋 Project Structure

```
.
├── cmd/
│   ├── cli/main.go           # CLI entry point
│   └── gui/main.go           # GUI entry point
├── internal/
│   ├── core/
│   │   ├── decrypt/service.go    # Decryption logic
│   │   └── report/service.go     # Report processing
│   ├── adapters/
│   │   ├── crypto/pkcs7.go       # PKCS7 implementation
│   │   ├── pdfgen/maroto.go      # PDF generator
│   │   └── storage/filesystem.go # File operations
│   └── ports/
│       ├── crypto.go             # Crypto interface
│       ├── pdf.go                # PDF interface
│       └── storage.go            # Storage interface
├── pkg/
│   └── models/report.go          # Data models
├── Makefile                      # Build automation
├── build.sh                      # Cross-platform build script
├── go.mod                        # Go dependencies
└── README.md                     # This file
```

## 🧪 Testing

```bash
# Run tests
make test

# Run tests with coverage
make test-coverage

# Format code
make fmt

# Run linter
make lint

# Run go vet
make vet
```

## 🐳 Docker Support

```bash
# Build Docker image
make docker-build

# Run in Docker
docker run -v $(pwd)/encrypt:/app/encrypt \
           -v $(pwd)/decrypt:/app/decrypt \
           -v $(pwd)/key:/app/key \
           ravro_dcrpt:latest
```

## 📝 Example Workflow

1. **Initialize directories**
   ```bash
   ravro_dcrpt --init
   ```

2. **Place files**
   - Copy encrypted reports (`.zip` or `.ravro`) to `encrypt/` directory
   - Copy your private key to `key/` directory

3. **Process reports**
   ```bash
   # CLI
   ravro_dcrpt
   
   # Or GUI
   ./ravro_dcrpt-gui
   ```

4. **Get results**
   - PDF files will be in `decrypt/` directory
   - Each report gets its own subdirectory

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

GNU General Public License v3.0

## 👥 Authors

### v2.0 Rewrite
- Ramin Farajpour Cami - Complete rewrite with Pure Go implementation

### Original v1.x
- Ravro Development Team (RDT)

## 📞 Contact

- Email: ramin.blackhat@gmail.com
- Alternate: farajpour@ravro.ir
- GitHub: [ravro-ir/ravro_dcrpt](https://github.com/ravro-ir/ravro_dcrpt)

## 🙏 Acknowledgments

- Mozilla for the excellent PKCS7 library
- Fyne team for the amazing GUI framework
- Maroto team for the PDF generation library
- Original Ravro Development Team

## 📊 Comparison with v1.x

### Build Size
- **v1.x**: ~50MB (with CGO dependencies)
- **v2.0**: ~15MB (pure Go, stripped)

### Build Time
- **v1.x**: Complex setup with OpenSSL and wkhtmltopdf
- **v2.0**: Simple `go build` - done! ✨

### Cross-Compilation
- **v1.x**: Nearly impossible (CGO for Windows from Linux)
- **v2.0**: Single command: `GOOS=windows go build` 🚀

## 🗺️ Roadmap

- [x] Pure Go crypto implementation
- [x] Pure Go PDF generation
- [x] CLI interface
- [x] GUI interface
- [x] Cross-platform builds
- [x] Persian/RTL support
- [ ] Batch processing improvements
- [ ] Web interface
- [ ] Report templates customization
- [ ] Digital signatures verification
- [ ] Multi-language support

---

**Made with ❤️ and Go**


## استفاده از HTML Template (پیشرفته)

اگر می‌خواهید از template HTML قبلی استفاده کنید، می‌توانید با یک library مانند `wkhtmltopdf` آن را به PDF تبدیل کنید:

### مرحله 1: نصب wkhtmltopdf

```bash
# Ubuntu/Debian
sudo apt-get install wkhtmltopdf

# macOS
brew install wkhtmltopdf

# Windows
# دانلود از https://wkhtmltopdf.org/downloads.html
```

### مرحله 2: تولید HTML

HTML template در `internal/adapters/pdfgen/html_template.go` موجود است و در آینده به صورت خودکار render خواهد شد.

### مرحله 3: تبدیل به PDF

```bash
wkhtmltopdf --encoding utf-8 report.html report.pdf
```

**توجه:** در نسخه فعلی، PDF به صورت خودکار با کتابخانه Maroto تولید می‌شود که Pure Go است و نیازی به dependency خارجی ندارد.

