# LEAP SSH Manager - Installation Script for Windows
# Usage: irm https://raw.githubusercontent.com/paramientos/leap/main/install.ps1 | iex

$ErrorActionPreference = "Stop"

$REPO = "paramientos/leap"
$BINARY_NAME = "leap"
$INSTALL_DIR = "$env:LOCALAPPDATA\Programs\LEAP"

Write-Host ""
Write-Host "⚡ LEAP SSH Manager Installer" -ForegroundColor Blue
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Blue
Write-Host ""

# Detect Architecture
$ARCH = if ([Environment]::Is64BitOperatingSystem) { "amd64" } else { "386" }
if ($env:PROCESSOR_ARCHITECTURE -eq "ARM64") {
    $ARCH = "arm64"
}

Write-Host "✓ Detected platform: windows-$ARCH" -ForegroundColor Green

# Get latest release
Write-Host "→ Fetching latest release..." -ForegroundColor Blue
try {
    $latestRelease = Invoke-RestMethod -Uri "https://api.github.com/repos/$REPO/releases/latest"
    $LATEST_VERSION = $latestRelease.tag_name
    Write-Host "✓ Latest version: $LATEST_VERSION" -ForegroundColor Green
} catch {
    Write-Host "✗ Failed to fetch latest version" -ForegroundColor Red
    exit 1
}

# Download URL
$DOWNLOAD_URL = "https://github.com/$REPO/releases/download/$LATEST_VERSION/$BINARY_NAME-$($LATEST_VERSION.TrimStart('v'))-windows-$ARCH.zip"

Write-Host "→ Downloading from: $DOWNLOAD_URL" -ForegroundColor Blue

# Create temp directory
$TMP_DIR = New-Item -ItemType Directory -Path "$env:TEMP\leap-install-$(Get-Random)" -Force

try {
    # Download
    $zipFile = Join-Path $TMP_DIR "$BINARY_NAME.zip"
    Invoke-WebRequest -Uri $DOWNLOAD_URL -OutFile $zipFile -UseBasicParsing
    Write-Host "✓ Downloaded successfully" -ForegroundColor Green

    # Extract
    Expand-Archive -Path $zipFile -DestinationPath $TMP_DIR -Force
    
    # Create install directory
    if (-not (Test-Path $INSTALL_DIR)) {
        New-Item -ItemType Directory -Path $INSTALL_DIR -Force | Out-Null
    }

    # Find the binary
    $binaryPath = Get-ChildItem -Path $TMP_DIR -Filter "$BINARY_NAME-windows-$ARCH.exe" -Recurse | Select-Object -First 1
    
    if (-not $binaryPath) {
        Write-Host "✗ Binary not found in archive" -ForegroundColor Red
        exit 1
    }

    # Copy to install directory
    $targetPath = Join-Path $INSTALL_DIR "$BINARY_NAME.exe"
    Copy-Item -Path $binaryPath.FullName -Destination $targetPath -Force
    Write-Host "✓ Installed to $targetPath" -ForegroundColor Green

    # Add to PATH if not already there
    $userPath = [Environment]::GetEnvironmentVariable("Path", "User")
    if ($userPath -notlike "*$INSTALL_DIR*") {
        Write-Host "→ Adding to PATH..." -ForegroundColor Blue
        [Environment]::SetEnvironmentVariable(
            "Path",
            "$userPath;$INSTALL_DIR",
            "User"
        )
        $env:Path = "$env:Path;$INSTALL_DIR"
        Write-Host "✓ Added to PATH" -ForegroundColor Green
    }

    # Verify installation
    Write-Host ""
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Green
    Write-Host "✓ Installation successful!" -ForegroundColor Green
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Green
    Write-Host ""
    Write-Host "  Installed to: $targetPath" -ForegroundColor White
    Write-Host ""
    Write-Host "Quick Start:" -ForegroundColor Blue
    Write-Host "  leap add        - Add a new SSH connection" -ForegroundColor Yellow
    Write-Host "  leap list       - List all connections" -ForegroundColor Yellow
    Write-Host "  leap            - Launch interactive TUI" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "Documentation: https://github.com/$REPO" -ForegroundColor Blue
    Write-Host ""
    Write-Host "⚠  Please restart your terminal to use 'leap' command" -ForegroundColor Yellow
    Write-Host ""

} catch {
    Write-Host "✗ Installation failed: $_" -ForegroundColor Red
    exit 1
} finally {
    # Cleanup
    Remove-Item -Path $TMP_DIR -Recurse -Force -ErrorAction SilentlyContinue
}
