# 🚀 Deployment Guide - Ravro Decryption Tool v2.0

## 📋 Table of Contents
- [Quick Start](#quick-start)
- [CLI Cross-Compilation](#cli-cross-compilation)
- [GUI Cross-Compilation](#gui-cross-compilation)
- [Build Scripts](#build-scripts)
- [Distribution](#distribution)

## ⚡ Quick Start

### Build for Current Platform

```bash
# Clone repository
git clone https://github.com/ravro-ir/ravro_dcrpt.git
cd ravro_dcrpt

# Build CLI and GUI
make build

# Or using build script
./build.sh
```

## 🖥️ CLI Cross-Compilation

CLI is **100% Pure Go** and can be cross-compiled easily from **any platform to any platform**.

### From Linux to All Platforms

```bash
# For Linux
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o build/ravro_dcrpt-linux-amd64 ./cmd/cli

# For Windows
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o build/ravro_dcrpt-windows-amd64.exe ./cmd/cli

# For macOS Intel
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o build/ravro_dcrpt-darwin-amd64 ./cmd/cli

# For macOS Apple Silicon (M1/M2)
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o build/ravro_dcrpt-darwin-arm64 ./cmd/cli
```

### Using Makefile

```bash
# Build for all platforms
make build-all

# Build for specific platform
make build-linux
make build-windows
make build-darwin

# Create distribution packages
make package
```

### Using Build Script

```bash
# Build for all platforms
./build.sh --all

# Build for specific platform
./build.sh --linux
./build.sh --windows
./build.sh --darwin

# Build only CLI
./build.sh --all --cli-only
```

## 🎨 GUI Cross-Compilation

GUI uses Fyne which requires platform-specific dependencies. There are two approaches:

### Approach 1: fyne-cross (Recommended for Cross-Compilation)

```bash
# Install fyne-cross
go install github.com/fyne-io/fyne-cross@latest

# Make sure Docker is running

# Build for Linux
fyne-cross linux -app-id com.ravro.dcrpt -icon ui/icon.png ./cmd/gui

# Build for Windows
fyne-cross windows -app-id com.ravro.dcrpt -icon ui/icon.png ./cmd/gui

# Build for macOS
fyne-cross darwin -app-id com.ravro.dcrpt -icon ui/icon.png ./cmd/gui

# Build for all platforms
fyne-cross linux windows darwin -app-id com.ravro.dcrpt -icon ui/icon.png ./cmd/gui
```

### Approach 2: Native Build (Simpler, but requires building on each platform)

#### On Linux
```bash
# Install development libraries
sudo apt-get install gcc libgl1-mesa-dev xorg-dev  # Debian/Ubuntu
sudo dnf install gcc mesa-libGL-devel libXcursor-devel libXrandr-devel libXinerama-devel libXi-devel  # Fedora

# Build GUI
go build -o build/ravro_dcrpt-gui ./cmd/gui
```

#### On Windows
```powershell
# Install MinGW-w64 for gcc
# Download from: https://www.msys2.org/

# Build GUI
go build -o build/ravro_dcrpt-gui.exe ./cmd/gui
```

#### On macOS
```bash
# Xcode Command Line Tools required
xcode-select --install

# Build GUI
go build -o build/ravro_dcrpt-gui ./cmd/gui
```

## 📦 Build Scripts

### Makefile Targets

```bash
make help           # Show all available commands
make deps           # Install dependencies
make build          # Build for current platform
make build-all      # Build for all platforms
make build-linux    # Build for Linux
make build-windows  # Build for Windows
make build-darwin   # Build for macOS
make package        # Create distribution packages
make test           # Run tests
make clean          # Clean build artifacts
make install        # Install CLI to system
```

### Build Script Options

```bash
./build.sh --help        # Show help
./build.sh               # Build for current platform
./build.sh --all         # Build for all platforms
./build.sh --linux       # Build for Linux only
./build.sh --windows     # Build for Windows only
./build.sh --darwin      # Build for macOS only
./build.sh --cli-only    # Build only CLI
./build.sh --gui-only    # Build only GUI
```

## 📤 Distribution

### Creating Release Packages

```bash
# Using Makefile
make package

# This creates:
# - ravro_dcrpt-2.0.0-linux-amd64.tar.gz
# - ravro_dcrpt-2.0.0-windows-amd64.zip
# - ravro_dcrpt-2.0.0-darwin-amd64.tar.gz
# - ravro_dcrpt-2.0.0-darwin-arm64.tar.gz
```

### Manual Packaging

#### Linux
```bash
cd build/linux
tar -czf ../ravro_dcrpt-2.0.0-linux-amd64.tar.gz ravro_dcrpt ravro_dcrpt-gui
```

#### Windows
```bash
cd build/windows
zip -r ../ravro_dcrpt-2.0.0-windows-amd64.zip ravro_dcrpt.exe ravro_dcrpt-gui.exe
```

#### macOS
```bash
cd build/darwin
tar -czf ../ravro_dcrpt-2.0.0-darwin-amd64.tar.gz ravro_dcrpt-amd64 ravro_dcrpt-gui-amd64
tar -czf ../ravro_dcrpt-2.0.0-darwin-arm64.tar.gz ravro_dcrpt-arm64 ravro_dcrpt-gui-arm64
```

## 🐳 Docker Build

### Using Docker for Consistent Builds

```bash
# Build Docker image
docker build -t ravro_dcrpt-builder .

# Run builds inside Docker
docker run -v $(pwd):/app ravro_dcrpt-builder make build-all
```

### Dockerfile Example

```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .

RUN apk add --no-cache make git
RUN make build-all

FROM alpine:latest
COPY --from=builder /app/build/* /usr/local/bin/
CMD ["ravro_dcrpt", "--help"]
```

## 🔍 Verification

### Verify Builds

```bash
# Check binary size
ls -lh build/

# Test CLI
./build/ravro_dcrpt --help

# Test GUI (requires display)
./build/ravro_dcrpt-gui
```

### Run Tests

```bash
# Unit tests
make test

# With coverage
make test-coverage

# Linting
make lint
```

## 📊 Build Sizes (Approximate)

| Platform | CLI | GUI |
|----------|-----|-----|
| Linux amd64 | ~12MB | ~18MB |
| Windows amd64 | ~12MB | ~18MB |
| macOS amd64 | ~12MB | ~18MB |
| macOS arm64 | ~12MB | ~18MB |

*Sizes after stripping debug symbols with `-ldflags="-s -w"`*

## 🚨 Common Issues

### Issue: `cannot find package`
**Solution**: Run `go mod tidy` to download dependencies

### Issue: GUI won't cross-compile
**Solution**: Use `fyne-cross` or build natively on target platform

### Issue: `gcc not found` on Linux
**Solution**: Install build-essential
```bash
sudo apt-get install build-essential
```

### Issue: GUI doesn't run on Linux
**Solution**: Install required libraries
```bash
sudo apt-get install libgl1-mesa-glx libxi6 libxrandr2 libxcursor1 libxinerama1
```

## 🎯 CI/CD Integration

### GitHub Actions Example

```yaml
name: Build and Release

on:
  push:
    tags:
      - 'v*'

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Build
        run: make build-all
      
      - name: Package
        run: make package
      
      - name: Create Release
        uses: softprops/action-gh-release@v1
        with:
          files: build/*.{tar.gz,zip}
```

## 📝 Best Practices

1. **Always test on target platform** before releasing
2. **Use stripped binaries** (`-ldflags="-s -w"`) for smaller size
3. **Include README** and license in distribution packages
4. **Sign binaries** for macOS and Windows (for production)
5. **Provide checksums** (SHA256) for all releases

## 🔗 Useful Links

- [Go Cross-Compilation](https://go.dev/doc/install/source#environment)
- [Fyne Cross](https://github.com/fyne-io/fyne-cross)
- [GoReleaser](https://goreleaser.com/) - Automated release tool

---

**Need Help?** Open an issue on GitHub or contact: ramin.blackhat@gmail.com

