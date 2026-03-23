// Package container wraps the Moby SDK to manage Docker containers
// for sailo workspaces.
package container

import (
	"context"
	"fmt"
	"log/slog"
)

// Client wraps the Docker/Moby API client.
type Client struct {
	logger *slog.Logger
}

// NewClient creates a Docker client connected to the local daemon.
func NewClient(logger *slog.Logger) (*Client, error) {
	return nil, fmt.Errorf("container client not yet implemented")
}

// Ping checks that the Docker daemon is reachable.
func (c *Client) Ping(ctx context.Context) error {
	return fmt.Errorf("container ping not yet implemented")
}

// CreateWorkspace creates and starts a container for a workspace.
func (c *Client) CreateWorkspace(ctx context.Context, opts CreateOptions) (string, error) {
	return "", fmt.Errorf("container creation not yet implemented")
}

// StopContainer stops a running container.
func (c *Client) StopContainer(ctx context.Context, containerID string) error {
	return fmt.Errorf("container stop not yet implemented")
}

// StartContainer starts a stopped container.
func (c *Client) StartContainer(ctx context.Context, containerID string) error {
	return fmt.Errorf("container start not yet implemented")
}

// RemoveContainer removes a container and its volumes.
func (c *Client) RemoveContainer(ctx context.Context, containerID string) error {
	return fmt.Errorf("container remove not yet implemented")
}

// ExecInContainer runs a command inside a running container.
func (c *Client) ExecInContainer(ctx context.Context, containerID string, cmd []string) error {
	return fmt.Errorf("container exec not yet implemented")
}

// CreateOptions configures a new workspace container.
type CreateOptions struct {
	Image       string
	Ports       map[int]int // container port → host port
	SSHAuthSock string
	EnvVars     map[string]string
	GHConfigDir string // path to gh config (mounted read-only)
}
