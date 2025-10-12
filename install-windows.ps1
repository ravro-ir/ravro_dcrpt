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
    $machinePath = [System.Environment]::GetEnvironmentVariable("Path", "Machine")
    $userPath = [System.Environment]::GetEnvironmentVariable("Path", "User")
    $env:Path = $machinePath + ";" + $userPath
    
    Print-Success "Chocolatey installed"
} else {
    Print-Success "Chocolatey is already installed"
}

Write-Host ""

# Install OpenSSL
Print-Status "Installing OpenSSL..."

# Check if already installed
$opensslLocations = @(
    "C:\OpenSSL-Win64",
    "C:\Program Files\OpenSSL-Win64",
    "C:\Program Files\OpenSSL"
)

$opensslFound = $false
foreach ($loc in $opensslLocations) {
    if (Test-Path "$loc\bin\openssl.exe") {
        Print-Success "OpenSSL is already installed at: $loc"
        $opensslFound = $true
        break
    }
}

if (-not $opensslFound) {
    Print-Status "Attempting to install OpenSSL..."
    
    # Try different versions from slproweb.com
    $OPENSSL_VERSIONS = @("3_4_0", "3_3_2", "3_3_1", "3_3_0", "3_2_0", "3_1_0")
    $OPENSSL_INSTALLED = $false
    
    foreach ($VERSION in $OPENSSL_VERSIONS) {
        try {
            Write-Host "  → Trying OpenSSL version $VERSION..." -ForegroundColor Cyan
            $OPENSSL_URL = "https://slproweb.com/download/Win64OpenSSL-${VERSION}.exe"
            $OPENSSL_INSTALLER = "$env:TEMP\openssl-installer.exe"
            
            # Try to download
            Invoke-WebRequest -Uri $OPENSSL_URL -OutFile $OPENSSL_INSTALLER -ErrorAction Stop -UseBasicParsing
            
            Print-Success "Downloaded OpenSSL $VERSION"
            Print-Status "Installing OpenSSL to C:\OpenSSL-Win64..."
            
            # Install silently
            Start-Process -FilePath $OPENSSL_INSTALLER -ArgumentList "/VERYSILENT /SP- /SUPPRESSMSGBOXES /DIR=C:\OpenSSL-Win64" -Wait -NoNewWindow
            
            if (Test-Path "C:\OpenSSL-Win64\bin\openssl.exe") {
                Print-Success "OpenSSL $VERSION installed successfully at C:\OpenSSL-Win64"
                $OPENSSL_INSTALLED = $true
                
                # Add to PATH
                $currentPath = [Environment]::GetEnvironmentVariable("Path", "Machine")
                if ($currentPath -notlike "*C:\OpenSSL-Win64\bin*") {
                    [Environment]::SetEnvironmentVariable("Path", "$currentPath;C:\OpenSSL-Win64\bin", "Machine")
                    Print-Success "Added OpenSSL to system PATH"
                }
                
                break
            }
        } catch {
            Write-Host "  ⚠ Version $VERSION not available, trying next..." -ForegroundColor Yellow
            continue
        }
    }
    
    # Fallback to chocolatey if direct download failed
    if (-not $OPENSSL_INSTALLED) {
        Print-Warning "Direct download failed, using Chocolatey..."
        choco install openssl -y
        
        if ($LASTEXITCODE -eq 0) {
            # Check where chocolatey installed it
            foreach ($loc in $opensslLocations) {
                if (Test-Path "$loc\bin\openssl.exe") {
                    Print-Success "OpenSSL installed via Chocolatey at: $loc"
                    
                    # Create junction to standard location if needed
                    if ($loc -ne "C:\OpenSSL-Win64" -and -not (Test-Path "C:\OpenSSL-Win64")) {
                        New-Item -ItemType Junction -Path "C:\OpenSSL-Win64" -Target $loc -ErrorAction SilentlyContinue
                        Print-Success "Created junction from C:\OpenSSL-Win64 to $loc"
                    }
                    
                    $OPENSSL_INSTALLED = $true
                    break
                }
            }
        }
    }
    
    if (-not $OPENSSL_INSTALLED) {
        Print-Error "Failed to install OpenSSL"
        Write-Host "Please install OpenSSL manually from: https://slproweb.com/products/Win32OpenSSL.html" -ForegroundColor Yellow
        exit 1
    }
}

# Install wkhtmltopdf
Print-Status "Installing wkhtmltopdf..."
choco install wkhtmltopdf -y --force
if ($LASTEXITCODE -eq 0) {
    Print-Success "wkhtmltopdf installed successfully"
} else {
    Print-Error "Failed to install wkhtmltopdf"
    exit 1
}

# Verify installations
Write-Host ""
Print-Status "Verifying installations..."

# Refresh environment PATH
$machinePath = [System.Environment]::GetEnvironmentVariable("Path", "Machine")
$userPath = [System.Environment]::GetEnvironmentVariable("Path", "User")
$env:Path = $machinePath + ";" + $userPath

# Check OpenSSL
$opensslFound = $false
if (Get-Command openssl -ErrorAction SilentlyContinue) {
    $opensslVersion = & openssl version
    Print-Success "OpenSSL: $opensslVersion"
    $opensslFound = $true
} else {
    # Try direct paths
    foreach ($loc in $opensslLocations) {
        if (Test-Path "$loc\bin\openssl.exe") {
            $opensslVersion = & "$loc\bin\openssl.exe" version
            Print-Success "OpenSSL: $opensslVersion (at $loc)"
            Print-Warning "OpenSSL not in PATH. Restart your terminal or add manually:"
            Write-Host "  $loc\bin" -ForegroundColor Cyan
            $opensslFound = $true
            break
        }
    }
    
    if (-not $opensslFound) {
        Print-Warning "OpenSSL not found. Please restart your terminal or computer."
    }
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

# Show installation paths
Write-Host "Installation Summary:" -ForegroundColor Blue
Write-Host "════════════════════" -ForegroundColor Blue

# OpenSSL location
foreach ($loc in $opensslLocations) {
    if (Test-Path "$loc\bin\openssl.exe") {
        Write-Host "  OpenSSL: $loc" -ForegroundColor Cyan
        break
    }
}

# wkhtmltopdf location
if (Test-Path "C:\Program Files\wkhtmltopdf\bin\wkhtmltopdf.exe") {
    Write-Host "  wkhtmltopdf: C:\Program Files\wkhtmltopdf" -ForegroundColor Cyan
}

Write-Host ""
Write-Host "Next steps:" -ForegroundColor Blue
Write-Host "  1. Download the latest release from GitHub" -ForegroundColor Cyan
Write-Host "  2. Extract the zip file:" -ForegroundColor Cyan
Write-Host "     Right-click → Extract All" -ForegroundColor Cyan
Write-Host "  3. Run the GUI:" -ForegroundColor Cyan
Write-Host "     Double-click ravro_dcrpt_gui.exe" -ForegroundColor Cyan
Write-Host ""
Write-Host "Important:" -ForegroundColor Yellow
Write-Host "  • Restart your terminal or PowerShell for PATH changes" -ForegroundColor Yellow
Write-Host "  • If the app shows DLL errors, make sure all dependencies" -ForegroundColor Yellow
Write-Host "    are in PATH or in the same folder as the executable" -ForegroundColor Yellow
Write-Host ""
Write-Host "Download link:" -ForegroundColor Blue
Write-Host "https://github.com/ravro-ir/ravro_dcrpt/releases" -ForegroundColor Cyan
Write-Host ""

# Offer to open the download page
$openBrowser = Read-Host "Open download page in browser? (Y/N)"
if ($openBrowser -eq "Y" -or $openBrowser -eq "y") {
    Start-Process "https://github.com/ravro-ir/ravro_dcrpt/releases"
}

Write-Host ""
pause