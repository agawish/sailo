// Package creds manages credential forwarding into workspace containers.
//
// Security model:
//   - SSH keys are forwarded via SSH_AUTH_SOCK (never copied to disk)
//   - Environment variables are explicitly allowlisted (no blanket passthrough)
//   - gh CLI config is mounted read-only
package creds

import (
	"fmt"
	"os"
)

// SSHAgentSocket returns the host's SSH agent socket path.
// Returns an error if SSH_AUTH_SOCK is not set.
func SSHAgentSocket() (string, error) {
	sock := os.Getenv("SSH_AUTH_SOCK")
	if sock == "" {
		return "", fmt.Errorf("SSH_AUTH_SOCK is not set; start an SSH agent or run `ssh-agent`")
	}
	return sock, nil
}

// GHConfigDir returns the path to the gh CLI config directory.
func GHConfigDir() string {
	if dir := os.Getenv("GH_CONFIG_DIR"); dir != "" {
		return dir
	}
	home, _ := os.UserHomeDir()
	return home + "/.config/gh"
}
