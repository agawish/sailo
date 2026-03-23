package workspace

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCanTransition_ValidTransitions(t *testing.T) {
	tests := []struct {
		from State
		to   State
	}{
		{StateCreating, StateRunning},
		{StateCreating, StateFailed},
		{StateRunning, StateStopped},
		{StateRunning, StateShipping},
		{StateRunning, StateFailed},
		{StateRunning, StateRemoved},
		{StateStopped, StateRunning},
		{StateStopped, StateRemoved},
		{StateShipping, StateArchived},
		{StateShipping, StateFailed},
		{StateArchived, StateRemoved},
		{StateFailed, StateRemoved},
	}
	for _, tt := range tests {
		t.Run(string(tt.from)+"→"+string(tt.to), func(t *testing.T) {
			assert.True(t, CanTransition(tt.from, tt.to))
		})
	}
}

func TestCanTransition_InvalidTransitions(t *testing.T) {
	tests := []struct {
		from State
		to   State
	}{
		{StateCreating, StateShipping},
		{StateCreating, StateArchived},
		{StateStopped, StateShipping},
		{StateArchived, StateRunning},
		{StateFailed, StateRunning},
		{StateRemoved, StateRunning},
	}
	for _, tt := range tests {
		t.Run(string(tt.from)+"→"+string(tt.to), func(t *testing.T) {
			assert.False(t, CanTransition(tt.from, tt.to))
		})
	}
}

func TestWorkspace_Transition(t *testing.T) {
	ws := &Workspace{
		ID:    "ws-test",
		State: StateCreating,
	}

	require.NoError(t, ws.Transition(StateRunning))
	assert.Equal(t, StateRunning, ws.State)

	require.NoError(t, ws.Transition(StateShipping))
	assert.Equal(t, StateShipping, ws.State)

	require.NoError(t, ws.Transition(StateArchived))
	assert.Equal(t, StateArchived, ws.State)
}

func TestWorkspace_Transition_Invalid(t *testing.T) {
	ws := &Workspace{
		ID:    "ws-test",
		State: StateCreating,
	}

	err := ws.Transition(StateShipping)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid state transition")
	assert.Equal(t, StateCreating, ws.State) // state unchanged
}
