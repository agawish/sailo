package container

import (
	"testing"

	"github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/assert"
)

func TestBuildEnvSlice(t *testing.T) {
	env := map[string]string{
		"NODE_ENV":         "development",
		"ANTHROPIC_API_KEY": "sk-test",
		"SSH_AUTH_SOCK":    "/run/ssh-agent.sock",
	}
	result := buildEnvSlice(env)

	assert.Len(t, result, 3)
	// Should be sorted by key
	assert.Equal(t, "ANTHROPIC_API_KEY=sk-test", result[0])
	assert.Equal(t, "NODE_ENV=development", result[1])
	assert.Equal(t, "SSH_AUTH_SOCK=/run/ssh-agent.sock", result[2])
}

func TestBuildEnvSlice_Empty(t *testing.T) {
	assert.Nil(t, buildEnvSlice(nil))
	assert.Nil(t, buildEnvSlice(map[string]string{}))
}

func TestBuildPortBindings(t *testing.T) {
	ports := map[int]int{3000: 3007, 8080: 3008}
	result := buildPortBindings(ports)

	assert.Len(t, result, 2)

	bindings3000 := result[nat.Port("3000/tcp")]
	assert.Len(t, bindings3000, 1)
	assert.Equal(t, "127.0.0.1", bindings3000[0].HostIP)
	assert.Equal(t, "3007", bindings3000[0].HostPort)

	bindings8080 := result[nat.Port("8080/tcp")]
	assert.Len(t, bindings8080, 1)
	assert.Equal(t, "127.0.0.1", bindings8080[0].HostIP)
	assert.Equal(t, "3008", bindings8080[0].HostPort)
}

func TestBuildPortBindings_Empty(t *testing.T) {
	result := buildPortBindings(nil)
	assert.Empty(t, result)
}

func TestBuildBinds(t *testing.T) {
	binds := buildBinds("/tmp/agent.sock", "/home/user/.config/gh")
	assert.Equal(t, []string{
		"/tmp/agent.sock:/run/ssh-agent.sock",
		"/home/user/.config/gh:/root/.config/gh:ro",
	}, binds)
}

func TestBuildBinds_SSHOnly(t *testing.T) {
	binds := buildBinds("/tmp/agent.sock", "")
	assert.Equal(t, []string{"/tmp/agent.sock:/run/ssh-agent.sock"}, binds)
}

func TestBuildBinds_Empty(t *testing.T) {
	binds := buildBinds("", "")
	assert.Nil(t, binds)
}

func TestShortID(t *testing.T) {
	assert.Equal(t, "abc123def456", shortID("abc123def4567890"))
	assert.Equal(t, "short", shortID("short"))
}
