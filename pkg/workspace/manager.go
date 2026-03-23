// Package workspace manages the lifecycle of isolated agent workspaces.
//
// A workspace is a Docker container + git clone + port range that provides
// full isolation for an AI coding agent. The Manager orchestrates creation,
// state transitions, and cleanup.
package workspace

import (
	"context"
	"fmt"
	"log/slog"
)

// Manager orchestrates workspace lifecycle operations.
type Manager struct {
	store  *Store
	logger *slog.Logger
}

// NewManager creates a workspace manager with the given store and logger.
func NewManager(store *Store, logger *slog.Logger) *Manager {
	return &Manager{
		store:  store,
		logger: logger,
	}
}

// Create provisions a new isolated workspace.
func (m *Manager) Create(ctx context.Context, task string, fromBranch string) (*Workspace, error) {
	return nil, fmt.Errorf("workspace creation not yet implemented")
}

// Stop halts a running workspace, preserving its state.
func (m *Manager) Stop(ctx context.Context, id string) error {
	return fmt.Errorf("workspace stop not yet implemented")
}

// Start resumes a stopped workspace.
func (m *Manager) Start(ctx context.Context, id string) error {
	return fmt.Errorf("workspace start not yet implemented")
}

// Remove destroys a workspace and frees its resources.
func (m *Manager) Remove(ctx context.Context, id string) error {
	return fmt.Errorf("workspace remove not yet implemented")
}

// List returns all workspaces matching the given filter.
func (m *Manager) List(ctx context.Context, includeArchived bool) ([]Workspace, error) {
	return nil, fmt.Errorf("workspace list not yet implemented")
}

// Get returns a single workspace by ID.
func (m *Manager) Get(ctx context.Context, id string) (*Workspace, error) {
	return nil, fmt.Errorf("workspace get not yet implemented")
}
