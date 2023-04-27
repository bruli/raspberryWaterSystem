docker-logs:
	docker logs -f water_system

tools-local: tool-golangci-lint tool-moq tool-fumpt	 tool-jsonschema tool-json-lint

tool-golangci-lint:
	devops/scripts/goget.sh github.com/golangci/golangci-lint/cmd/golangci-lint

tool-fumpt:
	devops/scripts/goget.sh mvdan.cc/gofumpt

tool-moq:
	devops/scripts/goget.sh github.com/matryer/moq

tool-jsonschema:
	devops/scripts/goget.sh github.com/atombender/go-jsonschema/...
	devops/scripts/goget.sh github.com/atombender/go-jsonschema/cmd/gojsonschema

tool-json-lint:
	go get github.com/santhosh-tekuri/jsonschema/cmd/jv

test:
	go test -race ./...

test-with-infra:
	go test -tags infra -race ./internal/infra/disk/... --count=1

test-integration:
	go test -tags integration -race ./internal/infra/telegram/... --count=1

test-functional:
	go test -tags functional -race ./functional_test/... --count=1

docker-up:
	docker-compose up -d --build water_system

docker-down:
	docker-compose down

docker-ps:
	docker-compose ps

docker-exec:
	docker exec -it water_system bash:q

lint:
	golangci-lint run
	devops/scripts/json-lint.sh
	go mod tidy -v && git --no-pager diff --quiet go.mod go.sum

clean:
	go fmt ./...

import-jsonschema:
	devops/scripts/import_jsonschema.sh

encryptVault:
	ansible-vault encrypt --vault-id raspberry_water_system@devops/ansible/password devops/ansible/inventories/production/group_vars/raspberry_water_system/vault.yml
decryptVault:
	ansible-vault decrypt --vault-id raspberry_water_system@devops/ansible/password devops/ansible/inventories/production/group_vars/raspberry_water_system/vault.yml

build:
	@make clean
	CC=arm-linux-gnueabi-gcc CGO_ENABLED=1 GOOS=linux GOARCH=arm GOARM=6 go build -a -ldflags "-s -w" -tags prod -buildvcs=false -o devops/ansible/assets/server ./cmd/server/

docker-exec-builder:
	docker build -t builder .
	docker run -it --rm -v $(shell pwd):/app builder bash

deploy:
	devops/scripts/deploy.sh
