COMPOSE_FILES := -f $(CURDIR)/compose/compose-data.yml \
                -f $(CURDIR)/compose/compose-identity.yml \
                -f $(CURDIR)/compose/compose-proxy.yml \
                -f $(CURDIR)/compose/compose-observability.yml \
                -f $(CURDIR)/compose/compose-web-server.yml \
                -f $(CURDIR)/compose/compose-ui.yml

E2E_FILES := $(COMPOSE_FILES) -f $(CURDIR)/compose/compose-e2e.yml

.PHONY: up down clean proto proto-ts recreate build test unit-test tidy go-get ca-init migrate dgraph-schema dgraph-seed

up:
	podman compose $(COMPOSE_FILES) up -d
	@printf "Waiting for Zitadel PAT..."
	@until podman exec compose_zitadel-login_1 cat /zitadel-data/login-client.pat >/dev/null 2>&1; do printf "."; sleep 2; done
	@echo " ready"
	@ROOTSTOCK_IDENTITY_ZITADEL_PAT=$$(podman exec compose_zitadel-login_1 cat /zitadel-data/login-client.pat | tr -d '\r\n') \
		podman compose $(COMPOSE_FILES) up -d web-server

down:
	podman compose $(COMPOSE_FILES) down

clean:
	podman compose $(COMPOSE_FILES) down --volumes

proto:
	podman compose $(COMPOSE_FILES) exec web-server sh -c "cd /proto && buf generate"

proto-ts:
	podman compose $(COMPOSE_FILES) exec ui sh -c "export PATH=/app/node_modules/.bin:\$$PATH && cd /proto && buf generate --template buf.gen.ts.yaml"

recreate:
	podman compose $(COMPOSE_FILES) up -d --force-recreate $(SERVICES)

build:
	podman compose $(E2E_FILES) build $(SERVICES)

unit-test:
	podman compose $(COMPOSE_FILES) exec web-server go test -p 1 -count=1 $(PKGS)

tidy:
	podman compose $(COMPOSE_FILES) exec web-server go mod tidy

go-get:
	podman compose $(COMPOSE_FILES) exec web-server go get $(PKG)

test:
	podman compose $(E2E_FILES) run --rm e2e

migrate:
	podman compose $(COMPOSE_FILES) exec web-server go run ./cmd/migrate

ca-init:
	@mkdir -p certs
	openssl ecparam -genkey -name prime256v1 -noout -out certs/ca.key
	openssl req -new -x509 -key certs/ca.key -out certs/ca.crt -days 3650 \
		-subj "/CN=Rootstock Dev CA"

dgraph-schema:
	podman compose $(COMPOSE_FILES) exec -T dgraph-alpha \
		curl -sf http://localhost:8080/alter -d "$$(cat web-server/repo/graph/schema/rootstock_runtime.graphql)"

dgraph-seed:
	{ echo '{ set {'; grep -v '^\s*#' web-server/repo/graph/seed/campaign_lifecycle.rdf | sed '/^\s*$$/d'; echo '} }'; } > /tmp/dgraph-seed.rdf
	podman cp /tmp/dgraph-seed.rdf compose_dgraph-alpha_1:/tmp/seed.rdf
	podman compose $(COMPOSE_FILES) exec dgraph-alpha \
		curl -sf 'http://localhost:8080/mutate?commitNow=true' -H 'Content-Type: application/rdf' --data-binary @/tmp/seed.rdf
	@rm -f /tmp/dgraph-seed.rdf
