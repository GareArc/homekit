.PHONY: docker.build
docker.build:
	@$(DOCKER) build -f $(DOCKERFILE) -t $(DEV_IMAGE):$(TAG) .