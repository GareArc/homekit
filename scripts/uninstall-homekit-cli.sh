#!/bin/bash
# HomeKit CLI Uninstallation Script
# Removes the HomeKit CLI binary installed by the installation scripts

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Default installation directory
INSTALL_DIR="${HOME}/.local/bin"
BINARY_NAME="homekit"
INSTALL_PATH="${INSTALL_DIR}/${BINARY_NAME}"

# Function to print colored output
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

# Main uninstallation function
main() {
    print_status "Uninstalling HomeKit CLI..."
    
    # Check if binary exists
    if [[ ! -f "$INSTALL_PATH" ]]; then
        print_warning "HomeKit CLI binary not found at ${INSTALL_PATH}"
        print_status "Nothing to uninstall."
        exit 0
    fi
    
    # Remove the binary
    print_status "Removing ${INSTALL_PATH}..."
    rm -f "$INSTALL_PATH"
    print_success "Removed HomeKit CLI binary"
    
    # Check if installation directory is empty (optional cleanup)
    if [[ -d "$INSTALL_DIR" ]] && [[ -z "$(ls -A "$INSTALL_DIR" 2>/dev/null)" ]]; then
        print_status "Installation directory ${INSTALL_DIR} is now empty"
        read -p "Remove empty directory ${INSTALL_DIR}? (y/N) " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            rmdir "$INSTALL_DIR"
            print_success "Removed empty directory ${INSTALL_DIR}"
        fi
    fi
    
    # Check for PATH modifications (informational only)
    if [[ ":$PATH:" == *":${INSTALL_DIR}:"* ]]; then
        print_warning "${INSTALL_DIR} is still in your PATH"
        print_status "If you added it manually for HomeKit CLI, you may want to remove it from your shell profile"
        print_status "Check and remove the following line from ~/.bashrc, ~/.zshrc, etc.:"
        echo "export PATH=\"${INSTALL_DIR}:\$PATH\""
    fi
    
    print_success "HomeKit CLI uninstallation completed!"
}

# Run main function
main "$@"
