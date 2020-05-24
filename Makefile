args = `arg="$(filter-out $@,$(MAKECMDGOALS))" && echo $${arg:-${1}}`
prepare_docker:
	@go build -o ./test/acceptance/server -i ./cmd/server/main.go
	@bash -c "cd ./test/acceptance && docker-compose up -d --build"

finish_docker:
	@bash -c "cd ./test/acceptance && docker-compose stop"

acceptance:
	@make finish_docker
	@make prepare_docker
	sleep 5
	@make migration_migrate
	@bash -c "GOTEST_PALETTE="red,blue" gotest ./test/acceptance -v"
	@make finish_docker

unit:
	@bash -c "cd internal && GOTEST_PALETTE="red,blue" gotest ./..."

coverage:
	@bash -c "cd internal && GOTEST_PALETTE="red,blue" gotest -coverprofile=coverage.out ./..."
	@bash -c "cd internal && go tool cover -html=coverage.out"

build:
	@cd cmd/server && CC=arm-linux-gnueabi-gcc CGO_ENABLED=1 GOOS=linux GOARCH=arm GOARM=6 go build

deploy:
	@make build
	@bash -c "cd cmd/server && scp server waterSystem:"

all_tests:
	@echo "executing unit tests..."
	@make unit
	@echo "\n"
	@echo "executing acceptance tests..."
	@make acceptance

migration_migrate:
	@migrate -database "mysql://raspberry:raspberry@tcp(localhost:3306)/raspberryWaterSystem" -path ./internal/infrastructure/migrations up

migration_create:
	@migrate create -ext sql -dir ./internal/infrastructure/migrations -seq $(call args,new_migration)