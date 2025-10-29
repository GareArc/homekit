.DEFAULT_GOAL := help

MAKE_FRAGMENTS := $(wildcard make/*.mk)
include $(MAKE_FRAGMENTS)

.PHONY: help
help: ## Show available targets
	@$(call print-help)

.PHONY: init deps fmt lint test generate build run clean dist release docker-build docker-push docker-shell
init: bootstrap.init ## Validate prerequisites and prepare environment

deps: go.deps ## Install Go tooling dependencies

tidy: go.tidy ## Tidy Go dependencies

fmt: go.fmt ## Format Go sources

lint: go.lint ## Run static analysis

test: go.test ## Run unit test suite

generate: go.generate ## Execute code generation hooks

build: go.build ## Build the CLI binary

run: go.run ## Run the CLI locally

clean: go.clean ## Remove build artifacts

dist: dist.package ## Build cross-platform distributables

release: dist.release ## Trigger release workflow (CI)

docker-build: container.build ## Build the dev container image

docker-push: container.push ## Push the dev container image

up: compose.dev.up ## Start the development container

down: compose.dev.down ## Stop the development container
