SHELL = /bin/sh

CONTAINER_RUNTIME ?= docker

.PHONY: check-container-runtime up down down-volumes run-tests test lint proto unit

up: check-container-runtime down ## Up containers with compose
	$(CONTAINER_RUNTIME) compose up --build -d

down: check-container-runtime ## Stop and remove containers
	$(CONTAINER_RUNTIME) compose down

down-volumes: check-container-runtime ## Stop and remove containers and all volumes
	$(CONTAINER_RUNTIME) compose down -v


lint: ## Run linters
	$(MAKE) -C repo-stat lint

proto: ## Compile protobuf files
	$(MAKE) -C repo-stat protobuf

unit: ## Run tests and generate coverage report
	$(MAKE) -C repo-stat test
	mv repo-stat/cover.html .

check-container-runtime: ## Check container runtime is available
ifeq (0,$(MAKELEVEL))
	@$(if $(strip $(CONTAINER_RUNTIME)),\
		$(info Using $(CONTAINER_RUNTIME) as container runtime),\
		$(error No container runtime found. Install Podman or Docker))
endif

help: ## Display this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
