package workspace

import "fmt"

// State represents the lifecycle state of a workspace.
type State string

const (
	StateCreating State = "creating"
	StateRunning  State = "running"
	StateStopped  State = "stopped"
	StateShipping State = "shipping"
	StateArchived State = "archived"
	StateFailed   State = "failed"
	StateRemoved  State = "removed"
)

// Workspace holds all metadata for an isolated agent workspace.
type Workspace struct {
	ID          string            `json:"id"`
	Task        string            `json:"task"`
	State       State             `json:"state"`
	Branch      string            `json:"branch"`
	ContainerID string            `json:"container_id"`
	Ports       map[int]int       `json:"ports"` // container port → host port
	FromBranch  string            `json:"from_branch"`
	CreatedAt   string            `json:"created_at"`
	UpdatedAt   string            `json:"updated_at"`
}

// validTransitions defines the allowed state machine transitions.
//
//	                sailo create
//	                    │
//	                    ▼
//	┌──────────┐  ┌──────────┐  ┌──────────┐
//	│ CREATING ├─►│ RUNNING  ├─►│ STOPPED  │──► RUNNING (restart)
//	└────┬─────┘  └────┬─────┘  └────┬─────┘
//	     │ (err)       │ ship       │ rm
//	     ▼             ▼            ▼
//	┌──────────┐  ┌──────────┐  ┌──────────┐
//	│  FAILED  │  │ SHIPPING │  │ REMOVED  │
//	└──────────┘  └────┬─────┘  └──────────┘
//	                   ▼
//	              ┌──────────┐
//	              │ ARCHIVED │
//	              └──────────┘
var validTransitions = map[State][]State{
	StateCreating: {StateRunning, StateFailed},
	StateRunning:  {StateStopped, StateShipping, StateFailed, StateRemoved},
	StateStopped:  {StateRunning, StateRemoved},
	StateShipping: {StateArchived, StateFailed},
	StateArchived: {StateRemoved},
	StateFailed:   {StateRemoved},
}

// CanTransition checks whether a state transition is valid.
func CanTransition(from, to State) bool {
	allowed, ok := validTransitions[from]
	if !ok {
		return false
	}
	for _, s := range allowed {
		if s == to {
			return true
		}
	}
	return false
}

// Transition attempts to move the workspace to a new state.
// Returns an error if the transition is invalid.
func (w *Workspace) Transition(to State) error {
	if !CanTransition(w.State, to) {
		return fmt.Errorf("invalid state transition: %s → %s", w.State, to)
	}
	w.State = to
	return nil
}
