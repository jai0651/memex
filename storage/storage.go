package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

type Command struct {
	Cmd       string   `json:"cmd"`
	Tags      []string `json:"tags"`
	Frequency int      `json:"frequency"`
	LastUsed  int64    `json:"last_used"`
}

type Storage struct {
	Commands []Command
	FilePath string
	mu       sync.Mutex
}

func NewStorage(filePath string) *Storage {
	return &Storage{
		FilePath: filePath,
		Commands: []Command{},
	}
}

func (s *Storage) Load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(s.FilePath)
	if os.IsNotExist(err) {
		s.Commands = []Command{}
		return nil
	}
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &s.Commands)
}

func (s *Storage) Save() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	dir := filepath.Dir(s.FilePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(s.Commands, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.FilePath, data, 0644)
}

func (s *Storage) AddCommand(cmdStr string, tags []string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if exists
	for i, c := range s.Commands {
		if c.Cmd == cmdStr {
			s.Commands[i].Frequency++
			s.Commands[i].LastUsed = time.Now().Unix()
			return
		}
	}

	// Add new
	s.Commands = append(s.Commands, Command{
		Cmd:       cmdStr,
		Tags:      tags,
		Frequency: 1,
		LastUsed:  time.Now().Unix(),
	})
}

func (s *Storage) GetCommands() []Command {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Return a copy to avoid race conditions if modified elsewhere
	cmds := make([]Command, len(s.Commands))
	copy(cmds, s.Commands)

	// Sort by Frequency desc, then LastUsed desc
	sort.Slice(cmds, func(i, j int) bool {
		if cmds[i].Frequency == cmds[j].Frequency {
			return cmds[i].LastUsed > cmds[j].LastUsed
		}
		return cmds[i].Frequency > cmds[j].Frequency
	})

	return cmds
}
