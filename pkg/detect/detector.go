// Package detect scans a project directory to determine the best Docker
// configuration for a workspace. It reuses existing Dockerfiles and
// docker-compose.yml files when available, falling back to language
// detection for base image selection.
package detect

import "log/slog"

// Result holds the detected project configuration.
type Result struct {
	// Dockerfile path if one exists (empty if none found)
	Dockerfile string

	// DockerCompose path if one exists
	DockerCompose string

	// DevContainer path if one exists
	DevContainer string

	// Detected programming language/framework
	Language string

	// Suggested base Docker image
	BaseImage string

	// Detected ports from project configuration
	Ports []int
}

// Detector scans a project to determine its Docker configuration.
type Detector struct {
	logger *slog.Logger
}

// NewDetector creates a project detector.
func NewDetector(logger *slog.Logger) *Detector {
	return &Detector{logger: logger}
}

// Detect scans the given directory and returns the detected configuration.
//
// Detection priority:
//  1. .sailo.yaml (explicit config)
//  2. Dockerfile / docker-compose.yml (reuse as-is)
//  3. devcontainer.json (honor devcontainer spec)
//  4. Language detection (package.json → node, go.mod → go, etc.)
func (d *Detector) Detect(projectDir string) (*Result, error) {
	result := &Result{}

	d.detectDockerfile(projectDir, result)
	d.detectCompose(projectDir, result)
	d.detectDevContainer(projectDir, result)
	d.detectLanguage(projectDir, result)
	d.detectPorts(projectDir, result)

	return result, nil
}
