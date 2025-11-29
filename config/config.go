package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	GhostText       bool   `yaml:"ghost_text"`
	FuzzySearch     bool   `yaml:"fuzzy_search"`
	Theme           string `yaml:"theme"`
	SuggestionLimit int    `yaml:"suggestion_limit"`
	CommandFile     string `yaml:"command_file"`
}

func DefaultConfig() Config {
	home, _ := os.UserHomeDir()
	return Config{
		GhostText:       true,
		FuzzySearch:     true,
		Theme:           "dark",
		SuggestionLimit: 10,
		CommandFile:     filepath.Join(home, ".memex", "commands.json"),
	}
}

func LoadConfig() (Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return DefaultConfig(), err
	}

	configPath := filepath.Join(home, ".memex", "config.yaml")
	data, err := os.ReadFile(configPath)
	if os.IsNotExist(err) {
		return DefaultConfig(), nil
	}
	if err != nil {
		return DefaultConfig(), err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return DefaultConfig(), err
	}

	// Ensure defaults if fields are missing
	if cfg.SuggestionLimit == 0 {
		cfg.SuggestionLimit = 10
	}
	if cfg.CommandFile == "" {
		cfg.CommandFile = filepath.Join(home, ".memex", "commands.json")
	}

	return cfg, nil
}
