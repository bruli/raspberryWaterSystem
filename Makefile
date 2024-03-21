define help
Usage: make <command>
Commands:
   help:                      Show this help information
   tools-local:               Install all tools
   tool-golangci-lint:        Install golangci linter
   tool-fumpt:                Install gofumpt tool
   tool-moq:                  Install moq tool
   tool-jsonschema:           Install gojsonschema tool
   tool-json-lint:            Install json lint linter
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

.PHONY: tools-local
tools-local: tool-golangci-lint tool-moq tool-fumpt	 tool-jsonschema tool-json-lint

.PHONY: tool-golangci-lint
tool-golangci-lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

.PHONY: tool-fumpt
tool-fumpt:
	go install mvdan.cc/gofumpt@latest

.PHONY: tool-moq
tool-moq:
	go get github.com/matryer/moq@latest

.PHONY: tool-jsonschema
tool-jsonschema:
	go get github.com/atombender/go-jsonschema/...
	go install github.com/atombender/go-jsonschema@latest

.PHONY: tool-json-lint
tool-json-lint:
	go get github.com/santhosh-tekuri/jsonschema/cmd/jv

.PHONY: test
test:
	go test -race ./...

.PHONY: test-with-infra
test-with-infra:
	go test -tags infra -race ./internal/infra/disk/... --count=1

.PHONY: test-integration
test-integration:
	go test -tags integration -race ./internal/infra/telegram/... --count=1

.PHONY: test-functional
test-functional:
	go test -tags functional -race ./tests/functional/... --count=1

.PHONY: docker-up
docker-up:
	docker-compose up -d --build water_system

.PHONY: docker-down
docker-down:
	docker-compose down

.PHONY: docker-ps
docker-ps:
	docker-compose ps

.PHONY: docker-exec
docker-exec:
	docker exec -it water_system bash

.PHONY: lint
lint:
	golangci-lint run
	devops/scripts/json-lint.sh
	go mod tidy -v && git --no-pager diff --quiet go.mod go.sum

.PHONY: clean
clean:
	go fmt ./...

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
