package workspace

import (
	"crypto/rand"
	"fmt"
	"regexp"
	"strings"
)

var nonAlphaNum = regexp.MustCompile(`[^a-z0-9]+`)

// generateID creates a unique workspace ID in the format ws-<8hex>.
func generateID() string {
	b := make([]byte, 4)
	rand.Read(b)
	return fmt.Sprintf("ws-%x", b)
}

// slugify converts a task description into a branch-safe slug.
// Lowercases, replaces non-alphanumeric chars with hyphens, truncates to 30 chars.
func slugify(task string) string {
	s := strings.ToLower(task)
	s = nonAlphaNum.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	if len(s) > 30 {
		s = s[:30]
		s = strings.TrimRight(s, "-")
	}
	if s == "" {
		s = "workspace"
	}
	return s
}
