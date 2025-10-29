.PHONY: bootstrap.init
bootstrap.init: ## Validate prerequisites and prepare local environment
	@command -v $(GO) >/dev/null || { echo "[error] Go toolchain not found" >&2; exit 1; }
	@command -v git >/dev/null || { echo "[error] git not found" >&2; exit 1; }
	@mkdir -p $(BIN_DIR)
	@echo "Go version: $$($(GO) version)"
	@echo "Workspace prepared at $(BIN_DIR)"

