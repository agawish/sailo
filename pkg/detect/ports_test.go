package detect

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/agawish/sailo/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func detectInDir(t *testing.T, dir string) *Result {
	t.Helper()
	d := NewDetector(testutil.NewTestLogger())
	result, err := d.Detect(dir)
	require.NoError(t, err)
	return result
}

func TestDetectPorts_Dockerfile(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "Dockerfile"), []byte(`
FROM node:22
EXPOSE 3000 8080
CMD ["node", "server.js"]
`), 0644))

	result := detectInDir(t, dir)
	assert.Equal(t, []int{3000, 8080}, result.Ports)
}

func TestDetectPorts_DockerfileWithProtocol(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "Dockerfile"), []byte(`
FROM nginx
EXPOSE 80/tcp 443/udp
`), 0644))

	result := detectInDir(t, dir)
	assert.Equal(t, []int{80, 443}, result.Ports)
}

func TestDetectPorts_DockerfileMultipleLines(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "Dockerfile"), []byte(`
FROM python:3.12
EXPOSE 5000
EXPOSE 5432
`), 0644))

	result := detectInDir(t, dir)
	assert.Equal(t, []int{5000, 5432}, result.Ports)
}

func TestDetectPorts_DockerCompose(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "docker-compose.yml"), []byte(`
services:
  web:
    image: nginx
    ports:
      - "3000:3000"
      - "8080:8080"
  db:
    image: postgres
    ports:
      - "5432:5432"
`), 0644))

	result := detectInDir(t, dir)
	assert.Contains(t, result.Ports, 3000)
	assert.Contains(t, result.Ports, 8080)
	assert.Contains(t, result.Ports, 5432)
}

func TestDetectPorts_DockerComposeHostOnly(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "docker-compose.yml"), []byte(`
services:
  app:
    image: node
    ports:
      - "3000"
`), 0644))

	result := detectInDir(t, dir)
	assert.Contains(t, result.Ports, 3000)
}

func TestDetectPorts_PackageJSON(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "package.json"), []byte(`{
  "scripts": {
    "start": "node server.js --port 3000",
    "dev": "vite --port 5173"
  }
}`), 0644))

	result := detectInDir(t, dir)
	assert.Contains(t, result.Ports, 3000)
	assert.Contains(t, result.Ports, 5173)
}

func TestDetectPorts_PackageJSONPortEnv(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "package.json"), []byte(`{
  "scripts": {
    "start": "PORT=4000 node server.js"
  }
}`), 0644))

	result := detectInDir(t, dir)
	assert.Contains(t, result.Ports, 4000)
}

func TestDetectPorts_EnvFile(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".env"), []byte(`
# Configuration
PORT=3000
DATABASE_URL=postgres://localhost
APP_PORT=8080
`), 0644))

	result := detectInDir(t, dir)
	assert.Contains(t, result.Ports, 3000)
	assert.Contains(t, result.Ports, 8080)
}

func TestDetectPorts_EnvFileIgnoresComments(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".env"), []byte(`
# PORT=9999
PORT=3000
`), 0644))

	result := detectInDir(t, dir)
	assert.Equal(t, []int{3000}, result.Ports)
}

func TestDetectPorts_MultipleSources(t *testing.T) {
	dir := t.TempDir()
	// Dockerfile with port 3000
	require.NoError(t, os.WriteFile(filepath.Join(dir, "Dockerfile"), []byte(`
FROM node:22
EXPOSE 3000
`), 0644))
	// .env with same port 3000 and another port
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".env"), []byte(`
PORT=3000
API_PORT=8080
`), 0644))

	result := detectInDir(t, dir)
	// Should be deduplicated and sorted
	assert.Equal(t, []int{3000, 8080}, result.Ports)
}

func TestDetectPorts_NoFiles(t *testing.T) {
	dir := t.TempDir()

	result := detectInDir(t, dir)
	assert.Empty(t, result.Ports)
}

func TestPortSummary(t *testing.T) {
	assert.Equal(t, "none detected", PortSummary(nil))
	assert.Equal(t, "none detected", PortSummary([]int{}))
	assert.Equal(t, "3000", PortSummary([]int{3000}))
	assert.Equal(t, "3000, 8080", PortSummary([]int{3000, 8080}))
}
