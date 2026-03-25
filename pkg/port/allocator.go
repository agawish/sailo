package port

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"
)

// Allocator assigns unique host ports to workspace containers.
type Allocator struct {
	minPort   int
	maxPort   int
	usedPorts func() ([]int, error)
	logger    *slog.Logger
}

// NewAllocator creates a port allocator for the given range.
// usedPorts is a function that returns all currently allocated host ports.
func NewAllocator(minPort, maxPort int, usedPorts func() ([]int, error), logger *slog.Logger) *Allocator {
	return &Allocator{
		minPort:   minPort,
		maxPort:   maxPort,
		usedPorts: usedPorts,
		logger:    logger,
	}
}

// Allocate finds a free host port for the given container port.
// It checks both the workspace store and the host OS for conflicts.
func (a *Allocator) Allocate(containerPort int) (int, error) {
	used, err := a.usedPorts()
	if err != nil {
		return 0, fmt.Errorf("query used ports: %w", err)
	}

	usedSet := make(map[int]bool, len(used))
	for _, p := range used {
		usedSet[p] = true
	}

	for port := a.minPort; port <= a.maxPort; port++ {
		if usedSet[port] {
			continue
		}
		available, err := IsPortAvailable(port)
		if err != nil {
			a.logger.Debug("error checking port availability", "port", port, "error", err)
			continue
		}
		if !available {
			continue
		}
		a.logger.Debug("allocated port", "container", containerPort, "host", port)
		return port, nil
	}

	return 0, fmt.Errorf("port range exhausted (%d-%d): %d ports allocated; run 'sailo ps' to check for stale workspaces",
		a.minPort, a.maxPort, len(used))
}

// Release is a no-op. Ports are freed automatically when workspaces are removed.
func (a *Allocator) Release(hostPort int) error {
	a.logger.Debug("port released", "port", hostPort)
	return nil
}

// AllocatedPorts returns all currently allocated host ports.
func (a *Allocator) AllocatedPorts() (map[int]int, error) {
	used, err := a.usedPorts()
	if err != nil {
		return nil, fmt.Errorf("query used ports: %w", err)
	}
	result := make(map[int]int, len(used))
	for _, p := range used {
		result[p] = p
	}
	return result, nil
}

// ParsePortRange parses a port range string like "3001-3999".
func ParsePortRange(s string) (min, max int, err error) {
	parts := strings.SplitN(s, "-", 2)
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid port range %q: expected format min-max", s)
	}
	min, err = strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return 0, 0, fmt.Errorf("invalid port range min %q: %w", parts[0], err)
	}
	max, err = strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return 0, 0, fmt.Errorf("invalid port range max %q: %w", parts[1], err)
	}
	if min > max {
		return 0, 0, fmt.Errorf("invalid port range: min %d > max %d", min, max)
	}
	if min < 1 || max > 65535 {
		return 0, 0, fmt.Errorf("invalid port range: must be between 1 and 65535")
	}
	return min, max, nil
}
