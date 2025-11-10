#!/bin/bash
# Local HomeKit CLI installer (for development/testing)
# This script installs from a local build instead of GitHub releases

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Default installation directory
WORK_DIR="$(pwd)"

INSTALL_DIR="${HOME}/.local/bin"
BINARY_NAME="homekit"
SOURCE_DIR="$WORK_DIR/homekit-cli/dist"

print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if source binary exists
if [[ ! -f "${SOURCE_DIR}/${BINARY_NAME}" ]]; then
    print_error "HomeKit CLI binary not found at ${SOURCE_DIR}/${BINARY_NAME}"
    print_status "Please build the project first: cd /work/homekit-cli && make build"
    exit 1
fi

print_status "Installing HomeKit CLI from local build..."

# Create installation directory if it doesn't exist
mkdir -p "$INSTALL_DIR"

# Copy binary to installation directory
cp "${SOURCE_DIR}/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
chmod +x "${INSTALL_DIR}/${BINARY_NAME}"

print_success "Installed ${BINARY_NAME} to ${INSTALL_DIR}/${BINARY_NAME}"

# Check if installation directory is in PATH
if [[ ":$PATH:" != *":${INSTALL_DIR}:"* ]]; then
    print_warning "${INSTALL_DIR} is not in your PATH"
    print_status "Add the following line to your shell profile (~/.bashrc, ~/.zshrc, etc.):"
    echo "export PATH=\"${INSTALL_DIR}:\$PATH\""
    echo
    print_status "Or run: echo 'export PATH=\"${INSTALL_DIR}:\$PATH\"' >> ~/.bashrc"
fi

# Verify installation
if command -v "$BINARY_NAME" >/dev/null 2>&1; then
    print_success "Installation verified! You can now run '${BINARY_NAME} --help'"
else
    print_warning "Installation completed, but ${BINARY_NAME} is not in your PATH"
    print_status "You may need to restart your terminal or run: source ~/.bashrc"
fi

print_success "HomeKit CLI local installation completed!"
