package search

import (
	"memex/storage"
	"sort"

	"github.com/sahilm/fuzzy"
)

type Match struct {
	Command storage.Command
	Score   int
	Index   int
}

func Search(query string, commands []storage.Command) []Match {
	if query == "" {
		// Return top commands by frequency if no query
		matches := make([]Match, len(commands))
		for i, cmd := range commands {
			matches[i] = Match{
				Command: cmd,
				Score:   0,
				Index:   i,
			}
		}
		return matches
	}

	// Create string list for fuzzy search
	targets := make([]string, len(commands))
	for i, cmd := range commands {
		targets[i] = cmd.Cmd
	}

	// Perform fuzzy search
	results := fuzzy.Find(query, targets)

	// Map back to Command objects
	matches := make([]Match, len(results))
	for i, r := range results {
		matches[i] = Match{
			Command: commands[r.Index],
			Score:   r.Score,
			Index:   r.Index,
		}
	}

	// Sort by Score (desc), then Frequency (desc)
	sort.Slice(matches, func(i, j int) bool {
		if matches[i].Score == matches[j].Score {
			return matches[i].Command.Frequency > matches[j].Command.Frequency
		}
		return matches[i].Score > matches[j].Score
	})

	return matches
}
