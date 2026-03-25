package workspace

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateID_Format(t *testing.T) {
	id := generateID()
	assert.Regexp(t, `^ws-[0-9a-f]{8}$`, id)
}

func TestGenerateID_Uniqueness(t *testing.T) {
	seen := make(map[string]bool)
	for i := 0; i < 1000; i++ {
		id := generateID()
		assert.False(t, seen[id], "duplicate ID: %s", id)
		seen[id] = true
	}
}

func TestSlugify(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"add dark mode", "add-dark-mode"},
		{"Fix BUG #123!!!", "fix-bug-123"},
		{"simple", "simple"},
		{"UPPERCASE TEXT", "uppercase-text"},
		{"  spaces  everywhere  ", "spaces-everywhere"},
		{"a very long task description that exceeds thirty characters limit", "a-very-long-task-description-t"},
		{"", "workspace"},
		{"---", "workspace"},
		{"hello---world", "hello-world"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.expected, slugify(tt.input))
		})
	}
}
