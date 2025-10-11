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

---

## 📞 پشتیبانی

در صورت بروز مشکل:
- [مشاهده Issues در GitHub](https://github.com/ravro-ir/ravro_dcrpt/issues)
- [مشاهده مستندات کامل](https://github.com/ravro-ir/ravro_dcrpt)

---

## 📝 لایسنس

این نرم‌افزار تحت لایسنس [LICENSE] منتشر شده است.

