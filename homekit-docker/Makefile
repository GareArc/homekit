.DEFAULT_GOAL := help

MAKE_FRAGMENTS := $(wildcard make/*.mk)
include $(MAKE_FRAGMENTS)

.PHONY: help dev.docker.build dev.docker.release
help: ## Show available targets
	@$(call print-help)

dev.docker.build: docker.build ## Build the development container image

dev.docker.release: 
	# for each file in the dev directory, build and push the image with tag based on the filename suffix
	for file in dev/*.Dockerfile; do \
		full_name=$(DEV_IMAGE):$$(basename $$file .Dockerfile) && \
		$(DOCKER) build -f $$file -t $$full_name . && \
		$(DOCKER) push $$full_name; \
	done
