# Workspace {{ .name }}

## How to use
Then start your development container with:

```bash
make up
```
`code` is the directory that will be mounted to the container. Put your codebase there using gh cli.
Notice that the container image does not contain docker in docker, so you may need to run some scripts from the host machine.

## How to stop

```bash
make down
```