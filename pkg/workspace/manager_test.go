package workspace

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/agawish/sailo/internal/testutil"
	"github.com/agawish/sailo/pkg/config"
	"github.com/agawish/sailo/pkg/detect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockDetector implements workspace.ProjectDetector for testing.
type mockDetector struct {
	Result *detect.Result
	Err    error
}

func (m *mockDetector) Detect(projectDir string) (*detect.Result, error) {
	if m.Result == nil {
		return &detect.Result{BaseImage: "ubuntu:24.04"}, m.Err
	}
	return m.Result, m.Err
}

func newTestManager(t *testing.T, mc *testutil.MockContainerClient, mp *testutil.MockPortAllocator, mg *testutil.MockGitOperator, md *mockDetector) (*Manager, *Store) {
	t.Helper()
	dbPath := filepath.Join(t.TempDir(), "test.db")
	store, err := NewStore(dbPath, testutil.NewTestLogger())
	require.NoError(t, err)
	t.Cleanup(func() { store.Close() })

	mgr := NewManager(ManagerConfig{
		Store:      store,
		Container:  mc,
		Ports:      mp,
		Git:        mg,
		Detector:   md,
		UserConfig: config.DefaultUserConfig(),
		Logger:     testutil.NewTestLogger(),
	})
	return mgr, store
}

func TestManager_Create_HappyPath(t *testing.T) {
	mc := &testutil.MockContainerClient{CreateID: "container-abc123456789"}
	mp := &testutil.MockPortAllocator{AllocatePort: 3007}
	mg := &testutil.MockGitOperator{}
	md := &mockDetector{}

	// Override GetRemoteURL for testing
	origGetRemoteURL := getRemoteURL
	getRemoteURL = func() (string, error) { return "git@github.com:user/repo.git", nil }
	defer func() { getRemoteURL = origGetRemoteURL }()

	mgr, store := newTestManager(t, mc, mp, mg, md)

	ws, err := mgr.Create(context.Background(), CreateOptions{
		Task:       "add dark mode",
		FromBranch: "main",
	})
	require.NoError(t, err)
	require.NotNil(t, ws)

	assert.Regexp(t, `^ws-[0-9a-f]{8}$`, ws.ID)
	assert.Equal(t, "add dark mode", ws.Task)
	assert.Equal(t, StateRunning, ws.State)
	assert.Contains(t, ws.Branch, "sailo/")
	assert.Contains(t, ws.Branch, "add-dark-mode")
	assert.Equal(t, "main", ws.FromBranch)
	assert.Equal(t, 1, mc.CreateCalls)
	assert.Equal(t, 1, mg.CloneCalls)

	// Verify persisted
	got, err := store.Get(ws.ID)
	require.NoError(t, err)
	assert.Equal(t, StateRunning, got.State)
}

func TestManager_Create_DockerDown(t *testing.T) {
	mc := &testutil.MockContainerClient{PingErr: fmt.Errorf("docker not reachable")}
	mp := &testutil.MockPortAllocator{}
	mg := &testutil.MockGitOperator{}
	md := &mockDetector{}

	mgr, _ := newTestManager(t, mc, mp, mg, md)

	_, err := mgr.Create(context.Background(), CreateOptions{Task: "test", FromBranch: "main"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "docker not reachable")
}

func TestManager_Create_PortExhaustion(t *testing.T) {
	mc := &testutil.MockContainerClient{CreateID: "container-abc"}
	mp := &testutil.MockPortAllocator{AllocateErr: fmt.Errorf("port range exhausted")}
	mg := &testutil.MockGitOperator{}
	md := &mockDetector{Result: &detect.Result{BaseImage: "node:22", Ports: []int{3000}}}

	origGetRemoteURL := getRemoteURL
	getRemoteURL = func() (string, error) { return "git@github.com:user/repo.git", nil }
	defer func() { getRemoteURL = origGetRemoteURL }()

	mgr, _ := newTestManager(t, mc, mp, mg, md)

	_, err := mgr.Create(context.Background(), CreateOptions{Task: "test", FromBranch: "main"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "port range exhausted")
}

func TestManager_Create_CloneFail_Cleanup(t *testing.T) {
	mc := &testutil.MockContainerClient{CreateID: "container-abc123456789"}
	mp := &testutil.MockPortAllocator{AllocatePort: 3007}
	mg := &testutil.MockGitOperator{CloneErr: fmt.Errorf("SSH auth failed")}
	md := &mockDetector{}

	origGetRemoteURL := getRemoteURL
	getRemoteURL = func() (string, error) { return "git@github.com:user/repo.git", nil }
	defer func() { getRemoteURL = origGetRemoteURL }()

	mgr, store := newTestManager(t, mc, mp, mg, md)

	_, err := mgr.Create(context.Background(), CreateOptions{Task: "test", FromBranch: "main"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "clone repository")

	// Verify cleanup happened
	assert.Equal(t, 1, mc.RemoveCalls, "container should be removed on failure")

	// Verify workspace record was deleted
	list, err := store.List(true)
	require.NoError(t, err)
	assert.Empty(t, list, "workspace record should be deleted on failure")
}

func TestManager_Create_ContainerFail_Cleanup(t *testing.T) {
	mc := &testutil.MockContainerClient{CreateErr: fmt.Errorf("image not found")}
	mp := &testutil.MockPortAllocator{AllocatePort: 3007}
	mg := &testutil.MockGitOperator{}
	md := &mockDetector{}

	origGetRemoteURL := getRemoteURL
	getRemoteURL = func() (string, error) { return "git@github.com:user/repo.git", nil }
	defer func() { getRemoteURL = origGetRemoteURL }()

	mgr, store := newTestManager(t, mc, mp, mg, md)

	_, err := mgr.Create(context.Background(), CreateOptions{Task: "test", FromBranch: "main"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "create container")

	// Container was never created so RemoveCalls should be 0 (cleanup skips if containerID is empty)
	assert.Equal(t, 0, mc.RemoveCalls)

	// But workspace record should still be cleaned up
	list, err := store.List(true)
	require.NoError(t, err)
	assert.Empty(t, list)
}

func TestManager_Stop_HappyPath(t *testing.T) {
	mc := &testutil.MockContainerClient{}
	mgr, store := newTestManager(t, mc, &testutil.MockPortAllocator{}, &testutil.MockGitOperator{}, &mockDetector{})

	ws := &Workspace{ID: "ws-test", Task: "test", State: StateRunning, ContainerID: "container-123"}
	require.NoError(t, store.Save(ws))

	err := mgr.Stop(context.Background(), "ws-test")
	require.NoError(t, err)

	got, _ := store.Get("ws-test")
	assert.Equal(t, StateStopped, got.State)
	assert.Equal(t, 1, mc.StopCalls)
}

func TestManager_Stop_WrongState(t *testing.T) {
	mc := &testutil.MockContainerClient{}
	mgr, store := newTestManager(t, mc, &testutil.MockPortAllocator{}, &testutil.MockGitOperator{}, &mockDetector{})

	ws := &Workspace{ID: "ws-test", Task: "test", State: StateCreating, ContainerID: "container-123"}
	require.NoError(t, store.Save(ws))

	err := mgr.Stop(context.Background(), "ws-test")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot stop")
}

func TestManager_Start_HappyPath(t *testing.T) {
	mc := &testutil.MockContainerClient{}
	mgr, store := newTestManager(t, mc, &testutil.MockPortAllocator{}, &testutil.MockGitOperator{}, &mockDetector{})

	ws := &Workspace{ID: "ws-test", Task: "test", State: StateStopped, ContainerID: "container-123"}
	require.NoError(t, store.Save(ws))

	err := mgr.Start(context.Background(), "ws-test")
	require.NoError(t, err)

	got, _ := store.Get("ws-test")
	assert.Equal(t, StateRunning, got.State)
}

func TestManager_Remove_RunningWorkspace(t *testing.T) {
	mc := &testutil.MockContainerClient{}
	mgr, store := newTestManager(t, mc, &testutil.MockPortAllocator{}, &testutil.MockGitOperator{}, &mockDetector{})

	ws := &Workspace{ID: "ws-test", Task: "test", State: StateRunning, ContainerID: "container-123"}
	require.NoError(t, store.Save(ws))

	err := mgr.Remove(context.Background(), "ws-test")
	require.NoError(t, err)

	assert.Equal(t, 1, mc.StopCalls, "should stop running container first")
	assert.Equal(t, 1, mc.RemoveCalls, "should remove container")

	got, _ := store.Get("ws-test")
	assert.Equal(t, StateRemoved, got.State)

	// Should not appear in list
	list, _ := store.List(false)
	assert.Empty(t, list)
}

func TestManager_Remove_StoppedWorkspace(t *testing.T) {
	mc := &testutil.MockContainerClient{}
	mgr, store := newTestManager(t, mc, &testutil.MockPortAllocator{}, &testutil.MockGitOperator{}, &mockDetector{})

	ws := &Workspace{ID: "ws-test", Task: "test", State: StateStopped, ContainerID: "container-123"}
	require.NoError(t, store.Save(ws))

	err := mgr.Remove(context.Background(), "ws-test")
	require.NoError(t, err)

	assert.Equal(t, 0, mc.StopCalls, "should not stop already-stopped container")
	assert.Equal(t, 1, mc.RemoveCalls, "should still remove container")
}

func TestManager_Remove_NotFound(t *testing.T) {
	mgr, _ := newTestManager(t, &testutil.MockContainerClient{}, &testutil.MockPortAllocator{}, &testutil.MockGitOperator{}, &mockDetector{})

	err := mgr.Remove(context.Background(), "ws-nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestManager_List_Empty(t *testing.T) {
	mgr, _ := newTestManager(t, &testutil.MockContainerClient{}, &testutil.MockPortAllocator{}, &testutil.MockGitOperator{}, &mockDetector{})

	list, err := mgr.List(context.Background(), false)
	require.NoError(t, err)
	assert.NotNil(t, list)
	assert.Empty(t, list)
}

func TestManager_Get(t *testing.T) {
	mgr, store := newTestManager(t, &testutil.MockContainerClient{}, &testutil.MockPortAllocator{}, &testutil.MockGitOperator{}, &mockDetector{})

	ws := &Workspace{ID: "ws-test", Task: "test", State: StateRunning}
	require.NoError(t, store.Save(ws))

	got, err := mgr.Get(context.Background(), "ws-test")
	require.NoError(t, err)
	assert.Equal(t, "ws-test", got.ID)
}
