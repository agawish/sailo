# CLAUDE.md - sAIlo Project Guide

## What is sAIlo?

sAIlo (sailo) is a CLI tool that creates isolated Docker workspaces for AI coding agents.
Each workspace gets its own container, git clone, port range, and SSH forwarding.
Agent-agnostic — attach Claude Code, Cursor, Codex, or any tool.

## Commands

- Build: `make build`
- Test: `make test`
- Lint: `make lint`
- Install locally: `make install`
- Run: `./bin/sailo --help`

## Architecture

```
sailo CLI (Go binary)
  ├── Workspace Manager   → creates/destroys containers via Moby SDK
  ├── Port Allocator      → assigns non-conflicting host ports per workspace
  ├── Git Manager         → git clone --depth=1 into Docker volumes
  ├── Project Detector    → reuses existing Dockerfile/docker-compose.yml
  ├── Credential Manager  → SSH agent forwarding + env var passthrough
  └── TUI                 → bubbletea dashboard for sailo ps
```

Workspaces are Docker containers. Each gets a full shallow clone, isolated ports,
and forwarded SSH agent. The host filesystem and git state are never touched.

## Code Style

- Go 1.22+ with standard library where possible
- 2-space indentation for non-Go config files (yaml, json)
- Go files use gofmt standard (tabs)
- Package names: short, lowercase, no underscores
- Error handling: wrap with `fmt.Errorf("context: %w", err)` — never swallow errors
- Logging: use `log/slog` structured logger everywhere
- CLI output: use `fmt.Fprintf(cmd.OutOrStdout(), ...)` for testability
- Tests: use testify for assertions, table-driven tests for multiple cases

## Project Structure

- `cmd/sailo/` — CLI entry point and command definitions
- `pkg/` — public packages (workspace, container, port, git, detect, creds, config, tui)
- `internal/` — private packages (test utilities)
- Each `pkg/` package should be independently testable

## Key Design Decisions

1. **Agent-agnostic**: sailo creates workspaces, users attach any agent via `sailo exec`
2. **Moby SDK**: uses Docker Engine API directly (github.com/docker/docker/client)
3. **Full git clone**: each workspace gets `git clone --depth=1` into a Docker volume
4. **SSH agent forwarding**: credentials are forwarded, never copied to disk
5. **SQLite state**: workspace metadata stored in ~/.sailo/workspaces.db
6. **Reuse Dockerfiles**: detect and reuse existing project Dockerfile/docker-compose.yml

## Git Conventions

- Present tense commit messages ("add workspace manager", not "added")
- Branch naming: feature/<name>, fix/<name>
- Keep commits focused — one logical change per commit
