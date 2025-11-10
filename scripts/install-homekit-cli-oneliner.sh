#!/bin/bash
# One-liner HomeKit CLI installer
# Usage: curl -fsSL https://raw.githubusercontent.com/GareArc/homekit/main/scripts/install-homekit-cli.sh | bash

set -euo pipefail

# Detect platform and install HomeKit CLI
PLATFORM=$(uname -s | tr '[:upper:]' '[:lower:]')-$(uname -m | tr '[:upper:]' '[:lower:]' | sed 's/x86_64/amd64/;s/aarch64/arm64/;s/armv7l/arm/')
VERSION=$(curl -fsSL https://api.github.com/repos/GareArc/homekit/releases/latest | grep -o '"tag_name": "[^"]*' | grep -o '[^"]*$')
INSTALL_DIR="${HOME}/.local/bin"
BINARY_NAME="homekit"

echo "Installing HomeKit CLI ${VERSION} for ${PLATFORM}..."

mkdir -p "$INSTALL_DIR"
TEMP_FILE=$(mktemp)

# For linux-amd64, use the 'homekit' binary directly (built by make build)
# For other platforms, use the platform-specific binary
if [ "$PLATFORM" = "linux-amd64" ]; then
    DOWNLOAD_URL="https://github.com/GareArc/homekit/releases/download/${VERSION}/${BINARY_NAME}"
else
    DOWNLOAD_URL="https://github.com/GareArc/homekit/releases/download/${VERSION}/${BINARY_NAME}-${PLATFORM}"
fi

HTTP_CODE=$(curl -fsSL -o "$TEMP_FILE" -w "%{http_code}" "$DOWNLOAD_URL" 2>/dev/null || echo "000")
if [ "$HTTP_CODE" != "200" ]; then
    echo "❌ Failed to download binary (HTTP $HTTP_CODE)"
    echo "   URL: $DOWNLOAD_URL"
    echo "   This usually means the release asset doesn't exist for this version."
    rm -f "$TEMP_FILE"
    exit 1
fi

# Verify it's a valid binary
if [ ! -s "$TEMP_FILE" ]; then
    echo "❌ Downloaded file is empty"
    rm -f "$TEMP_FILE"
    exit 1
fi

if command -v file >/dev/null 2>&1; then
    FILE_TYPE=$(file -b "$TEMP_FILE" 2>/dev/null || echo "")
    if [[ ! "$FILE_TYPE" =~ (ELF|Mach-O|PE32|executable) ]]; then
        echo "❌ Downloaded file is not a valid binary (type: $FILE_TYPE)"
        echo "   This usually means the release asset doesn't exist."
        rm -f "$TEMP_FILE"
        exit 1
    fi
fi

mv "$TEMP_FILE" "${INSTALL_DIR}/${BINARY_NAME}"
chmod +x "${INSTALL_DIR}/${BINARY_NAME}"

echo "✅ HomeKit CLI installed to ${INSTALL_DIR}/${BINARY_NAME}"

if [[ ":$PATH:" != *":${INSTALL_DIR}:"* ]]; then
    echo "⚠️  Add ${INSTALL_DIR} to your PATH:"
    echo "   echo 'export PATH=\"${INSTALL_DIR}:\$PATH\"' >> ~/.bashrc"
    echo "   source ~/.bashrc"
fi

echo "🚀 Run '${BINARY_NAME} --help' to get started!"
