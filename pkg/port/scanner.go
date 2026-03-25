package port

import (
	"errors"
	"fmt"
	"net"
	"syscall"
)

// IsPortAvailable checks whether a TCP port is free on the host.
func IsPortAvailable(port int) (bool, error) {
	ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		if isAddrInUse(err) {
			return false, nil
		}
		return false, fmt.Errorf("check port %d: %w", port, err)
	}
	ln.Close()
	return true, nil
}

// FindUsedPorts scans the given range and returns ports that are in use.
func FindUsedPorts(minPort, maxPort int) (map[int]bool, error) {
	used := make(map[int]bool)
	for p := minPort; p <= maxPort; p++ {
		available, err := IsPortAvailable(p)
		if err != nil {
			return nil, err
		}
		if !available {
			used[p] = true
		}
	}
	return used, nil
}

func isAddrInUse(err error) bool {
	var sysErr syscall.Errno
	if errors.As(err, &sysErr) {
		return sysErr == syscall.EADDRINUSE
	}
	return false
}
