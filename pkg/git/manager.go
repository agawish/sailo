// Package git manages git operations for workspace isolation.
//
// Each workspace gets a full shallow clone (git clone --depth=1) into a
// Docker volume, ensuring complete isolation from the host's git state.
package git

import (
	"context"
	"fmt"
	"log/slog"
)

// Manager handles git operations for workspaces.
type Manager struct {
	logger *slog.Logger
}

// NewManager creates a git manager.
func NewManager(logger *slog.Logger) *Manager {
	return &Manager{logger: logger}
}

// Clone performs a shallow clone of the repository into the workspace.
func (m *Manager) Clone(ctx context.Context, repoURL string, branch string, targetDir string) error {
	return fmt.Errorf("git clone not yet implemented")
}

// CreateBranch creates a new branch in the workspace.
func (m *Manager) CreateBranch(ctx context.Context, workspaceDir string, branchName string) error {
	return fmt.Errorf("git branch creation not yet implemented")
}

// Diff returns the git diff of changes in the workspace.
func (m *Manager) Diff(ctx context.Context, workspaceDir string, statOnly bool) (string, error) {
	return "", fmt.Errorf("git diff not yet implemented")
}

// CommitAll stages and commits all changes with the given message.
func (m *Manager) CommitAll(ctx context.Context, workspaceDir string, message string) error {
	return fmt.Errorf("git commit not yet implemented")
}

// Push pushes the current branch to the remote.
func (m *Manager) Push(ctx context.Context, workspaceDir string) error {
	return fmt.Errorf("git push not yet implemented")
}

// GetRemoteURL returns the remote origin URL for the current directory.
func GetRemoteURL() (string, error) {
	return "", fmt.Errorf("git remote URL detection not yet implemented")
}
