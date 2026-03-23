package container

import (
	"context"
	"fmt"
	"log/slog"
)

// Builder handles Docker image building for workspaces.
type Builder struct {
	logger *slog.Logger
}

// NewBuilder creates an image builder.
func NewBuilder(logger *slog.Logger) *Builder {
	return &Builder{logger: logger}
}

// Build builds a Docker image from a Dockerfile in the given context directory.
func (b *Builder) Build(ctx context.Context, contextDir string, dockerfile string) (string, error) {
	return "", fmt.Errorf("image building not yet implemented")
}
