# 🔨 راهنمای Cross-Compilation برای Windows

## ⚠️ مشکل فعلی

برای build کردن برای Windows از لینوکس، به موارد زیر نیاز است:

### 1. CLI (امکان‌پذیر اما پیچیده)

**مشکلات:**
- نیاز به OpenSSL برای MinGW (cross-compile)
- نیاز به wkhtmltopdf برای MinGW

**راه‌حل‌ها:**

#### روش 1: استفاده از Docker (توصیه می‌شود)
```bash
# استفاده از image آماده
docker run --rm -v "$PWD":/go/src/app \
    -w /go/src/app \
    dockercore/golang-cross:latest \
    sh -c "apt-get update && apt-get install -y mingw-w64 && \
           CGO_ENABLED=1 GOOS=windows GOARCH=amd64 \
           CC=x86_64-w64-mingw32-gcc \
           go build -o ravro_dcrpt.exe ./cmd/cli"
```

#### روش 2: نصب OpenSSL برای MinGW (زمان‌بر)
```bash
# 1. نصب MinGW
sudo apt-get install -y mingw-w64

# 2. دانلود و build OpenSSL برای Windows
cd /tmp
wget https://www.openssl.org/source/openssl-1.1.1w.tar.gz
tar xzf openssl-1.1.1w.tar.gz
cd openssl-1.1.1w

# Configure برای MinGW
./Configure mingw64 --cross-compile-prefix=x86_64-w64-mingw32- \
    --prefix=/usr/x86_64-w64-mingw32 no-shared

# Build (10-20 دقیقه)
make -j$(nproc)
sudo make install

# حالا build CLI
cd /path/to/ravro_dcrpt
CGO_ENABLED=1 GOOS=windows GOARCH=amd64 \
    CC=x86_64-w64-mingw32-gcc \
    PKG_CONFIG_PATH=/usr/x86_64-w64-mingw32/lib/pkgconfig \
    go build -o build/ravro_dcrpt.exe ./cmd/cli
```

#### روش 3: Build در Windows مستقیماً (ساده‌ترین!)
```bash
# در ویندوز:
# 1. نصب Go
# 2. نصب MSYS2 و MinGW
# 3. نصب OpenSSL
# 4. نصب wkhtmltopdf
# 5. Build:
go build -o ravro_dcrpt.exe ./cmd/cli
```

### 2. GUI (پیچیده‌تر)

برای GUI با Fyne، باید از `fyne-cross` استفاده کنید:

```bash
# نصب fyne-cross
go install fyne.io/fyne/v2/cmd/fyne@latest
go install github.com/fyne-io/fyne-cross@latest

# Build برای Windows
fyne-cross windows -arch=amd64 -app-id=ir.ravro.dcrpt ./cmd/gui

# خروجی در: fyne-cross/dist/windows-amd64/
```

## ✅ توصیه نهایی

**بهترین راه:**

### برای کاربران لینوکس:
```bash
# Build CLI و GUI برای لینوکس
go build -o build/ravro_dcrpt ./cmd/cli
go build -o build/ravro_dcrpt_gui ./cmd/gui
```

### برای کاربران Windows:
1. پروژه را در Windows clone کنید
2. مراحل نصب را از `README.md` دنبال کنید
3. Build بگیرید:
```cmd
go build -o ravro_dcrpt.exe .\cmd\cli
go build -o ravro_dcrpt_gui.exe .\cmd\gui
```

### برای توزیع:
- از GitHub Actions/CI استفاده کنید
- یا از Docker Multi-stage build
- یا هر پلتفرم را در OS خودش build کنید

## 📦 GitHub Actions مثال

```yaml
name: Build
on: [push, pull_request]
jobs:
  build:
    strategy:
      matrix:
        os: [ubuntu-latest, windows-latest, macos-latest]
    runs-on: ${{ matrix.os }}
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Install deps (Ubuntu)
        if: runner.os == 'Linux'
        run: |
          sudo apt-get update
          sudo apt-get install -y libssl-dev wkhtmltopdf
      
      - name: Install deps (Windows)
        if: runner.os == 'Windows'
        run: |
          choco install openssl wkhtmltopdf
      
      - name: Build CLI
        run: go build -o ravro_dcrpt ./cmd/cli
      
      - name: Build GUI
        run: go build -o ravro_dcrpt_gui ./cmd/gui
```

## 🎯 نتیجه‌گیری

Cross-compilation برای Go با CGO پیچیده است. بهترین راه:

1. ✅ **لینوکس**: مستقیماً build بگیرید
2. ✅ **Windows**: در Windows build بگیرید  
3. ✅ **macOS**: در macOS build بگیرید
4. ✅ **توزیع**: از CI/CD استفاده کنید

این روش ساده‌تر، سریع‌تر و قابل اطمینان‌تر است! 🚀
