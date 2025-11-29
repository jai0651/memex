# Memex Integration

# Detect Shell and apply appropriate integration

if [ -n "$ZSH_VERSION" ]; then
    # --- ZSH Integration ---
    smart_cli_run_zsh() {
        local output_file=$(mktemp)
        # Pass current buffer ($BUFFER) to memex
        # Use -- to ensure BUFFER is treated as an argument, not a flag
        $HOME/bin/memex --out "$output_file" -- "$BUFFER" < /dev/tty
        
        if [ -f "$output_file" ]; then
            local cmd=$(cat "$output_file")
            rm "$output_file"
            
            if [ -n "$cmd" ]; then
                # Replace buffer with selected command
                BUFFER="$cmd"
                CURSOR=${#BUFFER}
            fi
        fi
        zle redisplay
    }

    # Bind Ctrl+G to smart_cli_run_zsh
    zle -N smart_cli_run_zsh
    bindkey '^g' smart_cli_run_zsh

elif [ -n "$BASH_VERSION" ]; then
    # --- BASH Integration ---
    smart_cli_run_bash() {
        local output_file=$(mktemp)
        # Pass current buffer ($READLINE_LINE) to memex
        # Use -- to ensure READLINE_LINE is treated as an argument, not a flag
        $HOME/bin/memex --out "$output_file" -- "$READLINE_LINE" < /dev/tty
        
        if [ -f "$output_file" ]; then
            local cmd=$(cat "$output_file")
            rm "$output_file"
            
            if [ -n "$cmd" ]; then
                READLINE_LINE="$cmd"
                READLINE_POINT=${#cmd}
            fi
        fi
    }

    # Bind Ctrl+G to smart_cli_run_bash
    bind -x '"\C-g": smart_cli_run_bash'
fi

echo "Memex integration loaded. Press Ctrl+G to use."
