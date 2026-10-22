# نصب Ravro Decryption Tool

این راهنما روش نصب Ravro Decryption Tool را در پلتفرم‌های مختلف توضیح می‌دهد.

## 📦 Linux

### دانلود و نصب

```bash
# دانلود tarball
wget https://github.com/ravro-ir/ravro_dcrpt/releases/latest/download/ravro_dcrpt-linux-amd64.tar.gz

# استخراج
tar -xzf ravro_dcrpt-linux-amd64.tar.gz

# اجرا
./ravro_dcrpt_gui
```

**نصب dependency ها:**
```bash
# Ubuntu/Debian/Kali Linux
sudo apt-get install libgl1 libx11-6 libssl3

# Kali Linux (نصب کامل dependency ها)
sudo apt-get update
sudo apt-get install -y \
    libgl1-mesa-dev \
    libx11-dev \
    libxcursor-dev \
    libxrandr-dev \
    libxinerama-dev \
    libxi-dev \
    libssl-dev

# Fedora/RHEL
sudo dnf install mesa-libGL libX11 openssl-libs

# Arch Linux
sudo pacman -S mesa libx11 openssl
```

---

## 🐉 Kali Linux

### نصب ویژه برای Kali Linux

```bash
# دانلود نسخه مخصوص Kali
wget https://github.com/ravro-ir/ravro_dcrpt/releases/latest/download/ravro_dcrpt-kali-linux-amd64.tar.gz

# استخراج
tar -xzf ravro_dcrpt-kali-linux-amd64.tar.gz

# نصب dependency های Kali
sudo apt-get update
sudo apt-get install -y \
    libgl1-mesa-dev \
    libgl1-mesa-glx \
    xorg-dev \
    libx11-dev \
    libxcursor-dev \
    libxrandr-dev \
    libxinerama-dev \
    libxi-dev \
    libssl-dev \
    pkg-config

# اجرا
./ravro_dcrpt_gui
```

**نکات مهم برای Kali Linux:**
- اطمینان حاصل کنید که X11 یا Wayland در حال اجرا است
- برای استفاده در محیط headless، از VNC یا X11 forwarding استفاده کنید
- برای penetration testing، فایل‌های رمزگذاری شده را در محیط ایزوله بررسی کنید

**استفاده در محیط CLI:**
```bash
# رمزگشایی فایل
./ravro_dcrpt decrypt input.ravro output.pdf --key key.private

# نمایش اطلاعات فایل
./ravro_dcrpt info input.ravro
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

| فرمت | پلتفرم | توضیحات |
|------|--------|---------|
| `.tar.gz` | Linux | Binary + dependencies |
| `.tar.gz` | Kali Linux | Optimized for Kali + dependencies |
| `.tar.gz` | macOS | Application bundle |
| `.zip` | Windows | Executable + DLLs |

---

**🎉 از استفاده از Ravro Decryption Tool لذت ببرید!**
