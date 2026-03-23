# sAIlo Roadmap

## Phase 1: Foundation (Week 1-2)

Core data model and first working commands.

- [ ] SQLite workspace store (CRUD + auto-migration)
- [ ] Workspace state machine with validated transitions
- [ ] Project detector (Dockerfile, docker-compose.yml, devcontainer.json, language)
- [ ] Port detector (EXPOSE directives, compose ports, package.json scripts)
- [ ] `sailo init` command — detect project and create .sailo.yaml
- [ ] `sailo config show` / `sailo config set` commands
- [ ] Unit tests for state machine, detector, config parsing
- [ ] GitHub Actions CI (build + test on push)

## Phase 2: Core Workspace Lifecycle (Week 3-4)

Create and manage fully isolated workspaces.

- [ ] Moby SDK client wrapper (connect, ping, create, stop, start, remove, exec)
- [ ] Port allocator with SQLite persistence + transaction locking
- [ ] SSH agent forwarding into containers
- [ ] Env var passthrough (explicit allowlist from config)
- [ ] `sailo create` — full flow: detect → allocate ports → build image → create container → git clone → setup
- [ ] `sailo ps` — table output listing all workspaces
- [ ] `sailo exec` — run arbitrary commands in workspace container
- [ ] `sailo stop` / `sailo start` / `sailo rm` — lifecycle management
- [ ] Integration tests with testcontainers-go

## Phase 3: Visibility & Extraction (Week 5-6)

See what agents are doing and ship their work.

- [ ] `sailo logs` — stream container stdout/stderr (with --follow)
- [ ] `sailo diff` — git diff from inside workspace (full diff and --stat)
- [ ] `sailo preview` — open mapped port in default browser
- [ ] `sailo ship` — commit + push + create PR via gh CLI
- [ ] Auto-cleanup daemon for archived workspaces (configurable TTL)
- [ ] Orphaned container detection on `sailo ps`
- [ ] E2E tests for full create → exec → ship → rm workflow

## Phase 4: Polish & Distribution (Week 7-8)

Production-ready UX and distribution.

- [ ] Docker Compose integration for workspace services (postgres, redis)
- [ ] bubbletea TUI dashboard for `sailo ps` (interactive, live updates)
- [ ] bubbletea log viewer for `sailo logs` (scroll, search)
- [ ] Existing Dockerfile/docker-compose.yml/devcontainer.json reuse
- [ ] `.sailo.yaml` project config full support
- [ ] `~/.sailo/config.yaml` user config full support
- [ ] Homebrew formula for macOS distribution
- [ ] `go install` support
- [ ] README with getting started guide
- [ ] Demo GIF/video for README

## Future (v2+)

- [ ] `sailo clone` — fork a workspace mid-flight
- [ ] `sailo conflicts` — detect overlapping file changes across workspaces
- [ ] `sailo replay` — re-run task against updated main
- [ ] `sailo cost` — estimated API token cost per workspace
- [ ] `sailo web` — local web dashboard
- [ ] MCP server — expose workspaces as MCP resources
- [ ] Cloud workspaces — spin up on remote machines
- [ ] Team features — shared workspace visibility
- [ ] Podman/OrbStack runtime support
