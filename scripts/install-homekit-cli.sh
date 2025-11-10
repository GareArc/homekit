#!/bin/bash
# HomeKit CLI Installation Script
# One-line installation: curl -fsSL https://raw.githubusercontent.com/GareArc/homekit/main/scripts/install-homekit-cli.sh | bash

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
REPO="GareArc/homekit"

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

# Function to detect OS and architecture
detect_platform() {
    local os arch
    os=$(uname -s | tr '[:upper:]' '[:lower:]')
    arch=$(uname -m | tr '[:upper:]' '[:lower:]')
    
    # Normalize architecture names
    case "$arch" in
        x86_64|amd64)
            arch="amd64"
            ;;
        aarch64|arm64)
            arch="arm64"
            ;;
        armv7l)
            arch="arm"
            ;;
    esac
    
    echo "${os}-${arch}"
}

# Function to get latest release tag
get_latest_release() {
    local repo="$1"
    local api_url="https://api.github.com/repos/${repo}/releases/latest"
    
    if command -v curl >/dev/null 2>&1; then
        curl -fsSL "$api_url" | grep -o '"tag_name": "[^"]*' | grep -o '[^"]*$'
    elif command -v wget >/dev/null 2>&1; then
        wget -qO- "$api_url" | grep -o '"tag_name": "[^"]*' | grep -o '[^"]*$'
    else
        print_error "Neither curl nor wget is available. Please install one of them."
        exit 1
    fi
}

# Function to download and install binary
install_binary() {
    local repo="$1"
    local version="$2"
    local platform="$3"
    local install_dir="$4"
    local binary_name="$5"
    
    # Create installation directory if it doesn't exist
    mkdir -p "$install_dir"
    
    # Construct download URL
    # For linux-amd64, use the 'homekit' binary directly (built by make build)
    # For other platforms, use the platform-specific binary
    local download_url
    if [ "$platform" = "linux-amd64" ]; then
        download_url="https://github.com/${repo}/releases/download/${version}/${binary_name}"
    else
        download_url="https://github.com/${repo}/releases/download/${version}/${binary_name}-${platform}"
    fi
    
    print_status "Downloading ${binary_name} ${version} for ${platform}..."
    print_status "URL: ${download_url}"
    
    # Download the binary
    local temp_file=$(mktemp)
    local http_code
    if command -v curl >/dev/null 2>&1; then
        http_code=$(curl -fsSL -o "$temp_file" -w "%{http_code}" "$download_url" 2>/dev/null || echo "000")
        if [ "$http_code" != "200" ]; then
            print_error "Failed to download binary (HTTP $http_code). URL may be incorrect or release asset missing."
            print_error "Expected: ${download_url}"
            rm -f "$temp_file"
            exit 1
        fi
    elif command -v wget >/dev/null 2>&1; then
        if ! wget -qO "$temp_file" "$download_url" 2>/dev/null; then
            print_error "Failed to download binary from GitHub releases"
            rm -f "$temp_file"
            exit 1
        fi
    fi
    
    # Verify downloaded file is a valid binary
    if [ ! -s "$temp_file" ]; then
        print_error "Downloaded file is empty"
        rm -f "$temp_file"
        exit 1
    fi
    
    # Check if file is a valid binary (ELF, Mach-O, or PE)
    if command -v file >/dev/null 2>&1; then
        local file_type=$(file -b "$temp_file" 2>/dev/null || echo "")
        if [[ ! "$file_type" =~ (ELF|Mach-O|PE32|executable) ]]; then
            print_error "Downloaded file is not a valid binary. File type: ${file_type}"
            print_error "This usually means the release asset doesn't exist or the URL is incorrect."
            print_error "First 100 bytes of downloaded file:"
            head -c 100 "$temp_file" | cat -A
            echo
            rm -f "$temp_file"
            exit 1
        fi
    fi
    
    # Make binary executable
    chmod +x "$temp_file"
    
    # Move to installation directory
    local install_path="${install_dir}/${binary_name}"
    mv "$temp_file" "$install_path"
    
    print_success "Installed ${binary_name} to ${install_path}"
    
    # Check if installation directory is in PATH
    if [[ ":$PATH:" != *":${install_dir}:"* ]]; then
        print_warning "${install_dir} is not in your PATH"
        print_status "Add the following line to your shell profile (~/.bashrc, ~/.zshrc, etc.):"
        echo "export PATH=\"${install_dir}:\$PATH\""
        echo
        print_status "Or run: echo 'export PATH=\"${install_dir}:\$PATH\"' >> ~/.bashrc"
    fi
    
    # Verify installation
    if command -v "$binary_name" >/dev/null 2>&1; then
        print_success "Installation verified! You can now run '${binary_name} --help'"
    else
        print_warning "Installation completed, but ${binary_name} is not in your PATH"
        print_status "You may need to restart your terminal or run: source ~/.bashrc"
    fi
}

# Main installation function
main() {
    print_status "Installing HomeKit CLI..."
    
    # Detect platform
    local platform
    platform=$(detect_platform)
    print_status "Detected platform: ${platform}"
    
    # Get latest release
    print_status "Fetching latest release..."
    local version
    version=$(get_latest_release "$REPO")
    print_status "Latest version: ${version}"
    
    # Install binary
    install_binary "$REPO" "$version" "$platform" "$INSTALL_DIR" "$BINARY_NAME"
    
    print_success "HomeKit CLI installation completed!"
}

# Run main function
main "$@"
