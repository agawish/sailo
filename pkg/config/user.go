package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// UserConfig represents the ~/.sailo/config.yaml user configuration.
type UserConfig struct {
	Defaults       Defaults         `yaml:"defaults" mapstructure:"defaults"`
	EnvPassthrough []string         `yaml:"env_passthrough" mapstructure:"env_passthrough"`
	Agents         map[string]Agent `yaml:"agents,omitempty" mapstructure:"agents"`
	Git            GitConfig        `yaml:"git" mapstructure:"git"`
}

// Defaults holds default values for workspace creation.
type Defaults struct {
	From         string `yaml:"from" mapstructure:"from"`
	CleanupAfter string `yaml:"cleanup_after" mapstructure:"cleanup_after"`
	PortRange    string `yaml:"port_range" mapstructure:"port_range"`
}

// Agent defines an optional agent shortcut configuration.
type Agent struct {
	Command string   `yaml:"command" mapstructure:"command"`
	Mount   []string `yaml:"mount,omitempty" mapstructure:"mount"`
}

// GitConfig holds git-related settings.
type GitConfig struct {
	Credentials string `yaml:"credentials" mapstructure:"credentials"`
	AutoPush    bool   `yaml:"auto_push" mapstructure:"auto_push"`
}

// DefaultUserConfig returns the built-in default user configuration.
func DefaultUserConfig() *UserConfig {
	return &UserConfig{
		Defaults: Defaults{
			From:         "main",
			CleanupAfter: "24h",
			PortRange:    "3001-3999",
		},
		EnvPassthrough: []string{
			"ANTHROPIC_API_KEY",
			"OPENAI_API_KEY",
			"GITHUB_TOKEN",
		},
		Git: GitConfig{
			Credentials: "ssh-agent",
			AutoPush:    true,
		},
	}
}

// LoadUserConfig loads ~/.sailo/config.yaml.
// Returns default config if the file does not exist.
func LoadUserConfig() (*UserConfig, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("get home directory: %w", err)
	}

	path := filepath.Join(home, ".sailo", "config.yaml")
	return LoadUserConfigFrom(path)
}

// LoadUserConfigFrom loads user config from a specific path.
// Returns default config if the file does not exist.
func LoadUserConfigFrom(path string) (*UserConfig, error) {
	defaults := DefaultUserConfig()

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return defaults, nil
		}
		return nil, fmt.Errorf("read user config: %w", err)
	}

	var cfg UserConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse user config: %w", err)
	}

	// Backfill zero-value fields with defaults
	if cfg.Defaults.From == "" {
		cfg.Defaults.From = defaults.Defaults.From
	}
	if cfg.Defaults.CleanupAfter == "" {
		cfg.Defaults.CleanupAfter = defaults.Defaults.CleanupAfter
	}
	if cfg.Defaults.PortRange == "" {
		cfg.Defaults.PortRange = defaults.Defaults.PortRange
	}
	if cfg.Git.Credentials == "" {
		cfg.Git.Credentials = defaults.Git.Credentials
	}
	if len(cfg.EnvPassthrough) == 0 {
		cfg.EnvPassthrough = defaults.EnvPassthrough
	}

	return &cfg, nil
}

// SaveUserConfig writes the user config to ~/.sailo/config.yaml.
func SaveUserConfig(cfg *UserConfig) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("get home directory: %w", err)
	}

	path := filepath.Join(home, ".sailo", "config.yaml")
	return SaveUserConfigTo(path, cfg)
}

// SaveUserConfigTo writes user config to a specific path.
func SaveUserConfigTo(path string, cfg *UserConfig) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create config directory: %w", err)
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshal user config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("write user config: %w", err)
	}
	return nil
}
