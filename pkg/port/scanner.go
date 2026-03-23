package port

import "fmt"

// IsPortAvailable checks whether a host port is free.
func IsPortAvailable(port int) (bool, error) {
	return false, fmt.Errorf("port scanning not yet implemented")
}

// FindUsedPorts returns a set of ports currently in use on the host.
func FindUsedPorts() (map[int]bool, error) {
	return nil, fmt.Errorf("port scanning not yet implemented")
}
