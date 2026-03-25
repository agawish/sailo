package workspace

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/agawish/sailo/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestStore(t *testing.T) *Store {
	t.Helper()
	dbPath := filepath.Join(t.TempDir(), "sailo", "workspaces.db")
	store, err := NewStore(dbPath, testutil.NewTestLogger())
	require.NoError(t, err)
	t.Cleanup(func() { store.Close() })
	return store
}

func TestNewStore_CreatesDatabase(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "sailo", "workspaces.db")
	store, err := NewStore(dbPath, testutil.NewTestLogger())
	require.NoError(t, err)
	defer store.Close()

	assert.FileExists(t, dbPath)
}

func TestStore_SaveAndGet(t *testing.T) {
	store := newTestStore(t)

	ws := &Workspace{
		ID:          "ws-abc123",
		Task:        "add dark mode",
		State:       StateRunning,
		Branch:      "sailo/ws-abc123/add-dark-mode",
		ContainerID: "container-xyz",
		Ports:       map[int]int{3000: 3007, 8080: 3008},
		FromBranch:  "main",
	}

	require.NoError(t, store.Save(ws))

	got, err := store.Get("ws-abc123")
	require.NoError(t, err)

	assert.Equal(t, ws.ID, got.ID)
	assert.Equal(t, ws.Task, got.Task)
	assert.Equal(t, StateRunning, got.State)
	assert.Equal(t, ws.Branch, got.Branch)
	assert.Equal(t, ws.ContainerID, got.ContainerID)
	assert.Equal(t, ws.Ports, got.Ports)
	assert.Equal(t, ws.FromBranch, got.FromBranch)
	assert.NotEmpty(t, got.CreatedAt)
	assert.NotEmpty(t, got.UpdatedAt)
}

func TestStore_SaveAndGet_EmptyPorts(t *testing.T) {
	store := newTestStore(t)

	ws := &Workspace{
		ID:    "ws-empty",
		Task:  "test empty ports",
		State: StateCreating,
		Ports: map[int]int{},
	}

	require.NoError(t, store.Save(ws))

	got, err := store.Get("ws-empty")
	require.NoError(t, err)

	assert.NotNil(t, got.Ports)
	assert.Empty(t, got.Ports)
}

func TestStore_SaveAndGet_NilPorts(t *testing.T) {
	store := newTestStore(t)

	ws := &Workspace{
		ID:    "ws-nil",
		Task:  "test nil ports",
		State: StateCreating,
	}

	require.NoError(t, store.Save(ws))

	got, err := store.Get("ws-nil")
	require.NoError(t, err)

	assert.NotNil(t, got.Ports)
	assert.Empty(t, got.Ports)
}

func TestStore_Get_NotFound(t *testing.T) {
	store := newTestStore(t)

	_, err := store.Get("ws-nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestStore_List(t *testing.T) {
	store := newTestStore(t)

	workspaces := []Workspace{
		{ID: "ws-1", Task: "task 1", State: StateRunning},
		{ID: "ws-2", Task: "task 2", State: StateStopped},
		{ID: "ws-3", Task: "task 3", State: StateArchived},
		{ID: "ws-4", Task: "task 4", State: StateRemoved},
	}
	for i := range workspaces {
		require.NoError(t, store.Save(&workspaces[i]))
	}

	t.Run("exclude archived", func(t *testing.T) {
		list, err := store.List(false)
		require.NoError(t, err)
		assert.Len(t, list, 2) // running + stopped only
		ids := []string{list[0].ID, list[1].ID}
		assert.Contains(t, ids, "ws-1")
		assert.Contains(t, ids, "ws-2")
	})

	t.Run("include archived", func(t *testing.T) {
		list, err := store.List(true)
		require.NoError(t, err)
		assert.Len(t, list, 3) // running + stopped + archived, never removed
		ids := []string{list[0].ID, list[1].ID, list[2].ID}
		assert.Contains(t, ids, "ws-1")
		assert.Contains(t, ids, "ws-2")
		assert.Contains(t, ids, "ws-3")
	})
}

func TestStore_List_Empty(t *testing.T) {
	store := newTestStore(t)

	list, err := store.List(false)
	require.NoError(t, err)
	assert.NotNil(t, list)
	assert.Empty(t, list)
}

func TestStore_Delete(t *testing.T) {
	store := newTestStore(t)

	ws := &Workspace{ID: "ws-del", Task: "to delete", State: StateRunning}
	require.NoError(t, store.Save(ws))

	require.NoError(t, store.Delete("ws-del"))

	_, err := store.Get("ws-del")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestStore_Delete_NotFound(t *testing.T) {
	store := newTestStore(t)

	err := store.Delete("ws-ghost")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestStore_Save_UpdatesTimestamp(t *testing.T) {
	store := newTestStore(t)

	ws := &Workspace{ID: "ws-ts", Task: "timestamp test", State: StateCreating}
	require.NoError(t, store.Save(ws))

	first, err := store.Get("ws-ts")
	require.NoError(t, err)

	// RFC3339 has second-level precision, so wait past the boundary
	time.Sleep(1100 * time.Millisecond)

	ws.State = StateRunning
	require.NoError(t, store.Save(ws))

	second, err := store.Get("ws-ts")
	require.NoError(t, err)

	assert.Equal(t, first.CreatedAt, second.CreatedAt)
	assert.NotEqual(t, first.UpdatedAt, second.UpdatedAt)
}

func TestStore_MigrateIdempotent(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "workspaces.db")
	logger := testutil.NewTestLogger()

	store1, err := NewStore(dbPath, logger)
	require.NoError(t, err)

	ws := &Workspace{ID: "ws-migrate", Task: "persist", State: StateRunning}
	require.NoError(t, store1.Save(ws))
	store1.Close()

	store2, err := NewStore(dbPath, logger)
	require.NoError(t, err)
	defer store2.Close()

	got, err := store2.Get("ws-migrate")
	require.NoError(t, err)
	assert.Equal(t, "persist", got.Task)
}
