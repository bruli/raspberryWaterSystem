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
	go test -tags infra -race ./internal/infrastructure/disk/... --count=1
docker-up:
	docker-compose up -d --build server

docker-down:
	docker-compose down server

lint:
	golangci-lint run
	go mod tidy -v && git --no-pager diff --quiet go.mod go.sum