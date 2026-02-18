COMPOSE_FILES := -f $(CURDIR)/compose/compose-data.yml \
                -f $(CURDIR)/compose/compose-identity.yml \
                -f $(CURDIR)/compose/compose-proxy.yml \
                -f $(CURDIR)/compose/compose-observability.yml \
                -f $(CURDIR)/compose/compose-web-server.yml \
                -f $(CURDIR)/compose/compose-ui.yml

E2E_FILES := $(COMPOSE_FILES) -f $(CURDIR)/compose/compose-e2e.yml

.PHONY: up down proto recreate build test

up:
	podman compose $(COMPOSE_FILES) up -d

down:
	podman compose $(COMPOSE_FILES) down

proto:
	podman compose $(COMPOSE_FILES) exec web-server sh -c "cd /proto && buf generate"

recreate:
	podman compose $(COMPOSE_FILES) up -d --force-recreate $(SERVICES)

build:
	podman compose $(E2E_FILES) build $(SERVICES)

test:
	podman compose $(E2E_FILES) run --rm e2e
