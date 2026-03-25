package port

import (
	"net"
	"testing"

	"github.com/agawish/sailo/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAllocate_FirstPort(t *testing.T) {
	a := NewAllocator(9100, 9200, func() ([]int, error) { return nil, nil }, testutil.NewTestLogger())

	port, err := a.Allocate(3000)
	require.NoError(t, err)
	assert.Equal(t, 9100, port)
}

func TestAllocate_SkipsUsedPorts(t *testing.T) {
	used := []int{9100, 9101}
	a := NewAllocator(9100, 9200, func() ([]int, error) { return used, nil }, testutil.NewTestLogger())

	port, err := a.Allocate(3000)
	require.NoError(t, err)
	assert.Equal(t, 9102, port)
}

func TestAllocate_SkipsOSUsedPort(t *testing.T) {
	// Bind a port so the OS reports it as in use
	ln, err := net.Listen("tcp", "127.0.0.1:9100")
	if err != nil {
		t.Skip("cannot bind port 9100, skipping")
	}
	defer ln.Close()

	a := NewAllocator(9100, 9200, func() ([]int, error) { return nil, nil }, testutil.NewTestLogger())

	port, err := a.Allocate(3000)
	require.NoError(t, err)
	assert.NotEqual(t, 9100, port)
	assert.True(t, port >= 9101 && port <= 9200)
}

func TestAllocate_Exhausted(t *testing.T) {
	// Range of 1 port, already used
	a := NewAllocator(9100, 9100, func() ([]int, error) { return []int{9100}, nil }, testutil.NewTestLogger())

	_, err := a.Allocate(3000)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "port range exhausted")
}

func TestRelease_NoOp(t *testing.T) {
	a := NewAllocator(9100, 9200, func() ([]int, error) { return nil, nil }, testutil.NewTestLogger())
	assert.NoError(t, a.Release(9100))
}

func TestParsePortRange(t *testing.T) {
	tests := []struct {
		input   string
		min     int
		max     int
		wantErr bool
	}{
		{"3001-3999", 3001, 3999, false},
		{"1024-65535", 1024, 65535, false},
		{"8080-8080", 8080, 8080, false},
		{"invalid", 0, 0, true},
		{"100-50", 0, 0, true},   // min > max
		{"0-100", 0, 0, true},    // min < 1
		{"abc-def", 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			min, max, err := ParsePortRange(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.min, min)
				assert.Equal(t, tt.max, max)
			}
		})
	}
}
