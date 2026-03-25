// Package container wraps the Moby SDK to manage Docker containers
// for sailo workspaces.
package container

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"sort"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/moby/term"
)

// Client wraps the Docker/Moby API client.
type Client struct {
	docker *client.Client
	logger *slog.Logger
}

// CreateOptions configures a new workspace container.
type CreateOptions struct {
	WorkspaceID string            // used for container naming: sailo-<id>
	Image       string
	Ports       map[int]int       // container port → host port
	SSHAuthSock string
	EnvVars     map[string]string
	GHConfigDir string            // path to gh config (mounted read-only)
}

// ContainerState holds basic container status info.
type ContainerState struct {
	Running bool
	Status  string
}

// NewClient creates a Docker client connected to the local daemon.
func NewClient(logger *slog.Logger) (*Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("create docker client: %w", err)
	}
	return &Client{docker: cli, logger: logger}, nil
}

// Ping checks that the Docker daemon is reachable.
func (c *Client) Ping(ctx context.Context) error {
	_, err := c.docker.Ping(ctx)
	if err != nil {
		return fmt.Errorf("docker daemon not reachable: %w\n\nIs Docker running? Try: docker info", err)
	}
	return nil
}

// CreateWorkspace creates and starts a container for a workspace.
func (c *Client) CreateWorkspace(ctx context.Context, opts CreateOptions) (string, error) {
	if err := c.ensureImage(ctx, opts.Image); err != nil {
		return "", fmt.Errorf("ensure image %s: %w", opts.Image, err)
	}

	// Build env vars
	envSlice := buildEnvSlice(opts.EnvVars)

	// Build exposed ports and port bindings
	exposedPorts := nat.PortSet{}
	for cp := range opts.Ports {
		exposedPorts[nat.Port(fmt.Sprintf("%d/tcp", cp))] = struct{}{}
	}

	containerCfg := &container.Config{
		Image:        opts.Image,
		Env:          envSlice,
		Cmd:          []string{"sleep", "infinity"},
		Tty:          true,
		ExposedPorts: exposedPorts,
		Labels: map[string]string{
			"managed-by": "sailo",
		},
	}

	hostCfg := &container.HostConfig{
		PortBindings: buildPortBindings(opts.Ports),
		Binds:        buildBinds(opts.SSHAuthSock, opts.GHConfigDir),
	}

	name := "sailo-" + opts.WorkspaceID

	resp, err := c.docker.ContainerCreate(ctx, containerCfg, hostCfg, nil, nil, name)
	if err != nil {
		return "", fmt.Errorf("create container: %w", err)
	}

	if err := c.docker.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		c.docker.ContainerRemove(ctx, resp.ID, container.RemoveOptions{Force: true})
		return "", fmt.Errorf("start container: %w", err)
	}

	c.logger.Info("container created and started", "id", resp.ID[:12], "name", name)
	return resp.ID, nil
}

// StopContainer stops a running container.
func (c *Client) StopContainer(ctx context.Context, containerID string) error {
	timeout := 10
	if err := c.docker.ContainerStop(ctx, containerID, container.StopOptions{Timeout: &timeout}); err != nil {
		return fmt.Errorf("stop container %s: %w", shortID(containerID), err)
	}
	c.logger.Info("container stopped", "id", shortID(containerID))
	return nil
}

// StartContainer starts a stopped container.
func (c *Client) StartContainer(ctx context.Context, containerID string) error {
	if err := c.docker.ContainerStart(ctx, containerID, container.StartOptions{}); err != nil {
		return fmt.Errorf("start container %s: %w", shortID(containerID), err)
	}
	c.logger.Info("container started", "id", shortID(containerID))
	return nil
}

// RemoveContainer removes a container and its volumes.
func (c *Client) RemoveContainer(ctx context.Context, containerID string) error {
	if err := c.docker.ContainerRemove(ctx, containerID, container.RemoveOptions{
		Force:         true,
		RemoveVolumes: true,
	}); err != nil {
		return fmt.Errorf("remove container %s: %w", shortID(containerID), err)
	}
	c.logger.Info("container removed", "id", shortID(containerID))
	return nil
}

// ExecInContainer runs a command non-interactively inside a container.
// Returns an error if the command exits with a non-zero code, including its output.
func (c *Client) ExecInContainer(ctx context.Context, containerID string, cmd []string) error {
	_, err := c.execRaw(ctx, containerID, cmd)
	return err
}

// ExecWithOutput runs a command non-interactively and returns its stdout.
func (c *Client) ExecWithOutput(ctx context.Context, containerID string, cmd []string) (string, error) {
	return c.execRaw(ctx, containerID, cmd)
}

func (c *Client) execRaw(ctx context.Context, containerID string, cmd []string) (string, error) {
	execCfg := container.ExecOptions{
		Cmd:          cmd,
		AttachStdout: true,
		AttachStderr: true,
	}
	execID, err := c.docker.ContainerExecCreate(ctx, containerID, execCfg)
	if err != nil {
		return "", fmt.Errorf("create exec: %w", err)
	}

	resp, err := c.docker.ContainerExecAttach(ctx, execID.ID, container.ExecStartOptions{})
	if err != nil {
		return "", fmt.Errorf("attach exec: %w", err)
	}
	defer resp.Close()

	output, _ := io.ReadAll(resp.Reader)

	inspectResp, err := c.docker.ContainerExecInspect(ctx, execID.ID)
	if err != nil {
		return "", fmt.Errorf("inspect exec: %w", err)
	}
	if inspectResp.ExitCode != 0 {
		return "", fmt.Errorf("command %v exited with code %d: %s", cmd, inspectResp.ExitCode, strings.TrimSpace(string(output)))
	}

	return strings.TrimSpace(string(output)), nil
}

// ExecInteractive attaches stdin/stdout/stderr to a command in the container.
// Uses raw terminal mode with proper cleanup on signal/crash.
func (c *Client) ExecInteractive(ctx context.Context, containerID string, cmd []string) error {
	execCfg := container.ExecOptions{
		Cmd:          cmd,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
		WorkingDir:   "/workspace",
	}

	execID, err := c.docker.ContainerExecCreate(ctx, containerID, execCfg)
	if err != nil {
		return fmt.Errorf("create interactive exec: %w", err)
	}

	resp, err := c.docker.ContainerExecAttach(ctx, execID.ID, container.ExecStartOptions{Tty: true})
	if err != nil {
		return fmt.Errorf("attach interactive exec: %w", err)
	}
	defer resp.Close()

	// Set terminal to raw mode with signal-safe restore
	inFd, _ := term.GetFdInfo(os.Stdin)
	oldState, err := term.SetRawTerminal(inFd)
	if err == nil {
		// Restore on normal exit
		defer term.RestoreTerminal(inFd, oldState)

		// Also restore on signals
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt)
		go func() {
			<-sigCh
			term.RestoreTerminal(inFd, oldState)
		}()
		defer signal.Stop(sigCh)
	}

	// Bidirectional copy
	go io.Copy(resp.Conn, os.Stdin)
	io.Copy(os.Stdout, resp.Reader)

	return nil
}

// InspectContainer checks if a container is actually running.
func (c *Client) InspectContainer(ctx context.Context, containerID string) (ContainerState, error) {
	info, err := c.docker.ContainerInspect(ctx, containerID)
	if err != nil {
		return ContainerState{}, fmt.Errorf("inspect container %s: %w", shortID(containerID), err)
	}
	return ContainerState{
		Running: info.State.Running,
		Status:  info.State.Status,
	}, nil
}

// Close closes the underlying Docker client.
func (c *Client) Close() error {
	return c.docker.Close()
}

func (c *Client) ensureImage(ctx context.Context, img string) error {
	_, _, err := c.docker.ImageInspectWithRaw(ctx, img)
	if err == nil {
		c.logger.Debug("image already present", "image", img)
		return nil
	}

	c.logger.Info("pulling image", "image", img)
	reader, err := c.docker.ImagePull(ctx, img, image.PullOptions{})
	if err != nil {
		return fmt.Errorf("pull image %s: %w", img, err)
	}
	defer reader.Close()
	io.Copy(io.Discard, reader)
	return nil
}

// buildEnvSlice converts a map to KEY=VALUE slice.
func buildEnvSlice(env map[string]string) []string {
	if len(env) == 0 {
		return nil
	}
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	result := make([]string, 0, len(env))
	for _, k := range keys {
		result = append(result, fmt.Sprintf("%s=%s", k, env[k]))
	}
	return result
}

// buildPortBindings converts container→host port map to Moby's nat.PortMap.
// Binds to 127.0.0.1 only for security.
func buildPortBindings(ports map[int]int) nat.PortMap {
	portMap := nat.PortMap{}
	for containerPort, hostPort := range ports {
		cp := nat.Port(fmt.Sprintf("%d/tcp", containerPort))
		portMap[cp] = []nat.PortBinding{
			{HostIP: "127.0.0.1", HostPort: fmt.Sprintf("%d", hostPort)},
		}
	}
	return portMap
}

// buildBinds creates volume mount strings for SSH agent and gh config.
func buildBinds(sshAuthSock, ghConfigDir string) []string {
	var binds []string
	if sshAuthSock != "" {
		binds = append(binds, fmt.Sprintf("%s:/run/ssh-agent.sock", sshAuthSock))
	}
	if ghConfigDir != "" {
		binds = append(binds, fmt.Sprintf("%s:/root/.config/gh:ro", ghConfigDir))
	}
	return binds
}

func shortID(id string) string {
	if len(id) > 12 {
		return id[:12]
	}
	return id
}
