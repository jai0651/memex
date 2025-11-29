package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"memex/config"
	"memex/storage"
	"memex/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Flags
	addFlag := flag.String("add", "", "Add a new command")
	tagsFlag := flag.String("tags", "", "Tags for the new command (comma separated)")
	outFlag := flag.String("out", "", "Output file to write selected command to")
	listFlag := flag.Bool("list", false, "List all stored commands")
	flag.Parse()

	// Load Config & Storage
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	store := storage.NewStorage(cfg.CommandFile)
	if err := store.Load(); err != nil {
		// If file doesn't exist, it will be created on save
		// If other error, warn but continue
		if !os.IsNotExist(err) {
			fmt.Printf("Warning: could not load commands: %v\n", err)
		}
	}

	// Handle --list
	if *listFlag {
		cmds := store.GetCommands()
		if len(cmds) == 0 {
			fmt.Println("No commands stored.")
			return
		}
		fmt.Printf("%-30s %-20s %s\n", "COMMAND", "TAGS", "FREQUENCY")
		fmt.Println(strings.Repeat("-", 60))
		for _, c := range cmds {
			tags := strings.Join(c.Tags, ", ")
			fmt.Printf("%-30s %-20s %d\n", c.Cmd, tags, c.Frequency)
		}
		return
	}

	// Handle --add
	if *addFlag != "" {
		tags := []string{}
		if *tagsFlag != "" {
			tags = strings.Split(*tagsFlag, ",")
			for i := range tags {
				tags[i] = strings.TrimSpace(tags[i])
			}
		}
		store.AddCommand(*addFlag, tags)
		if err := store.Save(); err != nil {
			fmt.Printf("Error saving command: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Command added.")
		return
	}

	// Parse remaining args as query if no flag
	initialQuery := ""
	if len(flag.Args()) > 0 {
		initialQuery = strings.Join(flag.Args(), " ")
	}

	// TUI Mode
	p := tea.NewProgram(ui.NewModel(cfg, store, initialQuery), tea.WithAltScreen())
	m, err := p.Run()
	if err != nil {
		fmt.Printf("Error running TUI: %v\n", err)
		os.Exit(1)
	}

	// Handle Output
	finalModel := m.(ui.Model)
	if finalModel.OutputCmd != "" {
		// If output file specified, write to it
		if *outFlag != "" {
			if err := os.WriteFile(*outFlag, []byte(finalModel.OutputCmd), 0644); err != nil {
				fmt.Printf("Error writing output: %v\n", err)
				os.Exit(1)
			}
		} else {
			// Otherwise just print to stdout (useful for piping)
			fmt.Println(finalModel.OutputCmd)
		}
	}
}
