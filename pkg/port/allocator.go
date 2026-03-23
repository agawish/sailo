// Package port manages non-conflicting port allocation for workspaces.
package port

import (
	"fmt"
	"log/slog"
)

// Allocator assigns unique host ports to workspace containers.
type Allocator struct {
	minPort int
	maxPort int
	logger  *slog.Logger
}

// NewAllocator creates a port allocator for the given range.
func NewAllocator(minPort, maxPort int, logger *slog.Logger) *Allocator {
	return &Allocator{
		minPort: minPort,
		maxPort: maxPort,
		logger:  logger,
	}
}

// Allocate finds and reserves a free host port for the given container port.
// Uses SQLite transaction locking to prevent race conditions.
func (a *Allocator) Allocate(containerPort int) (int, error) {
	return 0, fmt.Errorf("port allocation not yet implemented")
}

// Release frees a previously allocated port.
func (a *Allocator) Release(hostPort int) error {
	return fmt.Errorf("port release not yet implemented")
}

// AllocatedPorts returns all currently allocated port mappings.
func (a *Allocator) AllocatedPorts() (map[int]int, error) {
	return nil, fmt.Errorf("port listing not yet implemented")
}
