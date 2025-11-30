#!/bin/bash

set -e

REPO_URL="https://github.com/jai0651/memex.git"
INSTALL_DIR="$HOME/bin"
CONFIG_DIR="$HOME/.memex"

echo "Installing Memex..."

# Check for Go
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed. Please install Go first."
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
go build -o memex main.go

# Create local bin directory if it doesn't exist
mkdir -p "$INSTALL_DIR"

# Move binary to $HOME/bin
echo "Moving binary to $INSTALL_DIR..."
mv memex "$INSTALL_DIR/"

# Add to PATH if not already there (temporary for this session)
export PATH=$INSTALL_DIR:$PATH

# Create config directory
mkdir -p "$CONFIG_DIR"

# Copy shell integration script
echo "Installing shell integration..."
cp shell_integration.sh "$CONFIG_DIR/"

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

# Cleanup if remote
if [ "$IS_REMOTE" = true ]; then
    echo "Cleaning up..."
    rm -rf "$TEMP_DIR"
fi

echo "Installation complete!"
echo ""
echo "To enable shell integration, add the following to your .zshrc or .bashrc:"
echo "source $CONFIG_DIR/shell_integration.sh"
