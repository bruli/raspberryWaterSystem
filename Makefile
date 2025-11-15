SHELL := /bin/bash

# ⚙️ Configuration
APP             ?= water_system
DOCKER_COMPOSE  := COMPOSE_BAKE=true docker compose

# Default goal
.DEFAULT_GOAL := help

# 📚 Declare all phony targets
.PHONY: docker-logs docker-down docker-exec docker-ps docker-up \
        test test-functional lint clean fmt help \
        encryptVault decryptVault build deploy security

# ────────────────────────────────────────────────────────────────
# 🐳 Docker
# ────────────────────────────────────────────────────────────────
docker-up:
	@set -euo pipefail; \
	echo "🚀 Starting services with Docker Compose..."; \
	$(DOCKER_COMPOSE) up -d

docker-down:
	@set -euo pipefail; \
	echo "🛑 Stopping and removing Docker Compose services..."; \
	$(DOCKER_COMPOSE) down

docker-ps:
	@set -euo pipefail; \
	echo "📋 Active services:"; \
	$(DOCKER_COMPOSE) ps

docker-exec:
	@set -euo pipefail; \
	test -n "${SVC:-}" || { echo "❌ Please specify SVC=<service>"; exit 2; }; \
	echo "🔎 Opening shell inside $$SVC..."; \
	$(DOCKER_COMPOSE) exec $$SVC sh

docker-logs:
	@set -euo pipefail; \
	echo "👀 Showing logs for container $(APP) (CTRL+C to exit)..."; \
	docker logs -f $(APP)

# ────────────────────────────────────────────────────────────────
# 🧹 Code quality: format, lint, tests
# ────────────────────────────────────────────────────────────────
fmt:
	@set -euo pipefail; \
	echo "👉 Formatting code with gofumpt..."; \
	go tool gofumpt -w .

security:
	@set -euo pipefail; \
	echo "👉 Check security"; \
	go tool govulncheck ./...

lint:
	@set -euo pipefail; \
	echo "🔍 Running golangci-lint..."; \
	go tool golangci-lint run ./...

test:
	@set -euo pipefail; \
	echo "🧪 Running unit tests (race, JSON → tparse)..."; \
	go test -race ./... -json -cover | go tool tparse -all

test-functional:
	@set -euo pipefail; \
	echo "🧪 Running functional tests..."; \
	# Example: adjust to your own functional test suite
	go test -tags=functional ./... -v

clean:
	@set -euo pipefail; \
	echo "🧹 Cleaning local artifacts..."; \
	rm -rf bin dist coverage .*cache || true; \
	go clean -testcache

# ────────────────────────────────────────────────────────────────
# 🔐 Ansible Vault
# ────────────────────────────────────────────────────────────────
encryptVault:
	@set -euo pipefail; \
	echo "🔐 Encrypting Ansible vault files..."; \
	ansible-vault encrypt --vault-id raspberry_water_system@devops/ansible/password devops/ansible/inventories/production/group_vars/raspberry_water_system/vault.yml

decryptVault:
	@set -euo pipefail; \
	echo "🔓 Decrypting Ansible vault files..."; \
	ansible-vault decrypt --vault-id raspberry_water_system@devops/ansible/password devops/ansible/inventories/production/group_vars/raspberry_water_system/vault.yml

# ────────────────────────────────────────────────────────────────
# 🏗️ Build & Deploy
# ────────────────────────────────────────────────────────────────
build: clean
	@set -euo pipefail; \
	echo "🏗️ Building ARM64 binary for Raspberry Pi..."; \
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 \
	go build -a -ldflags "-s -w" -tags prod -buildvcs=false \
	-o devops/ansible/assets/server ./cmd/server/

deploy: build decryptVault
	@set -euo pipefail; \
	echo "🚚 Deploying with Ansible (production inventory)..."; \
	ansible-playbook -i devops/ansible/inventories/production/hosts devops/ansible/deploy.yml; \
	$(MAKE) encryptVault

# ────────────────────────────────────────────────────────────────
# ℹ️ Help
# ────────────────────────────────────────────────────────────────
help:
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:' Makefile | awk -F':' '{print "  - " $$1}'
