# نصب Ravro Decryption Tool

این راهنما روش‌های مختلف نصب Ravro Decryption Tool را در پلتفرم‌های مختلف توضیح می‌دهد.

## 📦 Linux

### روش 1: AppImage (توصیه می‌شود)

AppImage یک فایل قابل اجرای مستقل است که نیازی به نصب ندارد:

```bash
# دانلود AppImage
wget https://github.com/ravro-ir/ravro_dcrpt/releases/latest/download/Ravro_Decryption_Tool-x86_64.AppImage

# قابل اجرا کردن
chmod +x Ravro_Decryption_Tool-x86_64.AppImage

# اجرا
./Ravro_Decryption_Tool-x86_64.AppImage
```

**مزایا:**
- ✅ نیازی به نصب ندارد
- ✅ همه dependency ها داخل آن است
- ✅ روی تمام توزیع‌های Linux کار می‌کند
- ✅ می‌توانید در هر مسیری اجرا کنید

### روش 2: DEB Package (Ubuntu/Debian)

برای Ubuntu، Debian و توزیع‌های مبتنی بر آن‌ها:

```bash
# دانلود DEB package
wget https://github.com/ravro-ir/ravro_dcrpt/releases/latest/download/ravro-decryption-tool-amd64.deb

# نصب
sudo dpkg -i ravro-decryption-tool-amd64.deb

# رفع مشکل dependency ها (در صورت نیاز)
sudo apt-get install -f

# اجرا
ravro_dcrpt_gui
```

**مزایا:**
- ✅ یکپارچه با سیستم
- ✅ می‌توانید از Application Menu اجرا کنید
- ✅ به راحتی قابل حذف است (`sudo apt remove ravro-decryption-tool`)

**حذف:**
```bash
sudo apt remove ravro-decryption-tool
```

### روش 3: Tarball (همه توزیع‌ها)

برای نصب دستی:

```bash
# دانلود tarball
wget https://github.com/ravro-ir/ravro_dcrpt/releases/latest/download/ravro_dcrpt-linux-amd64.tar.gz

# استخراج
tar -xzf ravro_dcrpt-linux-amd64.tar.gz

# اجرا
./ravro_dcrpt_gui
```

**توجه:** ممکن است نیاز به نصب dependency ها باشد:
```bash
# Ubuntu/Debian
sudo apt-get install libgl1 libx11-6 libssl3

# Fedora/RHEL
sudo dnf install mesa-libGL libX11 openssl-libs

# Arch Linux
sudo pacman -S mesa libx11 openssl
```

---

## 🍎 macOS

### دانلود و نصب

```bash
# برای Intel Macs
wget https://github.com/ravro-ir/ravro_dcrpt/releases/latest/download/ravro_dcrpt-darwin-amd64.tar.gz
tar -xzf ravro_dcrpt-darwin-amd64.tar.gz

# برای Apple Silicon (M1/M2/M3)
wget https://github.com/ravro-ir/ravro_dcrpt/releases/latest/download/ravro_dcrpt-darwin-arm64.tar.gz
tar -xzf ravro_dcrpt-darwin-arm64.tar.gz

# انتقال به Applications
mv "Ravro Decryption Tool.app" /Applications/

# اجرا
open "/Applications/Ravro Decryption Tool.app"
```

**نکته امنیتی:**
اگر با پیام "cannot be opened because it is from an unidentified developer" مواجه شدید:

```bash
xattr -cr "/Applications/Ravro Decryption Tool.app"
```

یا از System Preferences > Security & Privacy اجازه اجرا را بدهید.

---

## 🪟 Windows

### دانلود و نصب

1. دانلود فایل ZIP:
   ```
   https://github.com/ravro-ir/ravro_dcrpt/releases/latest/download/ravro_dcrpt-windows-amd64.zip
   ```

2. استخراج ZIP file

3. اجرای `ravro_dcrpt_gui.exe`

**نکته:** فایل‌های DLL باید در کنار EXE قرار داشته باشند.

---

## 🔧 نیازمندی‌های سیستم

### Linux
- **معماری:** x86_64 (64-bit)
- **Kernel:** 3.10+
- **Libraries:** OpenGL, X11, OpenSSL

### macOS
- **نسخه:** macOS 10.13+ (High Sierra یا جدیدتر)
- **معماری:** Intel (x86_64) یا Apple Silicon (ARM64)

### Windows
- **نسخه:** Windows 10/11 (64-bit)
- **معماری:** x86_64 (64-bit)

---

## 🐛 مشکلات رایج

### Linux: "error while loading shared libraries"

```bash
# نصب dependency های لازم
sudo apt-get install libgl1-mesa-glx libx11-6 libssl3
```

### Linux: AppImage اجرا نمی‌شود

```bash
# نصب FUSE
sudo apt-get install fuse libfuse2

# یا اجرا با extract mode
./Ravro_Decryption_Tool-x86_64.AppImage --appimage-extract-and-run
```

### macOS: "damaged and can't be opened"

```bash
# حذف quarantine attribute
xattr -cr "/Applications/Ravro Decryption Tool.app"
```

### Windows: "VCRUNTIME140.dll was not found"

دانلود و نصب [Microsoft Visual C++ Redistributable](https://aka.ms/vs/17/release/vc_redist.x64.exe)

---

## 📝 بررسی نسخه

بعد از نصب، می‌توانید نسخه برنامه را چک کنید:

```bash
# Linux
ravro_dcrpt_gui --version

# macOS
"/Applications/Ravro Decryption Tool.app/Contents/MacOS/ravro_dcrpt_gui" --version

# Windows
ravro_dcrpt_gui.exe --version
```

---

## 🆘 پشتیبانی

اگر مشکلی در نصب یا اجرا دارید:

1. [Issues](https://github.com/ravro-ir/ravro_dcrpt/issues) را بررسی کنید
2. یک [New Issue](https://github.com/ravro-ir/ravro_dcrpt/issues/new) ایجاد کنید
3. به [Releases](https://github.com/ravro-ir/ravro_dcrpt/releases) مراجعه کنید برای آخرین نسخه

---

## 📦 فرمت‌های Package

| فرمت | پلتفرم | مزایا |
|------|--------|-------|
| `.AppImage` | Linux | قابل حمل، بدون نیاز به نصب |
| `.deb` | Ubuntu/Debian | یکپارچه با سیستم، قابل update |
| `.tar.gz` | همه | ساده، قابل کنترل |
| `.app` | macOS | Native macOS application |
| `.zip` | Windows | استاندارد Windows |

---

**🎉 از استفاده از Ravro Decryption Tool لذت ببرید!**
