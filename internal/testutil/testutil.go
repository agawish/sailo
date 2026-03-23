// Package testutil provides shared test helpers for sailo tests.
package testutil

import (
	"log/slog"
	"os"
)

// NewTestLogger creates a logger suitable for tests.
func NewTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
}
