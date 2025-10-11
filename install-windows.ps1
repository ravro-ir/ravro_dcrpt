# Ravro Decryption Tool - Windows Installation Script
# This script installs all required dependencies for running the GUI

# Check if running as Administrator
$currentPrincipal = New-Object Security.Principal.WindowsPrincipal([Security.Principal.WindowsIdentity]::GetCurrent())
$isAdmin = $currentPrincipal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)

if (-not $isAdmin) {
    Write-Host ""
    Write-Host "========================================" -ForegroundColor Red
    Write-Host "  ERROR: Administrator Rights Required  " -ForegroundColor Red
    Write-Host "========================================" -ForegroundColor Red
    Write-Host ""
    Write-Host "Please run this script as Administrator:" -ForegroundColor Yellow
    Write-Host "  1. Right-click on PowerShell" -ForegroundColor Cyan
    Write-Host "  2. Select 'Run as Administrator'" -ForegroundColor Cyan
    Write-Host "  3. Run this script again" -ForegroundColor Cyan
    Write-Host ""
    pause
    exit 1
}

Write-Host ""
Write-Host "╔═══════════════════════════════════════════════════════╗" -ForegroundColor Blue
Write-Host "║   Ravro Decryption Tool - Dependency Installation    ║" -ForegroundColor Blue
Write-Host "║                    Windows                            ║" -ForegroundColor Blue
Write-Host "╚═══════════════════════════════════════════════════════╝" -ForegroundColor Blue
Write-Host ""

# Function to print status
function Print-Status {
    param([string]$message)
    Write-Host "→ $message" -ForegroundColor Blue
}

function Print-Success {
    param([string]$message)
    Write-Host "✓ $message" -ForegroundColor Green
}

function Print-Warning {
    param([string]$message)
    Write-Host "⚠ $message" -ForegroundColor Yellow
}

function Print-Error {
    param([string]$message)
    Write-Host "✗ $message" -ForegroundColor Red
}

# Check Windows version
$osVersion = [System.Environment]::OSVersion.Version
Write-Host "✓ Windows Version: $($osVersion.Major).$($osVersion.Minor)" -ForegroundColor Green
Write-Host ""

# Check if Chocolatey is installed
Print-Status "Checking for Chocolatey package manager..."
if (-not (Get-Command choco -ErrorAction SilentlyContinue)) {
    Print-Warning "Chocolatey is not installed"
    Print-Status "Installing Chocolatey..."
    
    # Set execution policy
    Set-ExecutionPolicy Bypass -Scope Process -Force
    
    # Install Chocolatey
    [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072
    Invoke-Expression ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))
    
    # Refresh environment
    $env:Path = [System.Environment]::GetEnvironmentVariable("Path","Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path","User")
    
    Print-Success "Chocolatey installed"
} else {
    Print-Success "Chocolatey is already installed"
}

Write-Host ""

# Install OpenSSL
Print-Status "Installing OpenSSL..."
if (Test-Path "C:\Program Files\OpenSSL-Win64\bin\openssl.exe") {
    Print-Success "OpenSSL is already installed"
} else {
    choco install openssl -y
    if ($LASTEXITCODE -eq 0) {
        Print-Success "OpenSSL installed successfully"
    } else {
        Print-Error "Failed to install OpenSSL"
        exit 1
    }
}

# Install wkhtmltopdf
Print-Status "Installing wkhtmltopdf..."
if (Test-Path "C:\Program Files\wkhtmltopdf\bin\wkhtmltopdf.exe") {
    Print-Success "wkhtmltopdf is already installed"
} else {
    choco install wkhtmltopdf -y
    if ($LASTEXITCODE -eq 0) {
        Print-Success "wkhtmltopdf installed successfully"
    } else {
        Print-Error "Failed to install wkhtmltopdf"
        exit 1
    }
}

# Refresh environment PATH
$env:Path = [System.Environment]::GetEnvironmentVariable("Path","Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path","User")

# Verify installations
Write-Host ""
Print-Status "Verifying installations..."

# Check OpenSSL
if (Get-Command openssl -ErrorAction SilentlyContinue) {
    $opensslVersion = & openssl version
    Print-Success "OpenSSL: $opensslVersion"
} else {
    Print-Warning "OpenSSL not found in PATH. You may need to restart your terminal."
}

# Check wkhtmltopdf
if (Get-Command wkhtmltopdf -ErrorAction SilentlyContinue) {
    $wkhtmlVersion = & wkhtmltopdf --version 2>&1 | Select-String "wkhtmltopdf" | Select-Object -First 1
    Print-Success "wkhtmltopdf: $wkhtmlVersion"
} else {
    Print-Warning "wkhtmltopdf not found in PATH. You may need to restart your terminal."
}

Write-Host ""
Write-Host "╔═══════════════════════════════════════════════════════╗" -ForegroundColor Green
Write-Host "║        Dependencies installed successfully! 🎉        ║" -ForegroundColor Green
Write-Host "╚═══════════════════════════════════════════════════════╝" -ForegroundColor Green
Write-Host ""
Write-Host "Next steps:" -ForegroundColor Blue
Write-Host "  1. Download the latest release from GitHub" -ForegroundColor Cyan
Write-Host "  2. Extract the zip file:" -ForegroundColor Cyan
Write-Host "     Right-click → Extract All" -ForegroundColor Cyan
Write-Host "  3. Run the GUI:" -ForegroundColor Cyan
Write-Host "     Double-click ravro_dcrpt_gui.exe" -ForegroundColor Cyan
Write-Host ""
Write-Host "Note: You may need to restart your terminal or computer" -ForegroundColor Yellow
Write-Host "      for PATH changes to take effect." -ForegroundColor Yellow
Write-Host ""
Write-Host "Download link:" -ForegroundColor Blue
Write-Host "https://github.com/ravro-ir/ravro_dcrpt/releases" -ForegroundColor Cyan
Write-Host ""
pause

