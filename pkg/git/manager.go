// Package git manages git operations for workspace isolation.
//
// Each workspace gets a full shallow clone (git clone --depth=1) into a
// Docker volume, ensuring complete isolation from the host's git state.
//
// All git operations (clone, branch, diff, commit, push) run inside
// the workspace container via docker exec. GetRemoteURL runs on the host.
package git

import (
	"context"
	"fmt"
	"log/slog"
	"os/exec"
	"strings"
)

// WorkspaceDir is the standard path where repos are cloned inside containers.
const WorkspaceDir = "/workspace"

// ContainerExecer is the interface needed to run commands inside a container.
type ContainerExecer interface {
	ExecInContainer(ctx context.Context, containerID string, cmd []string) error
	ExecWithOutput(ctx context.Context, containerID string, cmd []string) (string, error)
}

// Manager handles git operations for workspaces.
type Manager struct {
	container ContainerExecer
	logger    *slog.Logger
}

// NewManager creates a git manager.
func NewManager(container ContainerExecer, logger *slog.Logger) *Manager {
	return &Manager{container: container, logger: logger}
}

// Clone performs a shallow clone of the repository into the workspace container.
func (m *Manager) Clone(ctx context.Context, containerID string, repoURL string, branch string, targetDir string) error {
	cmd := []string{"git", "clone", "--depth=1", "--branch", branch, repoURL, targetDir}
	if err := m.container.ExecInContainer(ctx, containerID, cmd); err != nil {
		return fmt.Errorf("git clone: %w", err)
	}
	m.logger.Info("cloned repository", "repo", repoURL, "branch", branch, "dir", targetDir)
	return nil
}

// CreateBranch creates and checks out a new branch in the workspace container.
func (m *Manager) CreateBranch(ctx context.Context, containerID string, workspaceDir string, branchName string) error {
	cmd := []string{"git", "-C", workspaceDir, "checkout", "-b", branchName}
	if err := m.container.ExecInContainer(ctx, containerID, cmd); err != nil {
		return fmt.Errorf("create branch %s: %w", branchName, err)
	}
	m.logger.Info("created branch", "branch", branchName)
	return nil
}

// Diff returns the git diff of changes in the workspace container.
func (m *Manager) Diff(ctx context.Context, containerID string, workspaceDir string, statOnly bool) (string, error) {
	cmd := []string{"git", "-C", workspaceDir, "diff"}
	if statOnly {
		cmd = append(cmd, "--stat")
	}
	output, err := m.container.ExecWithOutput(ctx, containerID, cmd)
	if err != nil {
		return "", fmt.Errorf("git diff: %w", err)
	}
	return output, nil
}

// CommitAll stages and commits all changes with the given message.
func (m *Manager) CommitAll(ctx context.Context, containerID string, workspaceDir string, message string) error {
	addCmd := []string{"git", "-C", workspaceDir, "add", "-A"}
	if err := m.container.ExecInContainer(ctx, containerID, addCmd); err != nil {
		return fmt.Errorf("git add: %w", err)
	}
	commitCmd := []string{"git", "-C", workspaceDir, "commit", "-m", message}
	if err := m.container.ExecInContainer(ctx, containerID, commitCmd); err != nil {
		return fmt.Errorf("git commit: %w", err)
	}
	return nil
}

// Push pushes the current branch to the remote.
func (m *Manager) Push(ctx context.Context, containerID string, workspaceDir string) error {
	cmd := []string{"git", "-C", workspaceDir, "push", "-u", "origin", "HEAD"}
	if err := m.container.ExecInContainer(ctx, containerID, cmd); err != nil {
		return fmt.Errorf("git push: %w", err)
	}
	return nil
}

// GetRemoteURL returns the remote origin URL for the current directory.
// This runs on the HOST, not inside a container.
func GetRemoteURL() (string, error) {
	cmd := exec.Command("git", "remote", "get-url", "origin")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("get git remote URL: %w (are you in a git repository?)", err)
	}
	return strings.TrimSpace(string(out)), nil
}
