#!/bin/bash

set -e

echo "Installing Memex..."

# Build the binary
go build -o memex main.go

# Create local bin directory if it doesn't exist
mkdir -p $HOME/bin

# Move binary to $HOME/bin
echo "Moving binary to $HOME/bin..."
mv memex $HOME/bin/

# Add to PATH if not already there (temporary for this session)
export PATH=$HOME/bin:$PATH

# Create config directory
mkdir -p ~/.memex

# Create default config if not exists
if [ ! -f ~/.memex/config.yaml ]; then
    echo "Creating default config..."
    cat > ~/.memex/config.yaml <<EOF
ghost_text: true
fuzzy_search: true
theme: "dark"
suggestion_limit: 10
command_file: "$HOME/.memex/commands.json"
EOF
fi

# Create default commands file if not exists
if [ ! -f ~/.memex/commands.json ]; then
    echo "Creating default commands..."
    cat > ~/.memex/commands.json <<EOF
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

echo "Installation complete!"
echo ""
echo "To enable shell integration, add the following to your .zshrc or .bashrc:"
echo "source $(pwd)/shell_integration.sh"
