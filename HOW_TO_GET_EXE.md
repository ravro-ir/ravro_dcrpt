# 🪟 چطور فایل .exe برای Windows بگیریم؟

## وضعیت فعلی

❌ **الان فایل `.exe` نداریم!**

فایل‌های فعلی فقط برای لینوکس هستند:
- `build/ravro_dcrpt` (6.5 MB) - Linux ELF
- `build/ravro_dcrpt_gui` (22 MB) - Linux ELF

---

## ✅ سه راه برای گرفتن فایل `.exe`

### 🥇 راه اول: GitHub Actions (بهترین - خودکار)

**مزایا:**
- ✅ خودکار
- ✅ برای همه پلتفرم‌ها (Windows/Linux/macOS)
- ✅ قابل اطمینان
- ✅ فایل‌ها در GitHub Releases

**مراحل:**
```bash
# 1. مطمئن شو همه چی commit شده
git status

# 2. اگر تغییری هست، commit کن
git add .
git commit -m "Add Windows build support via GitHub Actions"

# 3. Push کن
git push

# 4. یک tag بزن برای Release
git tag v1.0.0

# 5. Tag را push کن
git push origin v1.0.0

# 6. صبر کن تا GitHub Actions build کنه (5-10 دقیقه)
# 7. برو به: https://github.com/YOUR_USERNAME/ravro_dcrpt/releases
# 8. فایل‌های .exe را دانلود کن!
```

**نتیجه:**
- `ravro_dcrpt-windows-amd64.exe` ✅
- `ravro_dcrpt_gui-windows-amd64.exe` ✅

---

### 🥈 راه دوم: Build در Windows (ساده‌ترین)

**مزایا:**
- ✅ خیلی ساده
- ✅ سریع
- ✅ بدون نیاز به setup پیچیده

**مراحل:**

#### 1. نصب Dependencies

```powershell
# در Windows PowerShell (به عنوان Administrator):

# 1. نصب Chocolatey (اگر نداری)
Set-ExecutionPolicy Bypass -Scope Process -Force
[System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072
iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))

# 2. نصب Go
choco install golang -y

# 3. نصب Git
choco install git -y

# 4. نصب MinGW
choco install mingw -y

# 5. نصب OpenSSL
choco install openssl -y

# 6. نصب wkhtmltopdf
choco install wkhtmltopdf -y

# 7. Restart PowerShell
```

#### 2. Clone پروژه

```powershell
cd C:\
git clone https://github.com/YOUR_USERNAME/ravro_dcrpt
cd ravro_dcrpt
```

#### 3. Build

```powershell
# Set environment variables
$env:PATH="C:/OpenSSL-Win64/bin;C:/wkhtmltox/bin;$env:PATH"
$env:CGO_CFLAGS="-IC:/OpenSSL-Win64/include -IC:/wkhtmltox/include"
$env:CGO_LDFLAGS="-LC:/OpenSSL-Win64/lib/VC/x64/MD -LC:/wkhtmltox/lib -L/C:/wkhtmltox/bin -lssl -lcrypto -lws2_32 -lcrypt32 -lwkhtmltox"

# Build CLI
go build -ldflags="-s -w" -o ravro_dcrpt.exe .\cmd\cli

# Build GUI (without console window)
go build -ldflags="-s -w -H windowsgui" -o ravro_dcrpt_gui.exe .\cmd\gui
```

**توضیحات:**
- فلگ `-s -w` حجم فایل را کاهش می‌دهد
- فلگ `-H windowsgui` پنجره CMD را برای برنامه GUI مخفی می‌کند

**نتیجه:**
- `ravro_dcrpt.exe` ✅
- `ravro_dcrpt_gui.exe` ✅

---

### 🥉 راه سوم: Cross-Compile از لینوکس (پیچیده)

**هشدار:** ⚠️ زمان‌بر (30+ دقیقه) و پیچیده!

**مراحل:**

```bash
# 1. نصب MinGW
sudo apt-get update
sudo apt-get install -y mingw-w64 gcc-mingw-w64-x86-64

# 2. دانلود و Build OpenSSL برای Windows (15-20 دقیقه!)
cd /tmp
wget https://www.openssl.org/source/openssl-1.1.1w.tar.gz
tar xzf openssl-1.1.1w.tar.gz
cd openssl-1.1.1w

# Configure برای MinGW
./Configure mingw64 \
    --cross-compile-prefix=x86_64-w64-mingw32- \
    --prefix=/usr/x86_64-w64-mingw32 \
    no-shared

# Build (صبر کن!)
make -j$(nproc)

# Install
sudo make install

# 3. Build CLI
cd /home/raminfp/GolandProjects/ravro_dcrpt
CGO_ENABLED=1 \
GOOS=windows \
GOARCH=amd64 \
CC=x86_64-w64-mingw32-gcc \
PKG_CONFIG_PATH=/usr/x86_64-w64-mingw32/lib/pkgconfig \
go build -o build/ravro_dcrpt.exe ./cmd/cli

# 4. Build GUI (نیاز به fyne-cross)
go install github.com/fyne-io/fyne-cross@latest
fyne-cross windows -arch=amd64 ./cmd/gui
```

---

## 🎯 توصیه من

### اگر GitHub داری:
→ **راه اول** (GitHub Actions) - خودکار و عالی!

### اگر Windows داری:
→ **راه دوم** (Build در Windows) - ساده و سریع!

### اگر فقط لینوکس داری و عجله داری:
→ **راه سوم** (Cross-compile) - اما زمان‌بر!

---

## ❓ سوالات متداول

**Q: چرا از لینوکس نمی‌تونم راحت build بگیرم؟**
A: چون از CGO + OpenSSL استفاده می‌کنیم که نیاز به native library داره.

**Q: بهترین روش چیه؟**
A: GitHub Actions! یک بار setup می‌کنی، بعد همیشه خودکار build می‌گیره.

**Q: اگر GitHub نداشته باشم؟**
A: در Windows build بگیر - خیلی ساده‌تره!

---

**موفق باشید! 🚀**
