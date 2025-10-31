# HomeKit CLI Developer Guide

This guide explains how to work on the HomeKit CLI codebase and the tooling that ships with it.

## Prerequisites

- Go 1.23 or newer (matches `go.mod`).
- Docker (required only for the optional development container workflow).
- Recommended: run `make deps` to install extra tooling such as `golangci-lint`.

Verify the local toolchain and bootstrap the workspace:

```bash
make init
make deps
```

## Repository Layout

```
homekit-cli/
├── assets/                # Embedded scripts and templates (go:embed)
├── cmd/homekit/           # Cobra entrypoint and version wiring
├── config/                # Example runtime configuration
├── docker/                # Portable development container assets
├── docs/                  # Developer documentation (this guide, spec, releases)
├── internal/
│   ├── assets/            # Asset manager with override + checksum support
│   ├── commands/          # CLI command groups (script, assets, docker, sys, ...)
│   ├── core/              # Runtime bootstrap (config, logging, dry-run)
│   ├── exec/              # External process runner
│   ├── plugins/           # Plugin discovery helpers
│   ├── shell/             # mvdan/sh-backed interpreter for embedded scripts
│   ├── templating/        # Helpers around Go text/template
│   └── ui/                # Lightweight terminal prompts
├── make/                  # Modular make fragments (Go, Docker, dist, bootstrap)
├── pkg/utils/             # Shared helper utilities
├── .github/workflows/     # CI pipelines
├── go.mod / go.sum        # Module definition
├── main.go                # Program entrypoint
└── Makefile               # Aggregates the make fragments
```

## Command Overview

The CLI surface is built with Cobra and initialised through `cmd/homekit/root.go`.

- Global flags: `--config`, `--log-level`, `--log-format`, `--no-color`, `--dry-run`.
- `homekit version` – emit build metadata (`table` or `json` output).
- `homekit script run|list` – execute local binaries or embedded scripts.
- `homekit assets list|extract|verify` – inspect bundled assets and export overrides.
- `homekit template render` – render embedded templates with merged YAML data.
- `homekit docker prune|images update` – quality-of-life Docker helpers.
- `homekit sys health` – show basic system metrics (load, memory, disk).
- `homekit plugins list` – discover executables prefixed with `homekit-cli-`.

Each subcommand retrieves the initialised runtime from context (see `internal/core/runtime.go`) to share configuration, logging, and dry-run settings.

## Build and Run

Compile the CLI locally:

```bash
make build
./dist/homekit --log-level debug --dry-run
```

Run straight from source during development:

```bash
go run ./cmd/homekit --log-level debug --dry-run
```

Generate cross-platform binaries and checksums:

```bash
make dist VERSION=0.1.0
# outputs dist/homekit-<goos>-<goarch>[.exe] + matching .sha256 files
```

## Configuration & Runtime Behaviour

- Default config path: `${XDG_CONFIG_HOME}/homekit/config.yaml` (override with `--config`).
- Reference file: `config/config.example.yaml`.
- Recognised keys today: `asset_overrides`, `plugin_paths`, `temp_dir`, `log_level`.
- Environment variables with the `HOMEKIT_` prefix take precedence (`HOMEKIT_LOG_LEVEL=debug`).
- Set `dry-run` via the flag to simulate side effects while still logging intent.

### Asset Overrides

If `asset_overrides` is configured (defaults to `~/.config/homekit/assets`), files placed under the matching namespace shadow embedded versions. For example:

```
~/.config/homekit/assets/
├── scripts/docker_prune_safe.sh
└── templates/docker-compose.yaml.tmpl
```

## Embedded Assets & Templates

Current embedded content:

- `assets/scripts/docker_prune_safe.sh` – placeholder shell script executed through the embedded interpreter.
- `assets/templates/docker-compose.yaml.tmpl` – minimal Compose template surfaced via `homekit template render`.

Useful commands:

```bash
homekit script list
homekit script run --embedded docker_prune_safe.sh
homekit assets extract templates docker-compose.yaml.tmpl ./out
homekit template render docker-compose.yaml.tmpl --data values.yaml --output ./docker-compose.yaml
```

### Programmatic Execution Examples

Embedded scripts can be invoked directly from code without writing them to disk:

```go
import (
    "fmt"
    "time"

    "github.com/homekit/homekit-cli/internal/assets"
    "github.com/homekit/homekit-cli/internal/shell"
)

mgr := assets.NewManager(assets.Embedded(), overrideDirectory(rt.Config))
rc, err := mgr.Open("scripts", "docker_prune_safe.sh")
if err != nil {
    return fmt.Errorf("open embedded script: %w", err)
}
defer rc.Close()

result, err := shell.Run(ctx, "docker_prune_safe.sh", rc, shell.Options{
    Args:          []string{"--dry-run"},
    Env:           map[string]string{"PRUNE_CONFIRM": "true"},
    Dir:           "/tmp",
    Timeout:       30 * time.Second,
    CaptureOutput: true,
})
if err != nil {
    return err
}
fmt.Print(result.Stdout)
```

When dispatching a single local command (or a script with a shebang), use the process runner:

```go
import (
    "context"
    "fmt"
    "os/exec"
    "time"

    executor "github.com/homekit/homekit-cli/internal/exec"
)

if _, err := exec.LookPath("git"); err != nil {
    return fmt.Errorf("git not installed: %w", err)
}

res, err := executor.Run(ctx, executor.Spec{
    Command:       "git",
    Args:          []string{"status", "--short"},
    Dir:           "/path/to/repo",
    Timeout:       10 * time.Second,
    CaptureOutput: true,
})
if err != nil {
    return err
}
fmt.Print(res.Stdout)
```

## Plugin Workflows

Plugins are external executables discoverable on `$PATH` (and any additional directories in `plugin_paths`) that start with the prefix `homekit-cli-` by default.

```bash
homekit plugins list
homekit plugins list --prefix hk --path ./build/plugins
```

The manager defined in `internal/plugins/manager.go` filters files for executability, removes the prefix for display, and returns descriptors that can be launched via `ExecProxy`.

## Development Container

The `docker/` directory ships a portable dev environment based on `ubuntu:24.04`:

- Installs Go (default `GO_VERSION=1.23.2`) and configures `GOPATH=/go`.
- Adds core build packages plus `golangci-lint`.
- Creates a non-root `dev` user with configurable UID/GID.

Common workflows:

```bash
make docker-build IMAGE=ghcr.io/<user>/homekit-cli-dev TAG=latest
make docker-shell IMAGE=ghcr.io/<user>/homekit-cli-dev TAG=latest
make docker-push IMAGE=ghcr.io/<user>/homekit-cli-dev TAG=latest
```

`docker/compose.dev.yml` provides a convenience wrapper that mounts the repository inside the container at `/work/homekit-cli`.

## Testing & Linting

```bash
make fmt          # gofmt + go fmt
make lint         # runs golangci-lint if available
make test         # writes coverage summary to dist/coverage.out
```

Run `GOFLAGS='-tags=integration'` (or similar) to forward extra flags through the make targets.

## Release Flow

1. Update docs/changelog as needed.
2. Build release artefacts: `make dist VERSION=1.0.0`.
3. Inspect the binaries and checksums under `dist/`.
4. Tag the commit and push: `git tag v1.0.0 && git push --tags`.
5. Publish container image if required: `make docker-push IMAGE=ghcr.io/<user>/homekit-cli-dev TAG=1.0.0`.
6. Attach artefacts to a GitHub release or allow CI to publish them.

The `dist.release` target remains a placeholder for future CI automation (e.g. triggering a workflow that consumes the generated artefacts).
