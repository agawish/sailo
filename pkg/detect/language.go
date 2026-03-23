package detect

import (
	"os"
	"path/filepath"
)

// languageMap maps marker files to language names and default base images.
var languageMap = map[string]struct {
	language  string
	baseImage string
}{
	"package.json":    {"node", "node:22-slim"},
	"go.mod":          {"go", "golang:1.22-bookworm"},
	"Cargo.toml":      {"rust", "rust:1-slim-bookworm"},
	"requirements.txt": {"python", "python:3.12-slim"},
	"pyproject.toml":  {"python", "python:3.12-slim"},
	"Gemfile":         {"ruby", "ruby:3.3-slim"},
	"pom.xml":         {"java", "eclipse-temurin:21-jdk"},
	"build.gradle":    {"java", "eclipse-temurin:21-jdk"},
}

func (d *Detector) detectLanguage(projectDir string, result *Result) {
	for marker, info := range languageMap {
		path := filepath.Join(projectDir, marker)
		if _, err := os.Stat(path); err == nil {
			result.Language = info.language
			if result.BaseImage == "" {
				result.BaseImage = info.baseImage
			}
			d.logger.Debug("detected language", "language", info.language, "marker", marker)
			return
		}
	}

	// fallback
	if result.BaseImage == "" {
		result.BaseImage = "ubuntu:24.04"
	}
}
