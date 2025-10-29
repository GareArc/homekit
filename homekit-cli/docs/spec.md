# HomeKit CLI — Project Structure

This specification captures the current layout and major subsystems of the HomeKit CLI codebase.

## Root Layout

```
homekit/
├── LICENSE
└── homekit-cli/       # Go CLI project
```

## homekit-cli Layout (high-level)

```
homekit-cli/
├── .dockerignore
├── .github/workflows/
├── assets/
│   ├── embed.go
│   ├── scripts/docker_prune_safe.sh
│   └── templates/docker-compose.yaml.tmpl
├── cmd/homekit/       # Cobra root command and version wiring
├── config/            # Example configuration
├── docker/            # Portable development container assets
├── docs/              # Project documentation (this spec, dev guide, releases)
├── internal/
│   ├── assets/        # Asset manager and override helpers
│   ├── commands/      # Cobra subcommands
│   ├── core/          # Runtime bootstrap (config + logging)
│   ├── exec/          # External process runner
│   ├── plugins/       # Plugin discovery
│   ├── shell/         # Embedded shell interpreter (mvdan/sh)
│   ├── templating/    # text/template helpers
│   └── ui/            # Minimal interactive prompts
├── make/              # Modular Make fragments
├── pkg/utils/         # Shared utility helpers
├── go.mod / go.sum
├── main.go
└── Makefile
```

## Directory Purpose Summary

| Path                   | Purpose                                                            |
| ---------------------- | ------------------------------------------------------------------ |
| `.github/workflows/`   | Build and release CI pipelines.                                    |
| `assets/`              | Scripts and templates bundled with the binary via `go:embed`.      |
| `cmd/homekit/`         | Cobra root command, persistent flags, and runtime bootstrap glue.  |
| `config/`              | Reference configuration (`config.example.yaml`).                   |
| `docker/`              | Development container (`Dockerfile.dev`, `compose.dev.yml`, docs). |
| `docs/`                | Contributor-facing documentation.                                  |
| `internal/assets/`     | Asset manager supporting overrides and checksums.                  |
| `internal/commands/`   | Implementation of CLI command groups (`script`, `assets`, etc.).   |
| `internal/core/`       | Configuration loading, logging setup, and runtime context.         |
| `internal/exec/`       | Thin wrapper around `os/exec` with timeout and dry-run support.    |
| `internal/plugins/`    | Discovers `homekit-cli-*` executables as plugins.                  |
| `internal/shell/`      | mvdan/sh-backed interpreter for embedded scripts.                  |
| `internal/templating/` | text/template renderer utilities.                                  |
| `internal/ui/`         | Terminal prompter helpers.                                         |
| `make/`                | Make fragments for Go builds, dist packaging, Docker, bootstrap.   |
| `pkg/utils/`           | Reusable helpers (env merging, key/value parsing).                 |

## CLI Command Surface

- `homekit` root flags: `--config`, `--log-level`, `--log-format`, `--no-color`, `--dry-run`.
- `homekit version`: print build metadata wired via `-ldflags`.
- `homekit script run|list`: execute local commands or embedded scripts via `internal/shell`.
- `homekit assets list|extract|verify`: inspect and export embedded assets with override support.
- `homekit template render`: render embedded templates with merged YAML data files.
- `homekit docker prune|images update`: quality-of-life Docker helpers.
- `homekit sys health`: show lightweight system metrics (CPU load, memory, disk).
- `homekit plugins list`: discover external executables matching the plugin prefix.

Each subcommand relies on the shared runtime initialized in `cmd/homekit/root.go`, exposing structured logging, config, and dry-run behaviour.

## Configuration & Overrides

- Default config path: `${XDG_CONFIG_HOME}/homekit/config.yaml` (see `core.DefaultConfigPath`).
- Example configuration lives in `config/config.example.yaml`.
- Environment variables prefixed with `HOMEKIT_` override config keys (`viper.AutomaticEnv`).
- Asset overrides are loaded from `asset_overrides` (defaults to `~/.config/homekit/assets`), allowing local files to shadow embedded content.
- Additional plugin search paths can be provided via the `plugin_paths` array.

## Embedded Assets

Current embedded artefacts (`assets/`):

| Namespace   | Files                      | Notes                                             |
| ----------- | -------------------------- | ------------------------------------------------- |
| `scripts`   | `docker_prune_safe.sh`     | Placeholder shell script executed via interpreter |
| `templates` | `docker-compose.yaml.tmpl` | Minimal compose template demonstrating rendering  |

The manager in `internal/assets/manager.go`:

- lists assets across embedded and override directories,
- exports assets to disk with executable permissions,
- verifies assets by SHA-256 checksum.

CLI usage examples:

```bash
homekit script list
homekit script run ./local-script.sh --timeout 30s
homekit script run --embedded docker_prune_safe.sh
homekit template render docker-compose.yaml.tmpl --data values.yaml --output ./docker-compose.yaml
homekit assets extract templates docker-compose.yaml.tmpl ./dist/templates
```

## Make Targets & Tooling

The top-level `Makefile` delegates to fragments in `make/`:

- `make init` – prerequisite checks (`bootstrap.mk`).
- `make deps|tidy|fmt|lint|test|generate` – common Go workflows (`go.mk`).
- `make build` – build the CLI into `dist/homekit` with version metadata.
- `make run PROFILE=<profile>` – run via `go run` with runtime flags.
- `make dist` – cross-platform builds and checksums (`dist.mk`).
- `make docker-build|docker-shell|docker-push` – dev-container lifecycle (`docker.mk`).

Environment variables such as `PROFILE`, `LOG_LEVEL`, `LOG_FORMAT`, `VERSION`, and `PLATFORMS` are configurable per invocation.

## Development Container

`docker/Dockerfile.dev` provides a portable Go toolchain image:

- Base: `ubuntu:24.04`.
- Installs common build prerequisites.
- Downloads Go (default `GO_VERSION=1.23.2`) and configures `GOPATH=/go`.
- Installs `golangci-lint` for static analysis.
- Creates a non-root `dev` user (UID/GID configurable via build args).

Supporting files:

- `.dockerignore` – keeps build context small.
- `compose.dev.yml` – optional docker compose wrapper that mounts the repo to `/work/homekit-cli`.
- `README.dev.md` – usage guide describing `make docker-build`, `make docker-shell`, and `make docker-push`.

This container is intended solely for CLI development; service-specific runtime images should live outside `homekit-cli`.
