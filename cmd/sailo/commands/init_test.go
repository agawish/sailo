package commands

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/agawish/sailo/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunInit_CreatesConfig(t *testing.T) {
	dir := t.TempDir()
	var buf bytes.Buffer

	err := runInit(&buf, dir, false, testutil.NewTestLogger())
	require.NoError(t, err)

	assert.FileExists(t, filepath.Join(dir, ".sailo.yaml"))
	assert.Contains(t, buf.String(), "Initialized .sailo.yaml")
}

func TestRunInit_DetectsLanguage(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module example\n\ngo 1.22\n"), 0644))

	var buf bytes.Buffer
	err := runInit(&buf, dir, false, testutil.NewTestLogger())
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "Language:   go")
	assert.Contains(t, output, "Base image: golang:1.22-bookworm")
}

func TestRunInit_AlreadyExists(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".sailo.yaml"), []byte("version: 1\n"), 0644))

	var buf bytes.Buffer
	err := runInit(&buf, dir, false, testutil.NewTestLogger())
	require.NoError(t, err)

	assert.Contains(t, buf.String(), "already exists")
}

func TestRunInit_ForceOverwrite(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".sailo.yaml"), []byte("version: 1\n"), 0644))

	var buf bytes.Buffer
	err := runInit(&buf, dir, true, testutil.NewTestLogger())
	require.NoError(t, err)

	assert.Contains(t, buf.String(), "Initialized .sailo.yaml")
}

func TestRunInit_DockerfileReused(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "Dockerfile"), []byte("FROM node:22\nEXPOSE 3000\n"), 0644))

	var buf bytes.Buffer
	err := runInit(&buf, dir, false, testutil.NewTestLogger())
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "Dockerfile: found (will be reused)")
	assert.Contains(t, output, "Ports:      3000")

	// Image should NOT be set in config when Dockerfile exists
	data, err := os.ReadFile(filepath.Join(dir, ".sailo.yaml"))
	require.NoError(t, err)
	assert.NotContains(t, string(data), "image:")
}

func TestRunInit_DetectsPorts(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".env"), []byte("PORT=4000\n"), 0644))

	var buf bytes.Buffer
	err := runInit(&buf, dir, false, testutil.NewTestLogger())
	require.NoError(t, err)

	assert.Contains(t, buf.String(), "Ports:      4000")
}
