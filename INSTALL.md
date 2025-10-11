# نصب و راه‌اندازی Ravro Decryption Tool

این راهنما نحوه نصب و اجرای ابزار رمزگشایی Ravro را برای سیستم‌عامل‌های مختلف توضیح می‌دهد.

## 📋 پیش‌نیازها

برای اجرای این نرم‌افزار، باید کتابخانه‌های زیر روی سیستم شما نصب باشند:

### Linux
- **OpenSSL** (libssl3 / libcrypto)
- **X11 Libraries** (libX11, libXcursor, libXrandr, libXinerama, libXi, libXxf86vm)
- **Mesa/OpenGL** (libGL)
- **wkhtmltopdf** (برای تولید PDF)

### macOS
- **Homebrew** (مدیر بسته macOS)
- **OpenSSL@3**
- **wkhtmltopdf** (برای تولید PDF)

### Windows
- **Chocolatey** (مدیر بسته Windows)
- **OpenSSL** (libssl / libcrypto)
- **wkhtmltopdf** (برای تولید PDF)

---

## 🚀 نصب خودکار (توصیه می‌شود)

### Linux

1. دانلود اسکریپت نصب:
```bash
curl -O https://raw.githubusercontent.com/ravro-ir/ravro_dcrpt/main/install-linux.sh
chmod +x install-linux.sh
```

2. اجرای اسکریپت نصب:
```bash
./install-linux.sh
```

این اسکریپت به صورت خودکار تمام dependency های مورد نیاز را برای توزیع‌های زیر نصب می‌کند:
- Ubuntu / Debian / Linux Mint / Pop!_OS
- Fedora / RHEL / CentOS / Rocky / AlmaLinux
- Arch Linux / Manjaro
- openSUSE / SLES

### macOS

1. دانلود اسکریپت نصب:
```bash
curl -O https://raw.githubusercontent.com/ravro-ir/ravro_dcrpt/main/install-macos.sh
chmod +x install-macos.sh
```

2. اجرای اسکریپت نصب:
```bash
./install-macos.sh
```

این اسکریپت به صورت خودکار:
- Homebrew را نصب می‌کند (در صورت نیاز)
- OpenSSL@3 را نصب می‌کند
- wkhtmltopdf را نصب می‌کند

### Windows

1. دانلود اسکریپت نصب:
   - به [صفحه GitHub](https://github.com/ravro-ir/ravro_dcrpt) بروید
   - فایل `install-windows.ps1` را دانلود کنید

2. اجرای اسکریپت نصب (به عنوان Administrator):
```powershell
# Right-click PowerShell → Run as Administrator
Set-ExecutionPolicy Bypass -Scope Process -Force
.\install-windows.ps1
```

این اسکریپت به صورت خودکار:
- Chocolatey را نصب می‌کند (در صورت نیاز)
- OpenSSL را نصب می‌کند
- wkhtmltopdf را نصب می‌کند

**نکته:** حتماً PowerShell را به عنوان Administrator اجرا کنید.

---

## 📦 دانلود نرم‌افزار

از صفحه [Releases](https://github.com/ravro-ir/ravro_dcrpt/releases) آخرین نسخه را دانلود کنید:

### Linux (x86_64)
```bash
wget https://github.com/ravro-ir/ravro_dcrpt/releases/latest/download/ravro_dcrpt-linux-amd64.tar.gz
tar -xzf ravro_dcrpt-linux-amd64.tar.gz
```

### macOS (Intel)
```bash
wget https://github.com/ravro-ir/ravro_dcrpt/releases/latest/download/ravro_dcrpt-darwin-amd64.tar.gz
tar -xzf ravro_dcrpt-darwin-amd64.tar.gz
```

### macOS (Apple Silicon)
```bash
wget https://github.com/ravro-ir/ravro_dcrpt/releases/latest/download/ravro_dcrpt-darwin-arm64.tar.gz
tar -xzf ravro_dcrpt-darwin-arm64.tar.gz
```

### Windows
از صفحه [Releases](https://github.com/ravro-ir/ravro_dcrpt/releases) فایل `ravro_dcrpt-windows-amd64.zip` را دانلود کرده و Extract کنید.

---

## ▶️ اجرای نرم‌افزار

### Linux

```bash
chmod +x ravro_dcrpt_gui
./ravro_dcrpt_gui
```

**نکته:** اطمینان حاصل کنید که display server (X11 یا Wayland) در حال اجراست.

### macOS

```bash
open "Ravro Decryption Tool.app"
```

یا با دابل‌کلیک روی فایل `.app` در Finder.

**هشدار امنیتی macOS:**
در اولین اجرا، macOS ممکن است پیام امنیتی نمایش دهد:

1. به `System Preferences` → `Security & Privacy` بروید
2. در تب `General` روی `Open Anyway` کلیک کنید
3. دوباره برنامه را اجرا کنید

### Windows

```powershell
# Extract کردن فایل zip
Expand-Archive -Path ravro_dcrpt-windows-amd64.zip -DestinationPath ravro_dcrpt

# اجرای برنامه
cd ravro_dcrpt
.\ravro_dcrpt_gui.exe
```

یا دابل‌کلیک روی `ravro_dcrpt_gui.exe` در File Explorer.

**هشدار امنیتی Windows:**
در اولین اجرا، Windows Defender SmartScreen ممکن است هشدار دهد:

1. روی `More info` کلیک کنید
2. روی `Run anyway` کلیک کنید

---

## 🔧 نصب دستی Dependencies

اگر اسکریپت نصب خودکار کار نکرد، می‌توانید به صورت دستی نصب کنید:

### Ubuntu/Debian

```bash
sudo apt-get update
sudo apt-get install -y \
    libgl1-mesa-glx \
    libx11-6 \
    libxcursor1 \
    libxrandr2 \
    libxinerama1 \
    libxi6 \
    libxxf86vm1 \
    libssl3 \
    wkhtmltopdf
```

### Fedora/RHEL

```bash
sudo dnf install -y \
    mesa-libGL \
    libX11 \
    libXcursor \
    libXrandr \
    libXinerama \
    libXi \
    libXxf86vm \
    openssl-libs \
    wkhtmltopdf
```

### Arch Linux

```bash
sudo pacman -S \
    mesa \
    libx11 \
    libxcursor \
    libxrandr \
    libxinerama \
    libxi \
    libxxf86vm \
    openssl \
    wkhtmltopdf
```

### macOS (با Homebrew)

```bash
# نصب Homebrew (در صورت نیاز)
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# نصب dependencies
brew install openssl@3
brew install --cask wkhtmltopdf
```

### Windows (با Chocolatey)

```powershell
# نصب Chocolatey (در صورت نیاز) - به عنوان Administrator
Set-ExecutionPolicy Bypass -Scope Process -Force
[System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072
iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))

# نصب dependencies
choco install -y openssl
choco install -y wkhtmltopdf
```

---

## 🐛 عیب‌یابی

### خطا: `libssl.so.3: cannot open shared object file`

**راه حل (Ubuntu/Debian):**
```bash
sudo apt-get install libssl3
```

**راه حل (Fedora/RHEL):**
```bash
sudo dnf install openssl-libs
```

### خطا: `cannot open display`

**راه حل:**
اطمینان حاصل کنید که X11 یا Wayland در حال اجراست و متغیر `DISPLAY` تنظیم شده است:
```bash
echo $DISPLAY
# باید چیزی مانند :0 یا :1 نمایش دهد
```

### خطا: wkhtmltopdf not found

**راه حل:**
```bash
# Linux
sudo apt-get install wkhtmltopdf

# macOS
brew install --cask wkhtmltopdf
```

### macOS: "App is damaged and can't be opened"

**راه حل:**
```bash
xattr -cr "Ravro Decryption Tool.app"
```

### Windows: DLL not found errors

**راه حل:**
اطمینان حاصل کنید که OpenSSL نصب شده است:
```powershell
choco install -y openssl
```

سپس فایل‌های DLL را کپی کنید:
```powershell
copy "C:\Program Files\OpenSSL-Win64\bin\*.dll" .
```

### Windows: wkhtmltopdf not found

**راه حل:**
```powershell
choco install -y wkhtmltopdf
```

---

## 📞 پشتیبانی

در صورت بروز مشکل:
- [مشاهده Issues در GitHub](https://github.com/ravro-ir/ravro_dcrpt/issues)
- [مشاهده مستندات کامل](https://github.com/ravro-ir/ravro_dcrpt)

---

## 📝 لایسنس

این نرم‌افزار تحت لایسنس [LICENSE] منتشر شده است.

