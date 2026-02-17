COMPOSE_FILES := -f $(CURDIR)/compose/compose-data.yml \
                -f $(CURDIR)/compose/compose-identity.yml \
                -f $(CURDIR)/compose/compose-proxy.yml \
                -f $(CURDIR)/compose/compose-web-server.yml

.PHONY: up down proto

up:
	podman compose $(COMPOSE_FILES) up -d

down:
	podman compose $(COMPOSE_FILES) down

proto:
	podman compose $(COMPOSE_FILES) exec web-server sh -c "cd /proto && buf generate"
