package detect

func (d *Detector) detectPorts(projectDir string, result *Result) {
	// TODO: parse Dockerfile EXPOSE directives
	// TODO: parse docker-compose.yml ports
	// TODO: parse package.json scripts for common port patterns
	// TODO: parse .env files for PORT variables

	// Default: no ports detected (user must configure in .sailo.yaml)
}
