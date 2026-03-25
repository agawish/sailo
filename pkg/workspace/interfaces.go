package workspace

import (
	"context"

	"github.com/agawish/sailo/pkg/container"
	"github.com/agawish/sailo/pkg/detect"
)

// ContainerClient is the interface for Docker container operations.
type ContainerClient interface {
	Ping(ctx context.Context) error
	CreateWorkspace(ctx context.Context, opts container.CreateOptions) (string, error)
	StopContainer(ctx context.Context, containerID string) error
	StartContainer(ctx context.Context, containerID string) error
	RemoveContainer(ctx context.Context, containerID string) error
	ExecInContainer(ctx context.Context, containerID string, cmd []string) error
	ExecWithOutput(ctx context.Context, containerID string, cmd []string) (string, error)
	InspectContainer(ctx context.Context, containerID string) (container.ContainerState, error)
}

// PortAllocator is the interface for port allocation.
type PortAllocator interface {
	Allocate(containerPort int) (int, error)
}

// GitOperator is the interface for git operations inside containers.
type GitOperator interface {
	Clone(ctx context.Context, containerID string, repoURL string, branch string, targetDir string) error
	CreateBranch(ctx context.Context, containerID string, workspaceDir string, branchName string) error
}

// ProjectDetector is the interface for project detection.
type ProjectDetector interface {
	Detect(projectDir string) (*detect.Result, error)
}
