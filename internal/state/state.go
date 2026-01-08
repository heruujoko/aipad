package state

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

const (
	StateType      = "state.json"
	AIPadDir       = ".aipad"
	ScratchpadFile = "scratchpad.md"
)

type ProviderConfig struct {
	ConfigFile string `json:"config_file"`
	RulesDir   string `json:"rules_dir"`
}

type State struct {
	Version         string                    `json:"version"`
	CurrentProvider string                    `json:"current_provider"`
	SessionID       string                    `json:"session_id"`
	CreatedAt       time.Time                 `json:"created_at"`
	LastSync        time.Time                 `json:"last_sync"`
	ContextHashes   []string                  `json:"context_hashes"`
	ContextHistory  []string                  `json:"context_history"`
	Providers       map[string]ProviderConfig `json:"providers"`
}

// GetStatePath returns the path to the state.json file
func GetStatePath() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(cwd, AIPadDir, StateType), nil
}

// InitAIPadDir creates the .aipad directory if it doesn't exist
func InitAIPadDir() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	aipadPath := filepath.Join(cwd, AIPadDir)
	return os.MkdirAll(aipadPath, 0755)
}

// EnsureScratchpad creates the scratchpad.md file if it doesn't exist
func EnsureScratchpad() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	scratchpadPath := filepath.Join(cwd, AIPadDir, ScratchpadFile)
	if _, err := os.Stat(scratchpadPath); os.IsNotExist(err) {
		_, err := os.Create(scratchpadPath)
		return err
	}
	return nil
}

// NewState creates a default state object
func NewState(provider string) *State {
	return &State{
		Version:         "1.0",
		CurrentProvider: provider,
		SessionID:       uuid.New().String(),
		CreatedAt:       time.Now(),
		LastSync:        time.Now(),
		ContextHashes:   []string{},
		ContextHistory:  []string{},
		Providers: map[string]ProviderConfig{
			"claude": {
				ConfigFile: "CLAUDE.md",
				RulesDir:   ".claude/rules/",
			},
			"antigravity": {
				ConfigFile: "AGENTS.md",
				RulesDir:   ".agent/rules/",
			},
		},
	}
}

// Save writes the state to disk
func (s *State) Save() error {
	path, err := GetStatePath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// Load reads the state from disk
func Load() (*State, error) {
	path, err := GetStatePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var s State
	err = json.Unmarshal(data, &s)
	if err != nil {
		return nil, err
	}

	return &s, nil
}
