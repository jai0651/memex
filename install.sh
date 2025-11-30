#!/bin/bash

set -e

REPO_OWNER="jai0651"
REPO_NAME="memex"
REPO_URL="https://github.com/$REPO_OWNER/$REPO_NAME.git"
INSTALL_DIR="$HOME/bin"
CONFIG_DIR="$HOME/.memex"
BINARY_NAME="memex"

echo "Installing Memex..."

# Detect OS and Arch
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

if [ "$ARCH" == "x86_64" ]; then
    ARCH="amd64"
elif [ "$ARCH" == "aarch64" ]; then
    ARCH="arm64"
fi

# Function to install from source
install_from_source() {
    echo "Attempting to build from source..."
    
    # Check for Go
    if ! command -v go &> /dev/null; then
        echo "Error: Go is not installed. Please install Go to build from source."
        exit 1
    fi

    # Check if running from source or need to clone
    IS_REMOTE=false
    if [ ! -f "main.go" ]; then
        IS_REMOTE=true
        if ! command -v git &> /dev/null; then
            echo "Error: Git is not installed. Please install Git first."
            exit 1
        fi
        
        echo "Cloning repository..."
        TEMP_DIR=$(mktemp -d)
        git clone "$REPO_URL" "$TEMP_DIR"
        cd "$TEMP_DIR"
    fi

    # Build the binary
    echo "Building Memex..."
    go build -o "$BINARY_NAME" main.go
    
    # Move binary
    mkdir -p "$INSTALL_DIR"
    mv "$BINARY_NAME" "$INSTALL_DIR/"
    
    # Copy shell integration
    mkdir -p "$CONFIG_DIR"
    cp shell_integration.sh "$CONFIG_DIR/"

    # Cleanup if remote
    if [ "$IS_REMOTE" = true ]; then
        rm -rf "$TEMP_DIR"
    fi
}

# Try to download release
echo "Checking for latest release..."
LATEST_RELEASE=$(curl -s "https://api.github.com/repos/$REPO_OWNER/$REPO_NAME/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -n "$LATEST_RELEASE" ]; then
    echo "Found release: $LATEST_RELEASE"
    DOWNLOAD_URL="https://github.com/$REPO_OWNER/$REPO_NAME/releases/download/$LATEST_RELEASE/memex-$OS-$ARCH"
    
    echo "Downloading binary from: $DOWNLOAD_URL"
    mkdir -p "$INSTALL_DIR"
    
    if curl -L -o "$INSTALL_DIR/$BINARY_NAME" "$DOWNLOAD_URL" --fail; then
        echo "Download successful."
        chmod +x "$INSTALL_DIR/$BINARY_NAME"
        
        # We still need the shell integration script
        # If we are remote, we need to fetch it
        mkdir -p "$CONFIG_DIR"
        if [ ! -f "shell_integration.sh" ]; then
             curl -s "https://raw.githubusercontent.com/$REPO_OWNER/$REPO_NAME/main/shell_integration.sh" -o "$CONFIG_DIR/shell_integration.sh"
        else
             cp shell_integration.sh "$CONFIG_DIR/"
        fi
    else
        echo "Download failed. Falling back to source build."
        install_from_source
    fi
else
    echo "No release found. Falling back to source build."
    install_from_source
fi

# Post-installation setup (Config & Commands)
mkdir -p "$CONFIG_DIR"

# Create default config if not exists
if [ ! -f "$CONFIG_DIR/config.yaml" ]; then
    echo "Creating default config..."
    cat > "$CONFIG_DIR/config.yaml" <<EOF
ghost_text: true
fuzzy_search: true
theme: "dark"
suggestion_limit: 10
command_file: "$CONFIG_DIR/commands.json"
EOF
fi

# Create default commands file if not exists
if [ ! -f "$CONFIG_DIR/commands.json" ]; then
    echo "Creating default commands..."
    cat > "$CONFIG_DIR/commands.json" <<EOF
[
  { "cmd": "git status", "tags": ["git"], "frequency": 1, "last_used": 0 },
  { "cmd": "git pull", "tags": ["git"], "frequency": 1, "last_used": 0 },
  { "cmd": "git push", "tags": ["git"], "frequency": 1, "last_used": 0 },
  { "cmd": "git log --oneline --graph", "tags": ["git"], "frequency": 1, "last_used": 0 },
  { "cmd": "docker ps", "tags": ["docker"], "frequency": 1, "last_used": 0 },
  { "cmd": "docker-compose up -d", "tags": ["docker"], "frequency": 1, "last_used": 0 },
  { "cmd": "ls -la", "tags": ["shell"], "frequency": 1, "last_used": 0 },
  { "cmd": "grep -r \"TODO\" .", "tags": ["shell"], "frequency": 1, "last_used": 0 },
  { "cmd": "find . -type f -name \"*.go\"", "tags": ["shell"], "frequency": 1, "last_used": 0 }
]
EOF
fi

# Add to PATH if not already there (temporary for this session)
export PATH=$INSTALL_DIR:$PATH

echo "Installation complete!"
echo ""
echo "To enable shell integration, add the following to your .zshrc or .bashrc:"
echo "source $CONFIG_DIR/shell_integration.sh"
