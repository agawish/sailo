package detect

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

var (
	exposeRe    = regexp.MustCompile(`(?i)^\s*EXPOSE\s+(.+)`)
	envPortRe   = regexp.MustCompile(`^(?:PORT|.*_PORT)\s*=\s*(\d+)`)
	scriptPortRe = regexp.MustCompile(`(?:--port|--PORT|-p)\s*[=\s]\s*(\d+)`)
	portEnvInScript = regexp.MustCompile(`PORT[=:]\s*(\d+)`)
)

func (d *Detector) detectPorts(projectDir string, result *Result) {
	seen := map[int]bool{}

	d.detectDockerfilePorts(projectDir, result, seen)
	d.detectComposePorts(projectDir, result, seen)
	d.detectPackageJSONPorts(projectDir, seen)
	d.detectEnvPorts(projectDir, seen)

	var ports []int
	for p := range seen {
		ports = append(ports, p)
	}
	sort.Ints(ports)
	result.Ports = ports
}

func (d *Detector) detectDockerfilePorts(projectDir string, result *Result, seen map[int]bool) {
	if result.Dockerfile == "" {
		return
	}

	data, err := os.ReadFile(result.Dockerfile)
	if err != nil {
		d.logger.Debug("could not read Dockerfile for port detection", "error", err)
		return
	}

	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := scanner.Text()
		matches := exposeRe.FindStringSubmatch(line)
		if matches == nil {
			continue
		}
		tokens := strings.Fields(matches[1])
		for _, tok := range tokens {
			// Strip protocol suffix like /tcp or /udp
			tok = strings.Split(tok, "/")[0]
			port, err := strconv.Atoi(tok)
			if err != nil || !isValidPort(port) {
				continue
			}
			seen[port] = true
		}
	}
}

func (d *Detector) detectComposePorts(projectDir string, result *Result, seen map[int]bool) {
	if result.DockerCompose == "" {
		return
	}

	data, err := os.ReadFile(result.DockerCompose)
	if err != nil {
		d.logger.Debug("could not read compose file for port detection", "error", err)
		return
	}

	var compose map[string]interface{}
	if err := yaml.Unmarshal(data, &compose); err != nil {
		d.logger.Debug("could not parse compose file", "error", err)
		return
	}

	services, ok := compose["services"].(map[string]interface{})
	if !ok {
		return
	}

	for _, svc := range services {
		svcMap, ok := svc.(map[string]interface{})
		if !ok {
			continue
		}
		portsRaw, ok := svcMap["ports"]
		if !ok {
			continue
		}
		portsList, ok := portsRaw.([]interface{})
		if !ok {
			continue
		}
		for _, p := range portsList {
			port := parseComposePort(p)
			if port > 0 && isValidPort(port) {
				seen[port] = true
			}
		}
	}
}

// parseComposePort extracts the container port from a compose port spec.
// Supports: "3000", "3000:3000", "8080:3000/tcp", 3000 (int)
func parseComposePort(v interface{}) int {
	switch p := v.(type) {
	case int:
		return p
	case string:
		// Strip protocol suffix
		p = strings.Split(p, "/")[0]
		// Take the right side of colon (container port)
		parts := strings.Split(p, ":")
		portStr := parts[len(parts)-1]
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return 0
		}
		return port
	}
	return 0
}

func (d *Detector) detectPackageJSONPorts(projectDir string, seen map[int]bool) {
	path := filepath.Join(projectDir, "package.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}

	var pkg map[string]interface{}
	if err := json.Unmarshal(data, &pkg); err != nil {
		return
	}

	scripts, ok := pkg["scripts"].(map[string]interface{})
	if !ok {
		return
	}

	for _, v := range scripts {
		script, ok := v.(string)
		if !ok {
			continue
		}
		for _, re := range []*regexp.Regexp{scriptPortRe, portEnvInScript} {
			for _, match := range re.FindAllStringSubmatch(script, -1) {
				port, err := strconv.Atoi(match[1])
				if err != nil || !isValidPort(port) {
					continue
				}
				seen[port] = true
			}
		}
	}
}

func (d *Detector) detectEnvPorts(projectDir string, seen map[int]bool) {
	path := filepath.Join(projectDir, ".env")
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "#") {
			continue
		}
		matches := envPortRe.FindStringSubmatch(line)
		if matches == nil {
			continue
		}
		port, err := strconv.Atoi(matches[1])
		if err != nil || !isValidPort(port) {
			continue
		}
		seen[port] = true
	}
}

func isValidPort(port int) bool {
	return port >= 1 && port <= 65535
}

// PortSummary returns a human-readable summary of detected ports.
func PortSummary(ports []int) string {
	if len(ports) == 0 {
		return "none detected"
	}
	strs := make([]string, len(ports))
	for i, p := range ports {
		strs[i] = fmt.Sprintf("%d", p)
	}
	return strings.Join(strs, ", ")
}
