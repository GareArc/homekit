# HomeKit CLI Installation Scripts

This directory contains installation scripts for the HomeKit CLI tool.

## Quick Installation (One-liner)

The simplest way to install HomeKit CLI is using the one-liner script:

```bash
curl -fsSL https://raw.githubusercontent.com/GareArc/homekit/main/scripts/install-homekit-cli.sh | bash
```

## Manual Installation

If you prefer to download and run the script manually:

```bash
# Download the script
curl -fsSL https://raw.githubusercontent.com/GareArc/homekit/main/scripts/install-homekit-cli.sh -o install-homekit-cli.sh

# Make it executable
chmod +x install-homekit-cli.sh

# Run the installation
./install-homekit-cli.sh
```

## What the Script Does

1. **Detects your platform** (Linux/macOS, amd64/arm64)
2. **Fetches the latest release** from GitHub
3. **Downloads the appropriate binary** for your platform
4. **Installs to `~/.local/bin`** (creates directory if needed)
5. **Makes the binary executable**
6. **Provides PATH setup instructions** if needed

## Installation Directory

By default, the script installs HomeKit CLI to `~/.local/bin/homekit`. This directory is commonly included in the PATH on most Linux distributions and macOS.

## Adding to PATH

If `~/.local/bin` is not in your PATH, add this line to your shell profile:

```bash
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc
```

For other shells, add the same line to:
- **Zsh**: `~/.zshrc`
- **Fish**: `~/.config/fish/config.fish`
- **PowerShell**: `$PROFILE`

## Verification

After installation, verify it works:

```bash
homekit --version
homekit --help
```

## Requirements

- `curl` or `wget` (for downloading)
- `bash` (for running the script)
- Internet connection (for fetching releases)

## Troubleshooting

### Binary not found after installation
- Ensure `~/.local/bin` is in your PATH
- Restart your terminal or run `source ~/.bashrc`

### Permission denied
- Make sure the script is executable: `chmod +x install-homekit-cli.sh`
- Check that `~/.local/bin` is writable

### Download failed
- Check your internet connection
- Verify the GitHub repository is accessible
- Try running the script again

## Script Features

- **Cross-platform**: Supports Linux and macOS on amd64 and arm64
- **Automatic updates**: Always installs the latest release
- **Safe installation**: Uses temporary files and proper error handling
- **User-friendly**: Provides clear status messages and instructions
- **Minimal dependencies**: Only requires curl/wget and bash
