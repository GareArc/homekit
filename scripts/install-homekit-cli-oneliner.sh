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
curl -fsSL "https://github.com/GareArc/homekit/releases/download/${VERSION}/${BINARY_NAME}-${PLATFORM}" -o "${INSTALL_DIR}/${BINARY_NAME}"
chmod +x "${INSTALL_DIR}/${BINARY_NAME}"

echo "✅ HomeKit CLI installed to ${INSTALL_DIR}/${BINARY_NAME}"

if [[ ":$PATH:" != *":${INSTALL_DIR}:"* ]]; then
    echo "⚠️  Add ${INSTALL_DIR} to your PATH:"
    echo "   echo 'export PATH=\"${INSTALL_DIR}:\$PATH\"' >> ~/.bashrc"
    echo "   source ~/.bashrc"
fi

echo "🚀 Run '${BINARY_NAME} --help' to get started!"
