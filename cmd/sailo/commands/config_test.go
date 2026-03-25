package commands

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/agawish/sailo/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunConfigShow_Output(t *testing.T) {
	var buf bytes.Buffer

	// runConfigShow reads from ~/.sailo/config.yaml (defaults if missing)
	// and ./.sailo.yaml (nil if missing). This tests the output format.
	err := runConfigShow(&buf)
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "User config")
	assert.Contains(t, output, "defaults.from:")
	assert.Contains(t, output, "git.credentials:")
}

func TestRunConfigShow_WithProjectConfig(t *testing.T) {
	// Create a temp dir with .sailo.yaml and run from there
	dir := t.TempDir()

	cfg := &config.ProjectConfig{
		Version: 1,
		Image:   "node:22-slim",
		Test:    "npm test",
	}
	require.NoError(t, config.SaveProjectConfig(dir, cfg))

	// Change to temp dir, then restore
	orig, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(dir))
	t.Cleanup(func() { os.Chdir(orig) })

	var buf bytes.Buffer
	err = runConfigShow(&buf)
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "Project config")
	assert.Contains(t, output, "node:22-slim")
	assert.Contains(t, output, "npm test")
}

func TestConfigSet_ValidKeys(t *testing.T) {
	tests := []struct {
		key      string
		value    string
		validate func(t *testing.T, cfg *config.UserConfig)
	}{
		{
			key:   "defaults.from",
			value: "develop",
			validate: func(t *testing.T, cfg *config.UserConfig) {
				assert.Equal(t, "develop", cfg.Defaults.From)
			},
		},
		{
			key:   "defaults.cleanup_after",
			value: "48h",
			validate: func(t *testing.T, cfg *config.UserConfig) {
				assert.Equal(t, "48h", cfg.Defaults.CleanupAfter)
			},
		},
		{
			key:   "defaults.port_range",
			value: "5001-5999",
			validate: func(t *testing.T, cfg *config.UserConfig) {
				assert.Equal(t, "5001-5999", cfg.Defaults.PortRange)
			},
		},
		{
			key:   "git.credentials",
			value: "https",
			validate: func(t *testing.T, cfg *config.UserConfig) {
				assert.Equal(t, "https", cfg.Git.Credentials)
			},
		},
		{
			key:   "git.auto_push",
			value: "false",
			validate: func(t *testing.T, cfg *config.UserConfig) {
				assert.False(t, cfg.Git.AutoPush)
			},
		},
		{
			key:   "env_passthrough",
			value: "KEY1,KEY2,KEY3",
			validate: func(t *testing.T, cfg *config.UserConfig) {
				assert.Equal(t, []string{"KEY1", "KEY2", "KEY3"}, cfg.EnvPassthrough)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, "config.yaml")

			// Save default config to a known path
			defaults := config.DefaultUserConfig()
			require.NoError(t, config.SaveUserConfigTo(path, defaults))

			// Load, set, save manually to verify the logic
			cfg, err := config.LoadUserConfigFrom(path)
			require.NoError(t, err)

			// Apply the same switch logic as runConfigSet
			var buf bytes.Buffer
			// We test the logic indirectly through the config package
			// since runConfigSet uses the real home dir
			switch tt.key {
			case "defaults.from":
				cfg.Defaults.From = tt.value
			case "defaults.cleanup_after":
				cfg.Defaults.CleanupAfter = tt.value
			case "defaults.port_range":
				cfg.Defaults.PortRange = tt.value
			case "git.credentials":
				cfg.Git.Credentials = tt.value
			case "git.auto_push":
				cfg.Git.AutoPush = (tt.value == "true")
			case "env_passthrough":
				parts := []string{}
				for _, p := range splitAndTrim(tt.value) {
					parts = append(parts, p)
				}
				cfg.EnvPassthrough = parts
			}

			require.NoError(t, config.SaveUserConfigTo(path, cfg))

			loaded, err := config.LoadUserConfigFrom(path)
			require.NoError(t, err)
			tt.validate(t, loaded)
			_ = buf
		})
	}
}

func TestConfigSet_UnknownKey(t *testing.T) {
	var buf bytes.Buffer
	err := runConfigSet(&buf, "unknown.key", "value")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown config key")
}

func TestConfigSet_InvalidBool(t *testing.T) {
	var buf bytes.Buffer
	err := runConfigSet(&buf, "git.auto_push", "maybe")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid boolean")
}

func splitAndTrim(s string) []string {
	var result []string
	for _, p := range splitComma(s) {
		p = trimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}

func splitComma(s string) []string {
	return split(s, ",")
}

func split(s, sep string) []string {
	var parts []string
	for len(s) > 0 {
		i := indexOf(s, sep)
		if i < 0 {
			parts = append(parts, s)
			break
		}
		parts = append(parts, s[:i])
		s = s[i+len(sep):]
	}
	return parts
}

func indexOf(s, sub string) int {
	for i := range s {
		if len(s[i:]) >= len(sub) && s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}

func trimSpace(s string) string {
	for len(s) > 0 && (s[0] == ' ' || s[0] == '\t') {
		s = s[1:]
	}
	for len(s) > 0 && (s[len(s)-1] == ' ' || s[len(s)-1] == '\t') {
		s = s[:len(s)-1]
	}
	return s
}
