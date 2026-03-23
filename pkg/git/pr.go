package git

import (
	"context"
	"fmt"
)

// PROptions configures pull request creation.
type PROptions struct {
	Title string
	Body  string
	Base  string
}

// CreatePR creates a pull request using the gh CLI inside the workspace container.
func CreatePR(ctx context.Context, workspaceDir string, opts PROptions) (string, error) {
	return "", fmt.Errorf("PR creation not yet implemented")
}
