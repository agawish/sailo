package detect

import (
	"os"
	"path/filepath"
)

func (d *Detector) detectDockerfile(projectDir string, result *Result) {
	path := filepath.Join(projectDir, "Dockerfile")
	if _, err := os.Stat(path); err == nil {
		result.Dockerfile = path
		d.logger.Debug("found Dockerfile", "path", path)
	}
}

func (d *Detector) detectCompose(projectDir string, result *Result) {
	candidates := []string{
		"docker-compose.yml",
		"docker-compose.yaml",
		"compose.yml",
		"compose.yaml",
	}
	for _, name := range candidates {
		path := filepath.Join(projectDir, name)
		if _, err := os.Stat(path); err == nil {
			result.DockerCompose = path
			d.logger.Debug("found Docker Compose", "path", path)
			return
		}
	}
}

func (d *Detector) detectDevContainer(projectDir string, result *Result) {
	path := filepath.Join(projectDir, ".devcontainer", "devcontainer.json")
	if _, err := os.Stat(path); err == nil {
		result.DevContainer = path
		d.logger.Debug("found devcontainer.json", "path", path)
	}
}
