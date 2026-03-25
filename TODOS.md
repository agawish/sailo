# TODOs

## Pre-built sailo base images
**What:** Create sailo/node:22, sailo/go:1.22, etc. with git+ssh pre-installed.
**Why:** Eliminates 10-30s `apt-get install git` on every `sailo create` for slim images.
**Depends on:** Phase 2 complete.
**Target:** Phase 4.

## Orphaned container detection on `sailo ps`
**What:** InspectContainer for each listed workspace and update stale states when container has crashed or been removed externally.
**Why:** DB can say "running" when the container is dead. `sailo ps` should show truth.
**Depends on:** Phase 2 InspectContainer method.
**Target:** Phase 3.

## SIGINT graceful cleanup during `sailo create`
**What:** Register signal handler that triggers cleanup() if user Ctrl-C's mid-create.
**Why:** Prevents orphaned containers and stale workspace records.
**Depends on:** Phase 2 Create flow with cleanup() closure.
**Target:** Phase 4.
