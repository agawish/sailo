package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- Project Config Tests ---

func TestLoadProjectConfig_FileExists(t *testing.T) {
	dir := t.TempDir()
	content := `version: 1
image: node:22-slim
ports:
  3000: auto
  8080: auto
env:
  NODE_ENV: development
setup:
  - npm install
test: npm test
`
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".sailo.yaml"), []byte(content), 0644))

	cfg, err := LoadProjectConfig(dir)
	require.NoError(t, err)
	require.NotNil(t, cfg)

	assert.Equal(t, 1, cfg.Version)
	assert.Equal(t, "node:22-slim", cfg.Image)
	assert.Equal(t, map[int]string{3000: "auto", 8080: "auto"}, cfg.Ports)
	assert.Equal(t, "development", cfg.Env["NODE_ENV"])
	assert.Equal(t, []string{"npm install"}, cfg.Setup)
	assert.Equal(t, "npm test", cfg.Test)
}

func TestLoadProjectConfig_FileMissing(t *testing.T) {
	dir := t.TempDir()

	cfg, err := LoadProjectConfig(dir)
	assert.NoError(t, err)
	assert.Nil(t, cfg)
}

func TestLoadProjectConfig_InvalidYAML(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".sailo.yaml"), []byte(":::bad yaml"), 0644))

	cfg, err := LoadProjectConfig(dir)
	assert.Error(t, err)
	assert.Nil(t, cfg)
	assert.Contains(t, err.Error(), "parse project config")
}

func TestLoadProjectConfig_DefaultVersion(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".sailo.yaml"), []byte("image: ubuntu:24.04\n"), 0644))

	cfg, err := LoadProjectConfig(dir)
	require.NoError(t, err)
	assert.Equal(t, 1, cfg.Version)
}

func TestSaveProjectConfig_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	original := &ProjectConfig{
		Version: 1,
		Image:   "golang:1.22-bookworm",
		Ports:   map[int]string{8080: "auto"},
		Setup:   []string{"go mod download"},
		Test:    "go test ./...",
	}

	require.NoError(t, SaveProjectConfig(dir, original))
	assert.FileExists(t, filepath.Join(dir, ".sailo.yaml"))

	loaded, err := LoadProjectConfig(dir)
	require.NoError(t, err)
	assert.Equal(t, original.Version, loaded.Version)
	assert.Equal(t, original.Image, loaded.Image)
	assert.Equal(t, original.Ports, loaded.Ports)
	assert.Equal(t, original.Setup, loaded.Setup)
	assert.Equal(t, original.Test, loaded.Test)
}

// --- User Config Tests ---

func TestLoadUserConfig_FileMissing(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nonexistent", "config.yaml")

	cfg, err := LoadUserConfigFrom(path)
	require.NoError(t, err)
	require.NotNil(t, cfg)

	// Should return defaults
	assert.Equal(t, "main", cfg.Defaults.From)
	assert.Equal(t, "24h", cfg.Defaults.CleanupAfter)
	assert.Equal(t, "3001-3999", cfg.Defaults.PortRange)
	assert.Equal(t, "ssh-agent", cfg.Git.Credentials)
	assert.True(t, cfg.Git.AutoPush)
	assert.Contains(t, cfg.EnvPassthrough, "ANTHROPIC_API_KEY")
}

func TestLoadUserConfig_FileExists(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	content := `defaults:
  from: develop
  cleanup_after: 48h
  port_range: 4001-4999
env_passthrough:
  - MY_SECRET
git:
  credentials: https
  auto_push: false
`
	require.NoError(t, os.WriteFile(path, []byte(content), 0644))

	cfg, err := LoadUserConfigFrom(path)
	require.NoError(t, err)

	assert.Equal(t, "develop", cfg.Defaults.From)
	assert.Equal(t, "48h", cfg.Defaults.CleanupAfter)
	assert.Equal(t, "4001-4999", cfg.Defaults.PortRange)
	assert.Equal(t, []string{"MY_SECRET"}, cfg.EnvPassthrough)
	assert.Equal(t, "https", cfg.Git.Credentials)
	assert.False(t, cfg.Git.AutoPush)
}

func TestLoadUserConfig_DefaultsMerged(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	// Only set one field; others should get defaults
	content := `defaults:
  from: develop
`
	require.NoError(t, os.WriteFile(path, []byte(content), 0644))

	cfg, err := LoadUserConfigFrom(path)
	require.NoError(t, err)

	assert.Equal(t, "develop", cfg.Defaults.From)
	assert.Equal(t, "24h", cfg.Defaults.CleanupAfter)
	assert.Equal(t, "3001-3999", cfg.Defaults.PortRange)
	assert.Equal(t, "ssh-agent", cfg.Git.Credentials)
	assert.Contains(t, cfg.EnvPassthrough, "ANTHROPIC_API_KEY")
}

func TestLoadUserConfig_InvalidYAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	require.NoError(t, os.WriteFile(path, []byte(":::bad"), 0644))

	_, err := LoadUserConfigFrom(path)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "parse user config")
}

func TestSaveUserConfig_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".sailo", "config.yaml")

	original := &UserConfig{
		Defaults: Defaults{
			From:         "develop",
			CleanupAfter: "12h",
			PortRange:    "5001-5999",
		},
		EnvPassthrough: []string{"MY_KEY"},
		Git: GitConfig{
			Credentials: "https",
			AutoPush:    false,
		},
	}

	require.NoError(t, SaveUserConfigTo(path, original))
	assert.FileExists(t, path)

	loaded, err := LoadUserConfigFrom(path)
	require.NoError(t, err)
	assert.Equal(t, original.Defaults, loaded.Defaults)
	assert.Equal(t, original.EnvPassthrough, loaded.EnvPassthrough)
	assert.Equal(t, original.Git, loaded.Git)
}
