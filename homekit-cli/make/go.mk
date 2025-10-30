GO_TOOLS ?= github.com/golangci/golangci-lint/cmd/golangci-lint@latest

LD_FLAGS := -X github.com/homekit/homekit-cli/cmd/homekit.version=$(VERSION) \
	-X github.com/homekit/homekit-cli/cmd/homekit.commit=$(COMMIT) \
	-X github.com/homekit/homekit-cli/cmd/homekit.date=$(DATE) \
	-X github.com/homekit/homekit-cli/cmd/homekit.source=$(PKG)

.PHONY: go.deps
go.deps: ## Install developer Go tools listed in GO_TOOLS
	@for tool in $(GO_TOOLS); do \
		if [ -n "$$tool" ]; then \
			echo "Installing $$tool"; \
			$(GO) install $$tool; \
		fi; \
	done

go.tidy: ## Tidy Go dependencies
	@$(GO) mod tidy

.PHONY: go.fmt
go.fmt: ## Format Go sources
	@$(GO) fmt ./...
	@gofmt -w $(shell find . -name '*.go' -not -path './vendor/*')

.PHONY: go.lint
go.lint: ## Run static analysis
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "[warn] golangci-lint not installed; run 'make deps'"; \
	fi

.PHONY: go.test
go.test: ## Execute unit tests with coverage summary
	@mkdir -p $(BIN_DIR)
	@$(GO) test ./... -coverprofile=$(BIN_DIR)/coverage.out
	@$(GO) tool cover -func=$(BIN_DIR)/coverage.out | tail -n 1

.PHONY: go.generate
go.generate: ## Run code generation hooks
	@$(GO) generate ./...

.PHONY: go.build
go.build: ## Build the CLI binary for the host platform
	@mkdir -p $(BIN_DIR)
	@$(GO) build -ldflags "$(LD_FLAGS)" -o $(BIN_DIR)/$(BIN) .
	@echo "Binary available at $(BIN_DIR)/$(BIN)"

.PHONY: go.run
go.run: ## Run the CLI with the configured profile
	@$(GO) run -ldflags "$(LD_FLAGS)" . --profile $(PROFILE) --log-level $(LOG_LEVEL) --log-format $(LOG_FORMAT)

.PHONY: go.clean
go.clean: ## Remove build artifacts
	@$(GO) clean ./...
	@rm -rf $(DIST_DIR)
