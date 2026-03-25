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

	"github.com/agawish/sailo/pkg/config"
	"github.com/agawish/sailo/pkg/container"
	"github.com/agawish/sailo/pkg/creds"
	"github.com/agawish/sailo/pkg/git"
)

// getRemoteURL is the function used to detect the git remote URL.
// It's a package variable so tests can override it.
var getRemoteURL = git.GetRemoteURL

// ManagerConfig groups all dependencies for the workspace Manager.
type ManagerConfig struct {
	Store      *Store
	Container  ContainerClient
	Ports      PortAllocator
	Git        GitOperator
	Detector   ProjectDetector
	UserConfig *config.UserConfig
	Logger     *slog.Logger
}

// Manager orchestrates workspace lifecycle operations.
type Manager struct {
	store      *Store
	container  ContainerClient
	ports      PortAllocator
	git        GitOperator
	detector   ProjectDetector
	userConfig *config.UserConfig
	logger     *slog.Logger
}

// NewManager creates a workspace manager with all dependencies.
func NewManager(cfg ManagerConfig) *Manager {
	return &Manager{
		store:      cfg.Store,
		container:  cfg.Container,
		ports:      cfg.Ports,
		git:        cfg.Git,
		detector:   cfg.Detector,
		userConfig: cfg.UserConfig,
		logger:     cfg.Logger,
	}
}

// CreateOptions configures workspace creation.
type CreateOptions struct {
	Task       string
	FromBranch string
	Image      string // override, empty = auto-detect
}

// Create provisions a new isolated workspace.
//
//	Flow: Ping → ID gen → detect → remote URL → allocate ports →
//	      save record → create container → install git → clone →
//	      branch → setup → transition to running
func (m *Manager) Create(ctx context.Context, opts CreateOptions) (*Workspace, error) {
	// 1. Validate Docker is available
	if err := m.container.Ping(ctx); err != nil {
		return nil, err
	}

	// 2. Generate workspace ID and branch name
	wsID := generateID()
	branch := fmt.Sprintf("sailo/%s/%s", wsID, slugify(opts.Task))

	// 3. Detect project configuration
	projectCfg, _ := config.LoadProjectConfig(".")
	detection, err := m.detector.Detect(".")
	if err != nil {
		return nil, fmt.Errorf("detect project: %w", err)
	}

	// 4. Determine base image
	img := detection.BaseImage
	if projectCfg != nil && projectCfg.Image != "" {
		img = projectCfg.Image
	}
	if opts.Image != "" {
		img = opts.Image
	}

	// 5. Detect remote URL (on host)
	repoURL, err := getRemoteURL()
	if err != nil {
		return nil, fmt.Errorf("detect git remote: %w", err)
	}

	// 6. Allocate host ports
	containerPorts := detection.Ports
	if projectCfg != nil && len(projectCfg.Ports) > 0 {
		containerPorts = nil
		for p := range projectCfg.Ports {
			containerPorts = append(containerPorts, p)
		}
	}
	portMap := make(map[int]int)
	for _, cp := range containerPorts {
		hp, err := m.ports.Allocate(cp)
		if err != nil {
			return nil, fmt.Errorf("allocate port for %d: %w", cp, err)
		}
		portMap[cp] = hp
	}

	// 7. Save workspace record (state: creating)
	ws := &Workspace{
		ID:         wsID,
		Task:       opts.Task,
		State:      StateCreating,
		Branch:     branch,
		Ports:      portMap,
		FromBranch: opts.FromBranch,
	}
	if err := m.store.Save(ws); err != nil {
		return nil, fmt.Errorf("save workspace: %w", err)
	}

	// cleanup removes container and workspace record on failure
	var containerID string
	cleanup := func() {
		if containerID != "" {
			m.container.RemoveContainer(ctx, containerID)
		}
		m.store.Delete(wsID)
	}

	// 8. Create and start container
	sshSock, _ := creds.SSHAgentSocket()
	ghConfig := creds.GHConfigDir()
	envVars := creds.FilterEnv(m.userConfig.EnvPassthrough)
	if sshSock != "" {
		envVars["SSH_AUTH_SOCK"] = "/run/ssh-agent.sock"
	}

	containerID, err = m.container.CreateWorkspace(ctx, container.CreateOptions{
		WorkspaceID: wsID,
		Image:       img,
		Ports:       portMap,
		SSHAuthSock: sshSock,
		EnvVars:     envVars,
		GHConfigDir: ghConfig,
	})
	if err != nil {
		cleanup()
		return nil, fmt.Errorf("create container: %w", err)
	}
	ws.ContainerID = containerID

	// 9. Install git (check-then-install)
	if err := m.container.ExecInContainer(ctx, containerID, []string{
		"sh", "-c", "which git > /dev/null 2>&1 || (apt-get update && apt-get install -y --no-install-recommends git openssh-client && rm -rf /var/lib/apt/lists/*)",
	}); err != nil {
		cleanup()
		return nil, fmt.Errorf("install git in container: %w", err)
	}

	// 10. SSH keyscan (non-fatal)
	m.container.ExecInContainer(ctx, containerID, []string{
		"sh", "-c", "mkdir -p /root/.ssh && ssh-keyscan github.com >> /root/.ssh/known_hosts 2>/dev/null",
	})

	// 11. Clone repository
	if err := m.git.Clone(ctx, containerID, repoURL, opts.FromBranch, git.WorkspaceDir); err != nil {
		cleanup()
		return nil, fmt.Errorf("clone repository: %w", err)
	}

	// 12. Create workspace branch
	if err := m.git.CreateBranch(ctx, containerID, git.WorkspaceDir, branch); err != nil {
		cleanup()
		return nil, fmt.Errorf("create branch: %w", err)
	}

	// 13. Run setup commands (non-fatal on failure)
	if projectCfg != nil {
		for _, cmd := range projectCfg.Setup {
			m.logger.Info("running setup command", "cmd", cmd)
			if err := m.container.ExecInContainer(ctx, containerID, []string{
				"sh", "-c", "cd /workspace && " + cmd,
			}); err != nil {
				m.logger.Warn("setup command failed (continuing)", "cmd", cmd, "error", err)
			}
		}
	}

	// 14. Transition to running
	if err := ws.Transition(StateRunning); err != nil {
		cleanup()
		return nil, fmt.Errorf("transition to running: %w", err)
	}
	if err := m.store.Save(ws); err != nil {
		return nil, fmt.Errorf("save workspace: %w", err)
	}

	m.logger.Info("workspace created",
		"id", wsID,
		"branch", branch,
		"container", containerID[:12],
		"ports", portMap,
	)
	return ws, nil
}

// Stop halts a running workspace, preserving its state.
func (m *Manager) Stop(ctx context.Context, id string) error {
	ws, err := m.store.Get(id)
	if err != nil {
		return fmt.Errorf("get workspace: %w", err)
	}
	if err := ws.Transition(StateStopped); err != nil {
		return fmt.Errorf("cannot stop workspace: %w", err)
	}
	if err := m.container.StopContainer(ctx, ws.ContainerID); err != nil {
		return fmt.Errorf("stop container: %w", err)
	}
	if err := m.store.Save(ws); err != nil {
		return fmt.Errorf("save workspace: %w", err)
	}
	m.logger.Info("workspace stopped", "id", id)
	return nil
}

// Start resumes a stopped workspace.
func (m *Manager) Start(ctx context.Context, id string) error {
	ws, err := m.store.Get(id)
	if err != nil {
		return fmt.Errorf("get workspace: %w", err)
	}
	if err := ws.Transition(StateRunning); err != nil {
		return fmt.Errorf("cannot start workspace: %w", err)
	}
	if err := m.container.StartContainer(ctx, ws.ContainerID); err != nil {
		return fmt.Errorf("start container: %w", err)
	}
	if err := m.store.Save(ws); err != nil {
		return fmt.Errorf("save workspace: %w", err)
	}
	m.logger.Info("workspace started", "id", id)
	return nil
}

// Remove destroys a workspace and frees its resources.
func (m *Manager) Remove(ctx context.Context, id string) error {
	ws, err := m.store.Get(id)
	if err != nil {
		return fmt.Errorf("get workspace: %w", err)
	}

	// Stop container first if running
	if ws.State == StateRunning && ws.ContainerID != "" {
		if err := m.container.StopContainer(ctx, ws.ContainerID); err != nil {
			m.logger.Warn("could not stop container before removal", "error", err)
		}
	}

	// Remove container
	if ws.ContainerID != "" {
		if err := m.container.RemoveContainer(ctx, ws.ContainerID); err != nil {
			m.logger.Warn("could not remove container", "error", err)
		}
	}

	// Transition to removed and save
	if err := ws.Transition(StateRemoved); err != nil {
		return fmt.Errorf("cannot remove workspace: %w", err)
	}
	if err := m.store.Save(ws); err != nil {
		return fmt.Errorf("save workspace: %w", err)
	}
	m.logger.Info("workspace removed", "id", id)
	return nil
}

// List returns all workspaces matching the given filter.
func (m *Manager) List(ctx context.Context, includeArchived bool) ([]Workspace, error) {
	return m.store.List(includeArchived)
}

// Get returns a single workspace by ID.
func (m *Manager) Get(ctx context.Context, id string) (*Workspace, error) {
	return m.store.Get(id)
}
