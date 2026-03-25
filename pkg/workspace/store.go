package workspace

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

// Store persists workspace metadata in SQLite.
type Store struct {
	db     *sql.DB
	dbPath string
	logger *slog.Logger
}

const createTableSQL = `CREATE TABLE IF NOT EXISTS workspaces (
	id           TEXT PRIMARY KEY,
	task         TEXT NOT NULL,
	state        TEXT NOT NULL DEFAULT 'creating',
	branch       TEXT NOT NULL DEFAULT '',
	container_id TEXT NOT NULL DEFAULT '',
	ports        TEXT NOT NULL DEFAULT '{}',
	from_branch  TEXT NOT NULL DEFAULT 'main',
	created_at   TEXT NOT NULL,
	updated_at   TEXT NOT NULL
);`

// NewStore creates a workspace store backed by SQLite at the given path.
func NewStore(dbPath string, logger *slog.Logger) (*Store, error) {
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("create store directory: %w", err)
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		db.Close()
		return nil, fmt.Errorf("set journal mode: %w", err)
	}

	if _, err := db.Exec(createTableSQL); err != nil {
		db.Close()
		return nil, fmt.Errorf("create workspaces table: %w", err)
	}

	logger.Debug("workspace store opened", "path", dbPath)
	return &Store{db: db, dbPath: dbPath, logger: logger}, nil
}

// Save persists a workspace to the database (upsert).
func (s *Store) Save(ws *Workspace) error {
	now := time.Now().UTC().Format(time.RFC3339)
	ws.UpdatedAt = now
	if ws.CreatedAt == "" {
		ws.CreatedAt = now
	}

	portsJSON, err := json.Marshal(ws.Ports)
	if err != nil {
		return fmt.Errorf("marshal ports: %w", err)
	}

	_, err = s.db.Exec(`INSERT OR REPLACE INTO workspaces
		(id, task, state, branch, container_id, ports, from_branch, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		ws.ID, ws.Task, string(ws.State), ws.Branch, ws.ContainerID,
		string(portsJSON), ws.FromBranch, ws.CreatedAt, ws.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("save workspace %s: %w", ws.ID, err)
	}

	s.logger.Debug("workspace saved", "id", ws.ID, "state", ws.State)
	return nil
}

// Get retrieves a workspace by ID.
func (s *Store) Get(id string) (*Workspace, error) {
	row := s.db.QueryRow(`SELECT id, task, state, branch, container_id, ports, from_branch, created_at, updated_at
		FROM workspaces WHERE id = ?`, id)

	ws, err := scanWorkspace(row)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("workspace not found: %s", id)
	}
	if err != nil {
		return nil, fmt.Errorf("get workspace %s: %w", id, err)
	}
	return ws, nil
}

// List returns all workspaces, optionally including archived ones.
// Removed workspaces are always excluded.
func (s *Store) List(includeArchived bool) ([]Workspace, error) {
	query := `SELECT id, task, state, branch, container_id, ports, from_branch, created_at, updated_at
		FROM workspaces WHERE state != ?`
	args := []interface{}{string(StateRemoved)}

	if !includeArchived {
		query += " AND state != ?"
		args = append(args, string(StateArchived))
	}

	query += " ORDER BY created_at DESC"

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("list workspaces: %w", err)
	}
	defer rows.Close()

	var workspaces []Workspace
	for rows.Next() {
		ws, err := scanWorkspace(rows)
		if err != nil {
			return nil, fmt.Errorf("scan workspace: %w", err)
		}
		workspaces = append(workspaces, *ws)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate workspaces: %w", err)
	}

	if workspaces == nil {
		workspaces = []Workspace{}
	}
	return workspaces, nil
}

// Delete removes a workspace from the database.
func (s *Store) Delete(id string) error {
	result, err := s.db.Exec("DELETE FROM workspaces WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("delete workspace %s: %w", id, err)
	}

	n, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check delete result: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("workspace not found: %s", id)
	}

	s.logger.Debug("workspace deleted", "id", id)
	return nil
}

// Close closes the underlying database connection.
func (s *Store) Close() error {
	return s.db.Close()
}

// UsedHostPorts returns all host ports allocated to non-removed workspaces.
func (s *Store) UsedHostPorts() ([]int, error) {
	rows, err := s.db.Query("SELECT ports FROM workspaces WHERE state != ?", string(StateRemoved))
	if err != nil {
		return nil, fmt.Errorf("query used ports: %w", err)
	}
	defer rows.Close()

	var allPorts []int
	for rows.Next() {
		var portsJSON string
		if err := rows.Scan(&portsJSON); err != nil {
			return nil, fmt.Errorf("scan ports: %w", err)
		}
		var portMap map[int]int
		if err := json.Unmarshal([]byte(portsJSON), &portMap); err != nil {
			continue // skip malformed entries
		}
		for _, hp := range portMap {
			allPorts = append(allPorts, hp)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate ports: %w", err)
	}
	return allPorts, nil
}

// scanner is an interface satisfied by both *sql.Row and *sql.Rows.
type scanner interface {
	Scan(dest ...interface{}) error
}

func scanWorkspace(s scanner) (*Workspace, error) {
	var ws Workspace
	var state, portsJSON string

	err := s.Scan(&ws.ID, &ws.Task, &state, &ws.Branch, &ws.ContainerID,
		&portsJSON, &ws.FromBranch, &ws.CreatedAt, &ws.UpdatedAt)
	if err != nil {
		return nil, err
	}

	ws.State = State(state)

	if err := json.Unmarshal([]byte(portsJSON), &ws.Ports); err != nil {
		return nil, fmt.Errorf("unmarshal ports: %w", err)
	}
	if ws.Ports == nil {
		ws.Ports = map[int]int{}
	}

	return &ws, nil
}
