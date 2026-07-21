SHELL = /bin/sh

CONTAINER_RUNTIME ?= docker

.PHONY: check-container-runtime up down down-volumes run-tests test lint proto unit

up: check-container-runtime down ## Up containers with compose
	$(CONTAINER_RUNTIME) compose up --build -d

down: check-container-runtime ## Stop and remove containers
	$(CONTAINER_RUNTIME) compose down

down-volumes: check-container-runtime ## Stop and remove containers and all volumes
	$(CONTAINER_RUNTIME) compose down -v

run-tests: ## Run tests container
	$(CONTAINER_RUNTIME) run --rm --network=host tests:latest

test: check-container-runtime ## Up containers and run tests
	@$(MAKE) down-volumes
	@$(MAKE) up
	@echo "Waiting for cluster to start" \
		&& for i in $$(seq 15); do \
			curl -sf http://localhost:28080 >/dev/null 2>&1 && break \
				|| true; \
			sleep 1; \
		done || { echo "Error: timeout"; exit 1; }
	@$(MAKE) run-tests; status=$$?; $(MAKE) down-volumes; exit $$status

lint: ## Run linters
	$(MAKE) -C repo-stat lint

proto: ## Compile protobuf files
	$(MAKE) -C repo-stat protobuf

unit: ## Run tests
	$(MAKE) -C repo-stat test

