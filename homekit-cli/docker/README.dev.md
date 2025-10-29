# Development Container

This image provides a portable environment for building and testing `homekit-cli`.

## Build

```bash
make docker-build IMAGE=ghcr.io/your-user/homekit-cli-dev TAG=latest
```

## Run an Interactive Shell

```bash
make docker-shell IMAGE=ghcr.io/your-user/homekit-cli-dev TAG=latest
```

## Publish

```bash
make docker-push IMAGE=ghcr.io/your-user/homekit-cli-dev TAG=latest
```

The container assumes the repository is mounted at `/work/homekit-cli` when running interactively.
