SHELL := /bin/bash
MAKEFLAGS += --warn-undefined-variables

GO          ?= go
BIN         ?= homekit
PKG         ?= github.com/homekit/homekit-cli
DIST_DIR    ?= dist
PROFILE     ?= default
LOG_LEVEL   ?= info
LOG_FORMAT  ?= console
REGISTRY    ?= garethcxy
IMAGE       ?= $(REGISTRY)/homekit-cli-dev
IMAGE_TAG   ?= latest
DOCKER      ?= docker
DATE        ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
COMMIT      ?= $(shell git rev-parse --short=12 HEAD 2>/dev/null || echo "unknown")
VERSION     ?= $(shell git describe --tags --dirty --always 2>/dev/null || echo "0.0.0-dev")
GOFLAGS     ?=
GO_BUILD_TAGS ?=

export GOFLAGS

ROOT_DIR := $(abspath $(dir $(lastword $(MAKEFILE_LIST)))/..)
BIN_DIR  := $(ROOT_DIR)/$(DIST_DIR)

# Utility for human-readable help output
define print-help
	@printf "\033[1mAvailable Targets\033[0m\n"
	@grep -hE '^[a-zA-Z0-9_/\.-]+:.*?##' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?##"} {printf "\033[36m%-24s\033[0m %s\n", $$1, $$2}'
endef

