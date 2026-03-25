package testutil

import (
	"context"

	"github.com/agawish/sailo/pkg/container"
)

// MockContainerClient implements workspace.ContainerClient for testing.
type MockContainerClient struct {
	PingErr    error
	CreateID   string
	CreateErr  error
	StopErr    error
	StartErr   error
	RemoveErr  error
	ExecErr    error
	ExecOutput string
	InspectState container.ContainerState
	InspectErr   error

	CreateCalls int
	RemoveCalls int
	ExecCalls   int
	StopCalls   int
}

func (m *MockContainerClient) Ping(ctx context.Context) error {
	return m.PingErr
}

func (m *MockContainerClient) CreateWorkspace(ctx context.Context, opts container.CreateOptions) (string, error) {
	m.CreateCalls++
	return m.CreateID, m.CreateErr
}

func (m *MockContainerClient) StopContainer(ctx context.Context, containerID string) error {
	m.StopCalls++
	return m.StopErr
}

func (m *MockContainerClient) StartContainer(ctx context.Context, containerID string) error {
	return m.StartErr
}

func (m *MockContainerClient) RemoveContainer(ctx context.Context, containerID string) error {
	m.RemoveCalls++
	return m.RemoveErr
}

func (m *MockContainerClient) ExecInContainer(ctx context.Context, containerID string, cmd []string) error {
	m.ExecCalls++
	return m.ExecErr
}

func (m *MockContainerClient) ExecWithOutput(ctx context.Context, containerID string, cmd []string) (string, error) {
	m.ExecCalls++
	return m.ExecOutput, m.ExecErr
}

func (m *MockContainerClient) InspectContainer(ctx context.Context, containerID string) (container.ContainerState, error) {
	return m.InspectState, m.InspectErr
}

// MockPortAllocator implements workspace.PortAllocator for testing.
type MockPortAllocator struct {
	AllocatePort int
	AllocateErr  error
	AllocateCalls int
}

func (m *MockPortAllocator) Allocate(containerPort int) (int, error) {
	m.AllocateCalls++
	return m.AllocatePort, m.AllocateErr
}

// MockGitOperator implements workspace.GitOperator for testing.
type MockGitOperator struct {
	CloneErr        error
	CreateBranchErr error
	CloneCalls      int
}

func (m *MockGitOperator) Clone(ctx context.Context, containerID string, repoURL string, branch string, targetDir string) error {
	m.CloneCalls++
	return m.CloneErr
}

func (m *MockGitOperator) CreateBranch(ctx context.Context, containerID string, workspaceDir string, branchName string) error {
	return m.CreateBranchErr
}

// Note: MockDetector is defined in pkg/workspace/manager_test.go
// to avoid import cycles with pkg/detect.
