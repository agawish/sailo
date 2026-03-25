package port

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsPortAvailable_FreePort(t *testing.T) {
	// Find a free port by binding to :0
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	port := ln.Addr().(*net.TCPAddr).Port
	ln.Close()

	available, err := IsPortAvailable(port)
	require.NoError(t, err)
	assert.True(t, available)
}

func TestIsPortAvailable_UsedPort(t *testing.T) {
	// Bind a port
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer ln.Close()
	port := ln.Addr().(*net.TCPAddr).Port

	available, err := IsPortAvailable(port)
	require.NoError(t, err)
	assert.False(t, available)
}
