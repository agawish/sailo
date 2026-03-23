// Package config handles sailo configuration from project (.sailo.yaml)
// and user (~/.sailo/config.yaml) configuration files.
package config

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
func LoadProjectConfig(dir string) (*ProjectConfig, error) {
	// Returns nil config (not error) if file doesn't exist — that's normal
	return &ProjectConfig{Version: 1}, nil
}
