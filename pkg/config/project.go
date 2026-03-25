// Package config handles sailo configuration from project (.sailo.yaml)
// and user (~/.sailo/config.yaml) configuration files.
package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// ProjectConfig represents the .sailo.yaml project configuration.
type ProjectConfig struct {
	Version  int               `yaml:"version" mapstructure:"version"`
	Image    string            `yaml:"image,omitempty" mapstructure:"image"`
	Services []string          `yaml:"services,omitempty" mapstructure:"services"`
	Ports    map[int]string    `yaml:"ports,omitempty" mapstructure:"ports"`
	Env      map[string]string `yaml:"env,omitempty" mapstructure:"env"`
	Setup    []string          `yaml:"setup,omitempty" mapstructure:"setup"`
	Test     string            `yaml:"test,omitempty" mapstructure:"test"`
}

// LoadProjectConfig loads .sailo.yaml from the given directory.
// Returns nil, nil if the file does not exist (not an error).
func LoadProjectConfig(dir string) (*ProjectConfig, error) {
	path := filepath.Join(dir, ".sailo.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("read project config: %w", err)
	}

	var cfg ProjectConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse project config: %w", err)
	}

	if cfg.Version == 0 {
		cfg.Version = 1
	}
	return &cfg, nil
}

// SaveProjectConfig writes a ProjectConfig to .sailo.yaml in the given directory.
func SaveProjectConfig(dir string, cfg *ProjectConfig) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshal project config: %w", err)
	}

	path := filepath.Join(dir, ".sailo.yaml")
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("write project config: %w", err)
	}
	return nil
}
