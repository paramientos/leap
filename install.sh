#!/bin/bash
# LEAP SSH Manager - Installation Script for Linux/macOS
# Usage: curl -fsSL https://raw.githubusercontent.com/paramientos/leap/main/install.sh | bash

set -e

REPO="paramientos/leap"
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="leap"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}"
echo "⚡ LEAP SSH Manager Installer"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo -e "${NC}"

# Detect OS and Architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $ARCH in
    x86_64)
        ARCH="amd64"
        ;;
    aarch64|arm64)
        ARCH="arm64"
        ;;
    *)
        echo -e "${RED}✗ Unsupported architecture: $ARCH${NC}"
        exit 1
        ;;
esac

case $OS in
    linux)
        PLATFORM="linux"
        ;;
    darwin)
        PLATFORM="darwin"
        ;;
    *)
        echo -e "${RED}✗ Unsupported OS: $OS${NC}"
        exit 1
        ;;
esac

echo -e "${GREEN}✓${NC} Detected platform: ${YELLOW}${PLATFORM}-${ARCH}${NC}"

# Get latest release version
echo -e "${BLUE}→${NC} Fetching latest release..."
LATEST_VERSION=$(curl -s "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_VERSION" ]; then
    echo -e "${RED}✗ Failed to fetch latest version${NC}"
    exit 1
fi

echo -e "${GREEN}✓${NC} Latest version: ${YELLOW}${LATEST_VERSION}${NC}"

# Download URL
DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${LATEST_VERSION}/${BINARY_NAME}-${LATEST_VERSION#v}-${PLATFORM}-${ARCH}.tar.gz"

echo -e "${BLUE}→${NC} Downloading from: ${DOWNLOAD_URL}"

# Create temp directory
TMP_DIR=$(mktemp -d)
cd "$TMP_DIR"

# Download and extract
if ! curl -fsSL "$DOWNLOAD_URL" -o "${BINARY_NAME}.tar.gz"; then
    echo -e "${RED}✗ Download failed${NC}"
    rm -rf "$TMP_DIR"
    exit 1
fi

echo -e "${GREEN}✓${NC} Downloaded successfully"

# Extract
tar -xzf "${BINARY_NAME}.tar.gz"

# Check if binary exists
if [ ! -f "${BINARY_NAME}-${PLATFORM}-${ARCH}" ]; then
    echo -e "${RED}✗ Binary not found in archive${NC}"
    rm -rf "$TMP_DIR"
    exit 1
fi

# Install
echo -e "${BLUE}→${NC} Installing to ${INSTALL_DIR}..."

if [ -w "$INSTALL_DIR" ]; then
    mv "${BINARY_NAME}-${PLATFORM}-${ARCH}" "${INSTALL_DIR}/${BINARY_NAME}"
    chmod +x "${INSTALL_DIR}/${BINARY_NAME}"
else
    echo -e "${YELLOW}⚠${NC}  Requires sudo for installation to ${INSTALL_DIR}"
    sudo mv "${BINARY_NAME}-${PLATFORM}-${ARCH}" "${INSTALL_DIR}/${BINARY_NAME}"
    sudo chmod +x "${INSTALL_DIR}/${BINARY_NAME}"
fi

# Cleanup
cd - > /dev/null
rm -rf "$TMP_DIR"

# Verify installation
if command -v leap &> /dev/null; then
    VERSION=$(leap --version 2>&1 | head -n 1)
    echo ""
    echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${GREEN}✓ Installation successful!${NC}"
    echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""
    echo -e "  ${VERSION}"
    echo ""
    echo -e "${BLUE}Quick Start:${NC}"
    echo -e "  ${YELLOW}leap add${NC}        - Add a new SSH connection"
    echo -e "  ${YELLOW}leap list${NC}       - List all connections"
    echo -e "  ${YELLOW}leap${NC}            - Launch interactive TUI"
    echo ""
    echo -e "${BLUE}Documentation:${NC} https://github.com/${REPO}"
    echo ""
else
    echo -e "${RED}✗ Installation failed - binary not found in PATH${NC}"
    exit 1
fi
