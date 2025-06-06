define help
Usage: make <command>
Commands:
   help:                      Show this help information
   tool-jsonschema:           Install gojsonschema tool
   test:                      Run unit tests
   test-with-infra:           Run infrastructure tests
   test-integration:          Run integration tests
   test-functional:           Run functional tests
   docker-up:                 Start docker containers
   docker-down:               Stop docker containers
   docker-ps:                 To watch all docker containers
   docker-exec                To entry into water system container
   lint:                      Execute go linter
   clean:                     To clean code
   fumpt:      	               Format code
   import-jsonschema:         Import and generate DTOS from json schemas
   encryptVault:              Encrypt vault secret file
   decryptVault:              Decrypt vault secret file
   build:                     Compile the project
   docker-exec-builder:       Start builder docker container and entry inside it. Build project here.
   deploy:                    Deploy the code to raspberry
endef
export help

.PHONY: help
help:
	@echo "$$help"

.PHONY: docker-logs
docker-logs:
	docker logs -f water_system


.PHONY: tool-jsonschema
tool-jsonschema:
	go get github.com/atombender/go-jsonschema/...
	go install github.com/atombender/go-jsonschema@latest

.PHONY: test
test:
	go test -race ./... -json|go tool tparse -all

.PHONY: test-with-infra
test-with-infra:
	go test -tags infra -race ./internal/infra/disk/... --count=1 -json|go tool tparse --all

.PHONY: test-integration
test-integration:
	go test -tags integration -race ./internal/infra/telegram/... --count=1

.PHONY: test-functional
test-functional:
	go test -tags functional -race ./tests/functional/... --count=1 -json|go tool tparse --all

.PHONY: docker-up
docker-up:
	docker compose up -d --build water_system

.PHONY: docker-down
docker-down:
	docker compose down

.PHONY: docker-ps
docker-ps:
	docker compose ps

.PHONY: docker-exec
docker-exec:
	docker exec -it water_system bash

.PHONY: lint
lint:
	go tool golangci-lint run
	devops/scripts/json-lint.sh
	go mod tidy -v && git --no-pager diff --quiet go.mod go.sum

.PHONY: clean
clean:
	go fmt ./...

.PHONY: fumpt
fumpt:
	go tool gofumpt -w -l .

.PHONY: import-jsonschema
import-jsonschema:
	devops/scripts/import_jsonschema.sh

.PHONY: encryptVault
encryptVault:
	ansible-vault encrypt --vault-id raspberry_water_system@devops/ansible/password devops/ansible/inventories/production/group_vars/raspberry_water_system/vault.yml

.PHONY: decryptVault
decryptVault:
	ansible-vault decrypt --vault-id raspberry_water_system@devops/ansible/password devops/ansible/inventories/production/group_vars/raspberry_water_system/vault.yml

.PHONY: build
build:
	@make clean
	CC=arm-linux-gnueabi-gcc CGO_ENABLED=1 GOOS=linux GOARCH=arm GOARM=6 go build -a -ldflags "-s -w" -tags prod -buildvcs=false -o devops/ansible/assets/server ./cmd/server/

.PHONY: docker-exec-builder
docker-exec-builder:
	docker build -t builder .
	docker run -it --rm -v $(shell pwd):/app builder bash

.PHONY: deploy
deploy:
	devops/scripts/deploy.sh
