# Release Workflow Overview

Use this document to expand the release automation once the project leaves the prototype stage.

1. Tag the repository (e.g. `git tag v0.1.0`).
2. Run `make dist VERSION=v0.1.0` to build archives and checksums locally or in CI.
3. Attach artifacts to a GitHub release or let CI publish them automatically.
4. Push the dev container image with `make docker-push IMAGE=ghcr.io/your-user/homekit-cli-dev TAG=v0.1.0`.
5. Update changelog and documentation accordingly.

Integrating with Goreleaser or other release tooling can reuse the `dist` artifacts when you are ready.
