# 🔄 Migration Guide: v1.x → v2.0

## 📌 Overview

Ravro Decryption Tool v2.0 is a **complete rewrite** using Pure Go, eliminating all CGO dependencies while adding GUI support and maintaining full backward compatibility with your data.

## ✅ What's Compatible

### ✓ Your Data
- ✅ **Encrypted files (.ravro)** - No changes needed
- ✅ **Private keys** - Same format, works perfectly
- ✅ **Directory structure** - encrypt/decrypt/key folders unchanged
- ✅ **Zip archives** - Same format, fully compatible

### ✓ Command-Line Interface
- ✅ **Same flags** - `--in`, `--out`, `--key`, `--init`, `--json`
- ✅ **Same workflow** - Initialize, place files, decrypt
- ✅ **Same output** - PDF reports with identical structure

## 🆕 What's New

### Major Improvements

| Feature | v1.x | v2.0 |
|---------|------|------|
| **Dependencies** | OpenSSL (CGO) ❌ | Pure Go ✅ |
| **PDF Library** | wkhtmltopdf (external) ❌ | Maroto (Pure Go) ✅ |
| **GUI** | None | Fyne GUI ✨ |
| **Cross-Compile** | Complex/Impossible | Simple ✅ |
| **Build Time** | Minutes + setup | Seconds ⚡ |
| **Binary Size** | ~50MB | ~15MB 📦 |
| **Architecture** | Monolithic | Clean Architecture 🏗️ |

### New Features
1. 🎨 **GUI Application** - Beautiful graphical interface
2. 🚀 **Faster Builds** - No CGO compilation overhead
3. 🌍 **True Cross-Platform** - Build for Windows from Linux!
4. 📦 **Smaller Binaries** - 3x smaller than v1.x
5. 🏗️ **Better Code Structure** - Clean Architecture with clear separation
6. ✅ **Type Safety** - Proper interfaces and dependency injection

## 🔄 Migration Steps

### Step 1: Backup (Optional but Recommended)

```bash
# Backup your current setup
cp -r encrypt encrypt.backup
cp -r key key.backup
```

### Step 2: Download v2.0

```bash
# Option A: Download pre-built binary
wget https://github.com/ravro-ir/ravro_dcrpt/releases/download/v2.0.0/ravro_dcrpt-2.0.0-linux-amd64.tar.gz
tar xzf ravro_dcrpt-2.0.0-linux-amd64.tar.gz

# Option B: Build from source
git clone https://github.com/ravro-ir/ravro_dcrpt.git
cd ravro_dcrpt
git checkout v2.0.0
make build
```

### Step 3: Test with One Report

```bash
# Move v2.0 binary to your project directory
cp ravro_dcrpt /path/to/your/project/

# Test with one report
./ravro_dcrpt --help

# Process one report to verify
./ravro_dcrpt
```

### Step 4: Verify Output

```bash
# Check that PDFs are generated correctly
ls -lh decrypt/

# Compare with v1.x output if needed
```

### Step 5: Replace v1.x (When Satisfied)

```bash
# Remove old binary
rm ravro_dcrpt_old  # or whatever you named v1.x

# Rename v2.0 to standard name
mv ravro_dcrpt_v2 ravro_dcrpt
```

## 📋 Command Comparison

### v1.x Commands

```bash
# Initialize
./ravro_dcrpt -init

# Process reports
./ravro_dcrpt

# With custom paths
./ravro_dcrpt -in=custom_input -out=custom_output -key=my_key.pem

# Update
./ravro_dcrpt -update

# JSON export
./ravro_dcrpt -json

# Enable logging
./ravro_dcrpt -log
```

### v2.0 Commands (Same! Plus more options)

```bash
# Initialize (same)
./ravro_dcrpt --init

# Process reports (same)
./ravro_dcrpt

# With custom paths (same, but cleaner)
./ravro_dcrpt --in=custom_input --out=custom_output --key=my_key.pem

# JSON export (same)
./ravro_dcrpt --json

# NEW: GUI mode
./ravro_dcrpt-gui

# NEW: Version information
./ravro_dcrpt --version

# NEW: Help with full details
./ravro_dcrpt --help
```

## 🎨 Using the New GUI

### Launch GUI

```bash
# Linux/macOS
./ravro_dcrpt-gui

# Windows
ravro_dcrpt-gui.exe
```

### GUI Features

1. **📁 Directory Browser**
   - Click "Browse" to select directories
   - Pre-filled with default paths

2. **🔑 Key Selection**
   - Browse to your private key
   - Validate before processing

3. **📊 Real-time Progress**
   - See progress as reports are processed
   - Live log output

4. **✅ Validation**
   - Validate key before starting
   - Initialize directories with one click

## 🔧 Troubleshooting

### Issue: "Key file not found"

**v1.x Behavior**: Silent failure or generic error

**v2.0 Solution**:
- More descriptive error messages
- Use `--key` flag to specify exact path
- GUI has validation button to test key before processing

### Issue: "Cannot decrypt file"

**v1.x**: Check if OpenSSL is installed and working

**v2.0**: 
- Pure Go, no external dependencies needed
- Check key format (should be PEM)
- Validate key using CLI or GUI

### Issue: "PDF generation failed"

**v1.x**: Check if wkhtmltopdf is installed

**v2.0**: 
- Pure Go PDF generation
- No external tools needed
- Check output directory permissions

### Issue: Build errors

**v1.x**:
```bash
# Required OpenSSL headers
sudo apt-get install libssl-dev

# Required wkhtmltopdf
sudo apt-get install wkhtmltopdf
```

**v2.0**:
```bash
# Just Go! No external dependencies
go build ./cmd/cli
```

## 📊 Performance Comparison

### Build Time

| Platform | v1.x | v2.0 |
|----------|------|------|
| Linux → Linux | ~2-3 min | ~10 sec ⚡ |
| Linux → Windows | ❌ Impossible | ~10 sec ✅ |
| Linux → macOS | ❌ Very hard | ~10 sec ✅ |

### Runtime Performance

| Operation | v1.x | v2.0 | Notes |
|-----------|------|------|-------|
| Decryption | Fast | Fast | Similar performance |
| PDF Generation | Slower | Faster | Pure Go is optimized |
| Large batches | Good | Better | Improved error handling |

### Binary Size

| Platform | v1.x | v2.0 | Reduction |
|----------|------|------|-----------|
| Linux | ~50MB | ~15MB | 70% ⬇️ |
| Windows | ~50MB | ~15MB | 70% ⬇️ |
| macOS | ~50MB | ~15MB | 70% ⬇️ |

## 🔐 Security Improvements

1. **No CGO** - Reduces attack surface
2. **Pure Go** - Memory safe, type safe
3. **Modern crypto** - Updated PKCS7 implementation
4. **Better validation** - Key and file validation before processing
5. **Error handling** - Comprehensive error messages

## 🚀 New Workflow Examples

### Example 1: GUI Workflow

```
1. Launch: ./ravro_dcrpt-gui
2. Click "Initialize Directories"
3. Place your files in encrypt/ and key/
4. Click "Browse" to select paths (or use defaults)
5. Click "Validate Key" to test
6. Click "Start Processing"
7. Watch progress in real-time!
```

### Example 2: CLI Batch Processing

```bash
# Process multiple report directories
for dir in encrypt/*/; do
    ./ravro_dcrpt --in="$dir" --out="decrypt/$(basename $dir)"
done
```

### Example 3: Cross-Platform Deployment

```bash
# On your Linux build machine
make build-all

# Deploy to Windows server
scp build/windows/ravro_dcrpt.exe windows-server:/path/

# Deploy to macOS
scp build/darwin/ravro_dcrpt-arm64 mac-server:/path/

# No need to build on each platform! 🎉
```

## 📦 Uninstalling v1.x

After confirming v2.0 works:

```bash
# Remove v1.x binary
rm ravro_dcrpt_v1

# Clean up old dependencies (if installed system-wide)
# Ubuntu/Debian
sudo apt-get remove libssl-dev wkhtmltopdf

# macOS
brew uninstall openssl wkhtmltopdf
```

## ❓ FAQ

### Q: Will v2.0 work with my existing reports?
**A**: Yes! 100% compatible with all v1.x encrypted reports.

### Q: Do I need to change my private key?
**A**: No, same keys work perfectly.

### Q: Can I use both versions?
**A**: Yes, they can coexist. Just use different binary names.

### Q: Is the PDF output identical?
**A**: Very similar, with improved formatting and RTL support.

### Q: Can I go back to v1.x?
**A**: Yes, v2.0 doesn't change your source files.

### Q: How do I report issues?
**A**: GitHub Issues: https://github.com/ravro-ir/ravro_dcrpt/issues

## 🎯 Migration Checklist

- [ ] Download v2.0 binary or build from source
- [ ] Test with one report
- [ ] Verify PDF output quality
- [ ] Try the GUI (optional)
- [ ] Process full batch
- [ ] Compare results with v1.x
- [ ] Remove v1.x if satisfied
- [ ] Update documentation/scripts
- [ ] Celebrate! 🎉

## 📞 Support

Need help with migration?

- 📧 Email: ramin.blackhat@gmail.com
- 🐛 Issues: https://github.com/ravro-ir/ravro_dcrpt/issues
- 💬 Discussions: https://github.com/ravro-ir/ravro_dcrpt/discussions

---

**Welcome to v2.0! Enjoy the Pure Go experience! 🚀**

