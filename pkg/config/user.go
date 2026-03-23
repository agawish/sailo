package config

// UserConfig represents the ~/.sailo/config.yaml user configuration.
type UserConfig struct {
	Defaults       Defaults          `yaml:"defaults" mapstructure:"defaults"`
	EnvPassthrough []string          `yaml:"env_passthrough" mapstructure:"env_passthrough"`
	Agents         map[string]Agent  `yaml:"agents,omitempty" mapstructure:"agents"`
	Git            GitConfig         `yaml:"git" mapstructure:"git"`
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

// LoadUserConfig loads ~/.sailo/config.yaml.
func LoadUserConfig() (*UserConfig, error) {
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
	}, nil
}
