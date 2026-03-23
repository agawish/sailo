package workspace

import (
	"fmt"
	"log/slog"
)

// Store persists workspace metadata in SQLite.
type Store struct {
	dbPath string
	logger *slog.Logger
}

// NewStore creates a workspace store backed by SQLite at the given path.
func NewStore(dbPath string, logger *slog.Logger) (*Store, error) {
	return nil, fmt.Errorf("workspace store not yet implemented (db: %s)", dbPath)
}

// Save persists a workspace to the database.
func (s *Store) Save(ws *Workspace) error {
	return fmt.Errorf("workspace save not yet implemented")
}

// Get retrieves a workspace by ID.
func (s *Store) Get(id string) (*Workspace, error) {
	return nil, fmt.Errorf("workspace get not yet implemented")
}

// List returns all workspaces, optionally including archived ones.
func (s *Store) List(includeArchived bool) ([]Workspace, error) {
	return nil, fmt.Errorf("workspace list not yet implemented")
}

// Delete removes a workspace from the database.
func (s *Store) Delete(id string) error {
	return fmt.Errorf("workspace delete not yet implemented")
}
