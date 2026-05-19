#!/bin/bash
set -e

REPO="darimuri/open-notebook-cli"
INSTALL_DIR="${HOME}/.local/bin"
BINARY_NAME="open-notebook"

# Detect OS
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
case "$OS" in
    darwin*) OS="darwin" ;;
    linux*) OS="linux" ;;
    *) echo "Unsupported OS: $OS"; exit 1 ;;
esac

# Detect architecture
ARCH="$(uname -m)"
case "$ARCH" in
    x86_64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

# Get latest version
LATEST=$(curl -sL "https://api.github.com/repos/${REPO}/releases/latest" | grep tag_name | cut -d'"' -f4)
if [ -z "$LATEST" ]; then
    echo "Failed to get latest version"
    exit 1
fi

# Determine binary name
if [ "$OS" = "windows" ]; then
    BINARY_NAME="open-notebook.exe"
fi
FILENAME="open-notebook-${OS}-${ARCH}${EXT}"

# Download URL
DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${LATEST}/${FILENAME}"

echo "Installing ${REPO} ${LATEST}..."
echo "Downloading ${DOWNLOAD_URL}..."

# Create install directory if not exists
mkdir -p "${INSTALL_DIR}"

# Download and install
curl -sL "$DOWNLOAD_URL" -o "${INSTALL_DIR}/${BINARY_NAME}"
chmod +x "${INSTALL_DIR}/${BINARY_NAME}"

echo "Installed to ${INSTALL_DIR}/${BINARY_NAME}"

# Verify
"${INSTALL_DIR}/${BINARY_NAME}" --version