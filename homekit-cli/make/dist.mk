PLATFORMS ?= linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64

.PHONY: dist.package
dist.package: ## Build cross-platform distributable binaries
	@mkdir -p $(DIST_DIR)
	@checksum() { if command -v shasum >/dev/null 2>&1; then shasum -a 256 "$$1"; else sha256sum "$$1"; fi; }; \
	for platform in $(PLATFORMS); do \
		IFS=/ read -r goos goarch <<<"$$platform"; \
		binary_name="$(BIN)-$$goos-$$goarch"; \
		if [ "$$goos" = "windows" ]; then binary_name="$$binary_name.exe"; fi; \
		echo "Building $$platform"; \
		GOOS=$$goos GOARCH=$$goarch $(GO) build -ldflags "$(LD_FLAGS)" -o $(DIST_DIR)/$$binary_name ./cmd/homekit; \
		checksum $(DIST_DIR)/$$binary_name > $(DIST_DIR)/$$binary_name.sha256; \
	done

.PHONY: dist.release
dist.release: ## Placeholder for CI-driven release automation
	@echo "Trigger release workflow via CI (see docs/releases.md)"
