# Memex

A full-featured command-line tool with real-time autocomplete, ghost-text suggestions, and fuzzy matching.

## Features

- **Real-time Autocomplete**: Inline ghost-text suggestions as you type.
- **Fuzzy Matching**: Find commands even with partial or typo-ed input.
- **Frequency Ranking**: Suggestions are ordered by usage frequency.
- **Command Storage**: Commands are stored in a simple JSON file in your home directory.
- **Shell Integration**: Seamlessly integrates with Zsh and Bash (Ctrl+G to trigger).

## Installation

### Option 1: Install from Binary (Recommended)

If you have downloaded a release (binary + script):

1.  **Move the binary** to a folder in your PATH (e.g., `~/bin`):
    ```bash
    mkdir -p ~/bin
    mv memex ~/bin/
    chmod +x ~/bin/memex
    ```

2.  **Save the integration script** (e.g., to `~/.memex/`):
    ```bash
    mkdir -p ~/.memex
    mv shell_integration.sh ~/.memex/
    ```

3.  **Enable Shell Integration:**
    Add this to your `~/.zshrc` or `~/.bashrc`:
    ```bash
    source ~/.memex/shell_integration.sh
    ```

### Option 2: Build from Source

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/yourusername/memex.git
    cd memex
    ```

2.  **Run the install script:**
    ```bash
    chmod +x install.sh
    ./install.sh
    ```
    *(This builds the binary and sets everything up automatically)*

3.  **Enable Shell Integration:**
    Add this to your `~/.zshrc` or `~/.bashrc`:
    ```bash
    source ~/repos/projects/memex/shell_integration.sh
    ```

## File Locations

-   **Binary**: `~/bin/memex`
-   **Configuration**: `~/.memex/config.yaml`
-   **Command Database**: `~/.memex/commands.json`
    -   You can manually edit this JSON file to add/remove commands or modify tags.

## Usage

### Interactive Mode (TUI)

1.  **Type & Complete**:
    -   Type a command prefix (e.g., `git`).
    -   Press **Ctrl+G**.
    -   The Memex window opens with suggestions.
    -   Use **Up/Down Arrows** to navigate.
    -   Press **Enter** to select.
    -   The command is pasted into your shell prompt.

### CLI Commands

-   **Add a command:**
    ```bash
    memex --add "kubectl get pods" --tags "k8s"
    ```
    *Tip: You can also just type a new command in the TUI and press Enter; it will be saved automatically.*

-   **List all commands:**
    ```bash
    memex --list
    ```

## Configuration

Edit `~/.memex/config.yaml` to customize:

```yaml
ghost_text: true       # Enable/disable ghost text
fuzzy_search: true     # Enable/disable fuzzy matching
theme: "dark"          # UI theme
suggestion_limit: 10   # Max suggestions to show
```
