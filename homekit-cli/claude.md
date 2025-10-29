# HomeKit CLI — High-Level Prompt (No Code)

Use this as a **prompt spec** for an AI to generate a multi-module Go CLI project and its supporting tooling. It intentionally contains **no code** and focuses on structure, behavior, and deliverables.

## 1) Objective

Create a modern, extensible CLI ("homekit-cli") for home server operations that:

- Is built with a contemporary Go CLI framework.
- Is easy to extend with internal subcommands and external plugins.
- Can execute embedded scripts and render embedded templates.
- Ships as single binaries for major platforms and can be developed anywhere via a portable dev environment image.

## 2) Scope & Non-Goals

- In scope: CLI, embedded assets, plugin mechanism, script runner, templating, dev container image for development only, build/release automation.
- Out of scope: Service/runtime Dockerfiles for apps—those live elsewhere in the repository.

## 3) Repository Structure (Conceptual)

- Root repository: contains multiple resources (scripts, service Dockerfiles, etc.). The CLI lives in a subfolder named **homekit-cli**.
- The **homekit-cli** project is multi-modular and includes:
  - A command surface for common home server tasks.
  - Internal libraries for config, execution, logging, plugins, asset management, and templating.
  - A small set of embedded, auditable scripts and templates to provide out-of-the-box functionality.
  - A dev environment image definition housed under a local "docker" folder within the CLI project (for development convenience only).
  - A Makefile system for orchestrating the lifecycle of the project. It should be multi-module with utils and helper libs rather than a single huge Makefile.

## 4) Multi-Module Design (within `homekit-cli`)

- **Module: core** — owns lifecycle, configuration loading, global flags, logging initialization, versioning, and common error handling.
- **Module: commands** — organizes subcommands by domain (e.g., script, docker, backup, system). Each command is independent and registered with the root command.
- **Module: exec** — process execution utilities: timeouts, environment control, redaction policy, result capture, and dry-run functionality.
- **Module: plugins** — discovers and delegates to external executables that follow a naming convention; handles safety controls and manifest handshakes where available.
- **Module: assets** — manages access to embedded scripts and templates; supports listing, selective extraction, and verification via checksums; honors user overrides from a configuration path.
- **Module: templating** — renders text-based templates using structured data from flags, environment, and files.
- **Module: ui** — optional helpers for confirmations, prompts, and progress spinners; adheres to non-interactive mode rules for CI.
- **Module: pkg/utils** — reusable utilities shared across modules, with clear boundaries and minimal coupling.

**Design constraints**

- Internal modules must be importable without side effects and should expose narrowly-scoped interfaces.
- Cross-module communication happens through simple data structures and context propagation.
- Public utilities in `pkg` must remain stable and documented; internal modules can evolve faster.

## 5) Command Surface (Initial Set)

- **script run** — executes a path or an embedded asset; supports dry-run, timeout, environment variables, working directory, and interpreter selection.
- **assets list | extract | verify** — enumerates embedded assets, supports exporting to a directory, and prints provenance information and checksums.
- **template render** — renders an embedded or local template to a destination with data provided via flags and/or files.
- **docker prune / images update** — quality-of-life helpers that shell out safely; designed to be extendable.
- **sys health** — lightweight checks for disk, load, memory; optional integration points for SMART or service probes.

## 6) Plugin System (Exec-Based)

- Subcommand resolution:
  - On unknown subcommand, search for executables prefixed with a defined CLI prefix in the user path and a user-configured plugins directory.
  - Invoke with pass-through flags and arguments; forward selected environment variables for context.
- Safety model:
  - Restrict plugin directories by default; require an explicit opt-in to permit arbitrary paths.
  - Provide a verification command that can display plugin origin and hashes when available.
- Plugin discoverability:
  - Offer a command to list installed plugins and expose basic introspection (name, location, version if provided).

## 7) Embedded Assets: Scripts & Templates

- Provide a curated set of baseline scripts and templates embedded into the binary; ensure they are auditable and versioned.
- Support user overrides in a config directory; the CLI always prefers the user’s modified copy when present.
- Document a workflow for inspecting, exporting, and updating embedded content.

## 8) Configuration & Profiles

- Default configuration location under a standard user config directory.
- Support profiles (e.g., default, lab, prod) that alter defaults and targets; allow selection via a flag or environment variable.
- Environment variable overrides follow a consistent namespace and are documented in help output.

## 9) Development Environment Image (Conceptual)

- Provide a portable, versioned dev environment image used only for building and hacking on the CLI.
- Include commonly required toolchains; do not include secrets or user-specific customizations.
- Document how to pull the image, mount the repository, and start a development session.

## 10) Makefile System (Multi-Module, No Code)

Design a thin, multi-file Make system that orchestrates the lifecycle across modules without embedding logic.

**Structure**

- A **top-level Makefile** in `homekit-cli` with only phony orchestration targets.
- Per-module include files (for example: a make fragment for core, commands, exec, plugins, assets, templating) that define module-specific targets and variables.
- A make fragment dedicated to the dev container workflow.

**Global targets**

- **init** — bootstraps the environment and validates prerequisites.
- **deps** — fetches tool dependencies used by build, lint, generation, or release flows.
- **fmt** — applies formatting and basic style fixes across all modules.
- **lint** — runs static analysis across all modules.
- **test** — runs unit tests with coverage and reports summary.
- **generate** — refreshes derived artifacts (documentation, completions, embedded indexes, metadata).
- **build** — compiles binaries for the host platform with version metadata; outputs to a dist directory.
- **run** — executes the CLI in place with the selected profile and logging mode; useful for smoke tests.
- **clean** — removes build and temporary outputs.
- **dist** — prepares distributable archives and checksums for supported platforms and architectures.
- **release** — orchestrates the release flow: tagging, building, packaging, and publishing artifacts via CI.
- **docker-build** — builds the dev environment image.
- **docker-push** — pushes the dev environment image to a registry.
- **docker-shell** — launches an interactive shell within the dev image with repository mounts.

**Variables (overridable)**

- VERSION, COMMIT, DATE — version metadata; default to inferred values when not provided.
- BIN, PKG, DIST_DIR — binary name, module import path, and output directory.
- REGISTRY, IMAGE, IMAGE_TAG — container registry coordinates for the dev image.
- PROFILE, LOG_LEVEL, LOG_FORMAT — runtime defaults for development and CI.

**Conventions**

- All targets are phony, deterministic, and safe to re-run.
- Heavy lifting belongs to dedicated tools or scripts; Make only orchestrates.
- CI should call the same targets used locally; no bespoke CI-only steps.

## 11) Quality & Reliability

- Clear, helpful help texts; consistent global flags; readable error messages with actionable guidance.
- Exit codes map cleanly to failure modes (validation error vs. execution error vs. plugin error).
- Logging supports readable human output and structured machine output.
- Provide basic telemetry hooks that are off by default and documented (if telemetry is added later).

## 12) Release Management

- Cross-platform builds for common OS/architectures.
- Checksums for all release artifacts.
- Optional containerized build steps for reproducibility.
- Change log generation tied to tags; release notes summarize notable changes.

## 13) AI Tasking (What to Ask the AI to Do)

- Scaffold the multi-module project: create modules, wiring, and minimal root command with global options.
- Implement module contracts: define how core, exec, plugins, assets, templating, and commands interact.
- Add the command surface: script, assets, template, docker helpers, and system health.
- Implement embedded asset management and user override rules.
- Create the multi-file Make system with the targets and variables described above; ensure it supports both local and CI.
- Produce a developer guide that explains how to work with profiles, plugins, and embedded assets, and how to use the dev container.
- Add tests for exec timeouts, redaction behavior, plugin dispatch, and configuration precedence.
- Provide initial release configuration and documentation for tagging, building, and publishing artifacts.

## 14) Acceptance Criteria

- The CLI builds and runs on major platforms with consistent help and global flags.
- The Make targets operate end-to-end locally and in CI, using the same entry points.
- Plugins resolve correctly, with safe defaults and clear opt-outs for extended behavior.
- Embedded assets are discoverable, exportable, and verifiable; user overrides work as described.
- Distributable archives and checksums are produced for supported platforms.
- The dev environment image is documented and can be used to develop the CLI in a clean environment.

— End of high-level prompt —
