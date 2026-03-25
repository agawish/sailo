package commands

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/agawish/sailo/pkg/config"
	"github.com/agawish/sailo/pkg/container"
	"github.com/agawish/sailo/pkg/detect"
	gitpkg "github.com/agawish/sailo/pkg/git"
	"github.com/agawish/sailo/pkg/port"
	"github.com/agawish/sailo/pkg/workspace"
)

// deps holds initialized dependencies shared across commands.
var deps struct {
	manager   *workspace.Manager
	container *container.Client
	store     *workspace.Store
	logger    *slog.Logger
}

func initDeps() error {
	// 1. Logger
	level := slog.LevelWarn
	if verbose {
		level = slog.LevelDebug
	}
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level}))
	deps.logger = logger

	// 2. SQLite store
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("get home directory: %w", err)
	}
	dbPath := filepath.Join(home, ".sailo", "workspaces.db")
	store, err := workspace.NewStore(dbPath, logger)
	if err != nil {
		return fmt.Errorf("open workspace store: %w", err)
	}
	deps.store = store

	// 3. User config
	userCfg, err := config.LoadUserConfig()
	if err != nil {
		return fmt.Errorf("load user config: %w", err)
	}

	// 4. Parse port range
	minPort, maxPort, err := port.ParsePortRange(userCfg.Defaults.PortRange)
	if err != nil {
		return fmt.Errorf("parse port range: %w", err)
	}

	// 5. Container client
	containerClient, err := container.NewClient(logger)
	if err != nil {
		return fmt.Errorf("create container client: %w", err)
	}
	deps.container = containerClient

	// 6. Port allocator
	portAllocator := port.NewAllocator(minPort, maxPort, store.UsedHostPorts, logger)

	// 7. Git manager
	gitManager := gitpkg.NewManager(containerClient, logger)

	// 8. Detector
	detector := detect.NewDetector(logger)

	// 9. Workspace manager
	deps.manager = workspace.NewManager(workspace.ManagerConfig{
		Store:      store,
		Container:  containerClient,
		Ports:      portAllocator,
		Git:        gitManager,
		Detector:   detector,
		UserConfig: userCfg,
		Logger:     logger,
	})

	return nil
}

var verbose bool
