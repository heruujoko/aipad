package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	ConfigFile      = "providers.json"
	AIPadConfigDir  = ".aipad"
	HomeConfigDir   = ".aipad"
	HomeConfigFile  = "providers.json"
)

// CustomProviderConfig defines a custom provider configuration
type CustomProviderConfig struct {
	Name       string `json:"name"`
	ConfigFile string `json:"config_file"`
	RulesDir   string `json:"rules_dir"`
	Enabled    bool   `json:"enabled"`
}

// CustomProviders holds the custom provider configurations
type CustomProviders struct {
	Providers []CustomProviderConfig `json:"providers"`
}

// GetCustomProvidersPath returns the path to the custom providers config file
// It checks both the local project directory and the user's home directory
func GetCustomProvidersPath() (string, error) {
	// First, check local project directory
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	localPath := filepath.Join(cwd, AIPadConfigDir, ConfigFile)

	// Check if local config exists
	if _, err := os.Stat(localPath); err == nil {
		return localPath, nil
	}

	// Fall back to home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	homeConfigPath := filepath.Join(homeDir, HomeConfigDir, HomeConfigFile)

	return homeConfigPath, nil
}

// LoadCustomProviders loads custom provider configurations from disk
func LoadCustomProviders() (*CustomProviders, error) {
	configPath, err := GetCustomProvidersPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get config path: %w", err)
	}

	// If config doesn't exist, return empty providers
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return &CustomProviders{Providers: []CustomProviderConfig{}}, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var customProviders CustomProviders
	if err := json.Unmarshal(data, &customProviders); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &customProviders, nil
}

// SaveCustomProviders saves custom provider configurations to disk
func SaveCustomProviders(customProviders *CustomProviders) error {
	configPath, err := GetCustomProvidersPath()
	if err != nil {
		return fmt.Errorf("failed to get config path: %w", err)
	}

	// Ensure directory exists
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := json.MarshalIndent(customProviders, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// AddProvider adds a new custom provider configuration
func AddProvider(name, configFile, rulesDir string) error {
	customProviders, err := LoadCustomProviders()
	if err != nil {
		return err
	}

	// Check if provider already exists
	for _, p := range customProviders.Providers {
		if p.Name == name {
			return fmt.Errorf("provider '%s' already exists", name)
		}
	}

	// Add new provider
	newProvider := CustomProviderConfig{
		Name:       name,
		ConfigFile: configFile,
		RulesDir:   rulesDir,
		Enabled:    true,
	}
	customProviders.Providers = append(customProviders.Providers, newProvider)

	return SaveCustomProviders(customProviders)
}

// RemoveProvider removes a custom provider configuration
func RemoveProvider(name string) error {
	customProviders, err := LoadCustomProviders()
	if err != nil {
		return err
	}

	// Find and remove the provider
	found := false
	newProviders := make([]CustomProviderConfig, 0, len(customProviders.Providers))
	for _, p := range customProviders.Providers {
		if p.Name != name {
			newProviders = append(newProviders, p)
		} else {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("provider '%s' not found", name)
	}

	customProviders.Providers = newProviders
	return SaveCustomProviders(customProviders)
}

// ListProviders returns all custom provider configurations
func ListProviders() ([]CustomProviderConfig, error) {
	customProviders, err := LoadCustomProviders()
	if err != nil {
		return nil, err
	}
	return customProviders.Providers, nil
}

// GetProvider returns a specific custom provider configuration
func GetProvider(name string) (*CustomProviderConfig, error) {
	customProviders, err := LoadCustomProviders()
	if err != nil {
		return nil, err
	}

	for _, p := range customProviders.Providers {
		if p.Name == name && p.Enabled {
			return &p, nil
		}
	}

	return nil, fmt.Errorf("provider '%s' not found", name)
}

// MergeWithBuiltinProviders merges custom providers with builtin providers
func MergeWithBuiltinProviders(builtinProviders map[string]struct{}) (map[string]bool, error) {
	customProviders, err := LoadCustomProviders()
	if err != nil {
		return nil, err
	}

	// Start with builtin providers
	result := make(map[string]bool)
	for name := range builtinProviders {
		result[name] = true
	}

	// Add custom providers
	for _, p := range customProviders.Providers {
		if p.Enabled {
			result[p.Name] = true
		}
	}

	return result, nil
}

// GetCustomProviderConfigMap returns a map of custom provider configurations
// keyed by provider name
func GetCustomProviderConfigMap() (map[string]CustomProviderConfig, error) {
	customProviders, err := LoadCustomProviders()
	if err != nil {
		return nil, err
	}

	result := make(map[string]CustomProviderConfig)
	for _, p := range customProviders.Providers {
		if p.Enabled {
			result[p.Name] = p
		}
	}

	return result, nil
}
