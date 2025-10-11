# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [v2.0.2] - 2025-10-11

### Fixed
- **Windows CI/CD OpenSSL compilation errors**
  - Resolved `openssl/bio.h: No such file or directory` issue
  - Fixed OpenSSL installation with fallback mechanisms
  - Added support for multiple OpenSSL versions (3.3.2, 3.3.1, 3.3.0, 3.2.0, 3.1.0)
  - Implemented automatic junction creation for non-standard OpenSSL paths

### Improved
- **Multi-layer OpenSSL installation strategy**
  - Direct download from slproweb.com as primary method
  - Fallback to Chocolatey if direct download fails
  - Automatic PATH configuration
- **Enhanced error handling and debugging**
  - Better error messages with visual indicators (✅, ⚠️, ❌)
  - Automatic detection of lib directory structure (lib/VC/x64/MD, lib/VC, lib)
  - Verified OpenSSL headers and libraries before build

### Documentation
- **Updated Windows installation documentation**
  - Enhanced `install-windows.ps1` with intelligent OpenSSL installation
  - Added comprehensive troubleshooting section in INSTALL.md
  - Improved README.md Windows installation instructions
  - Added manual installation alternatives

---

## [v2.0.1] - 2025-02-08

### Added
- Enhanced Persian (Farsi) date formatting in PDFs
- Improved error messages for better debugging

### Fixed
- Minor bug fixes in PDF generation
- Improved stability on Windows

---

## [v2.0.0] - 2025-02-02

### 🏗️ Architecture Overhaul
- ✨ **Complete rewrite** using Clean Architecture pattern
- 📦 **Modular design** with clear separation of concerns
- 🔌 **Port & Adapter pattern** for better testability
- 🧹 **Removed legacy code** (23 old files cleaned up)
- 📁 **New project structure**: `cmd/`, `internal/`, `pkg/`

### 🎨 GUI Improvements
- 🖥️ **Added full-featured GUI application** with Fyne framework
- 📂 **Large file browser dialogs** (1000×700) for easier navigation
- ✅ **Real-time validation** and status updates
- 📊 **Live processing logs** and progress tracking
- 🎯 **Window size optimization** (800×600 main window)

### 📄 PDF Generation Enhancements
- 📅 **Automatic date conversion** - Gregorian to Persian (Shamsi/Jalali)
- 💰 **Formatted amounts** with thousand separators (e.g., 10,500,000 ریال)
- 📎 **Attachment tables** - Files displayed in organized tables with types
- 🎨 **Modern styling** - Professional design with Vazirmatn font
- 🌐 **Google Fonts integration** - Beautiful Persian typography
- ⚖️ **Conditional sections** - Judge info only shown when available
- 🔧 **Improved wkhtmltopdf integration** with external resource loading

### 🔐 Crypto & Security
- 🔧 **Fixed OpenSSL Windows CGO integration**
- 📝 **Added proper include headers** for SSL initialization
- 🔑 **Improved key validation** and error handling
- 🛡️ **PKCS7 decryption** with native performance

### 📚 Documentation
- 📝 **Comprehensive BUILD.md** documentation
- 🔨 **Platform-specific build instructions**
- 🌍 **Cross-compilation guide**
- 🐛 **Troubleshooting section**
- 🏗️ **Improved Makefile** with clear targets

### 🐛 Bug Fixes
- ✅ Fixed attachment decryption (all `.ravro` files now processed)
- ✅ Fixed JSON field mapping (camelCase vs snake_case)
- ✅ Fixed report directory detection (direct folder processing)
- ✅ Removed unnecessary debug HTML files
- ✅ Corrected PDF generation with proper data population

### 🛠️ Developer Experience
- 🚀 **Faster build times** with optimized dependencies
- 📦 **`go mod tidy`** for clean dependency management
- 🎯 **Better error messages** and logging
- 🔍 **Improved code organization** and readability

---

## [v1.0.4] - 2024-12-15

### Changed
- Use CGO for OpenSSL and wkhtmltopdf integration
- Improved native library bindings

---

## [v1.0.3] - 2024-11-20

### Added
- Multi-zip file decryption support
- Enhanced path handling

### Improved
- Improved key selection process
- Comprehensive error handling
- Code refactoring

---

## [v1.0.2] - 2024-10-05

### Added
- Added logging capabilities
- Implemented loading spinner
- Added update functionality
- JSON conversion support
- Project packaging

### Improved
- Improved PDF generation performance

### Fixed
- Bug fixes

---

## [v1.0.1] - 2024-09-10

### Added
- Initial release
- Basic PKCS7 decryption
- PDF report generation
- CLI interface

---

[v2.0.2]: https://github.com/ravro-ir/ravro_dcrpt/releases/tag/v2.0.2
[v2.0.1]: https://github.com/ravro-ir/ravro_dcrpt/releases/tag/v2.0.1
[v2.0.0]: https://github.com/ravro-ir/ravro_dcrpt/releases/tag/v2.0.0
[v1.0.4]: https://github.com/ravro-ir/ravro_dcrpt/releases/tag/v1.0.4
[v1.0.3]: https://github.com/ravro-ir/ravro_dcrpt/releases/tag/v1.0.3
[v1.0.2]: https://github.com/ravro-ir/ravro_dcrpt/releases/tag/v1.0.2
[v1.0.1]: https://github.com/ravro-ir/ravro_dcrpt/releases/tag/v1.0.1

