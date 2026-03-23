# sAIlo

**Workspace isolation layer for AI coding agents.**

When multiple AI agents work on the same codebase, everything collides — ports, git state, file locks.
sAIlo creates isolated Docker workspaces so each agent gets its own container, git clone, and port range.
Agent-agnostic — attach Claude Code, Cursor, Codex, or any tool.

```
$ sailo create "add dark mode to settings page" --from=main
  ✓ Workspace ws-7f3a created (port 3007, branch sailo/ws-7f3a/dark-mode)

$ sailo create "fix pagination bug on /users" --from=main
  ✓ Workspace ws-9b1c created (port 3008, branch sailo/ws-9b1c/fix-pagination)

$ sailo ps
  ID       TASK                          STATUS   PORT   BRANCH
  ws-7f3a  add dark mode to settings     running  3007   sailo/ws-7f3a/dark-mode
  ws-9b1c  fix pagination bug on /users  running  3008   sailo/ws-9b1c/fix-pagination

$ sailo exec ws-7f3a -- claude-code
$ sailo diff ws-7f3a
$ sailo ship ws-7f3a
  ✓ PR #47 created: "Add dark mode to settings page"
```

## What sAIlo Does

- **Isolated workspaces**: Each agent gets its own Docker container with full git clone
- **No port conflicts**: Automatic port mapping — container:3000 gets a unique host port
- **No git conflicts**: Each workspace works on its own branch in its own clone
- **Ship to PR**: `sailo ship` commits, pushes, and creates a pull request
- **Agent-agnostic**: Works with any AI coding tool (Claude Code, Cursor, Codex, bash)
- **Reuses your Dockerfile**: Detects and reuses existing project Docker configuration

## Architecture

```
sailo CLI (Go binary)
  ├── Workspace Manager   → container lifecycle via Moby/Docker API
  ├── Port Allocator      → non-conflicting host port assignment
  ├── Git Manager         → git clone --depth=1 into Docker volumes
  ├── Project Detector    → reuses existing Dockerfile/docker-compose.yml
  ├── Credential Manager  → SSH agent forwarding + env var passthrough
  └── TUI Dashboard       → interactive workspace overview
```

## Install

```bash
# From source
go install github.com/agawish/sailo/cmd/sailo@latest

# Or build locally
git clone https://github.com/agawish/sailo.git
cd sailo
make build
./bin/sailo --help
```

**Requirements**: Docker (or any Moby-compatible engine) must be running.

## Quick Start

```bash
# Initialize sailo in your project
cd your-project
sailo init

# Create an isolated workspace
sailo create "add user authentication"

# Attach an agent
sailo exec ws-abc1 -- claude-code

# See what changed
sailo diff ws-abc1

# Ship it
sailo ship ws-abc1
```

## Commands

| Command | Description |
|---------|-------------|
| `sailo init` | Initialize sailo in current project |
| `sailo create <task>` | Create isolated workspace |
| `sailo ps` | List all workspaces |
| `sailo exec <id> -- <cmd>` | Run command in workspace |
| `sailo logs <id>` | Stream workspace logs |
| `sailo diff <id>` | Show changes in workspace |
| `sailo preview <id>` | Open mapped port in browser |
| `sailo ship <id>` | Extract changes into PR |
| `sailo stop <id>` | Stop workspace (preserve state) |
| `sailo start <id>` | Resume stopped workspace |
| `sailo rm <id>` | Remove workspace and cleanup |
| `sailo config` | Manage configuration |

## Configuration

### Project config (`.sailo.yaml`)

```yaml
version: 1
image: node:22-slim
services:
  - postgres:16
ports:
  3000: auto
setup:
  - npm install
test: npm test
```

### User config (`~/.sailo/config.yaml`)

```yaml
defaults:
  from: main
  port_range: 3001-3999
env_passthrough:
  - ANTHROPIC_API_KEY
  - GITHUB_TOKEN
```

## Security

- SSH keys are forwarded via SSH agent (never copied to disk)
- Environment variables use an explicit allowlist
- Workspace containers have no access to host filesystem
- Each workspace uses an isolated Docker bridge network

## Development

```bash
make build    # Build binary
make test     # Run tests
make lint     # Run linters
```

## License

MIT
