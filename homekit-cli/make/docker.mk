DEV_DOCKERFILE ?= docker/Dockerfile.dev
DEV_COMPOSE    ?= docker/compose.dev.yml


.PHONY: container.build
container.build: ## Build the development container image
	@$(DOCKER) build -f $(DEV_DOCKERFILE) -t $(IMAGE):$(IMAGE_TAG) ..

.PHONY: container.push
container.push: ## Push the development container image to a registry
	@$(DOCKER) push $(IMAGE):$(IMAGE_TAG)

compose.dev.up: ## Start the development container
	@IMAGE=$(IMAGE) TAG=$(IMAGE_TAG) $(DOCKER) compose -f $(DEV_COMPOSE) up -d

compose.dev.down: ## Stop the development container
	@IMAGE=$(IMAGE) TAG=$(IMAGE_TAG) $(DOCKER) compose -f $(DEV_COMPOSE) down