.EXPORT_ALL_VARIABLES:
TEST_POSTGRES_DRIVER=postgres
TEST_POSTGRES_HOST=127.0.0.1
TEST_POSTGRES_PORT=5431
TEST_POSTGRES_USER=postgres
TEST_POSTGRES_NAME=postgres
TEST_POSTGRES_PASSWORD=postgres
TEST_REDIS_ADDR=127.0.0.1:6378


dependency:
	@go get -v ./...


live:
	@gin -b go-gin-starterkit -p 8081 -a 8080 run main.go


unit: dependency
	@go test -race -v -short ./...

unit_ci:
	@go test -race -coverprofile=coverage.txt -covermode=atomic -v -short ./...


integration: dependency docker_up
	@go test -race -v -run Integration ./...
	@$(MAKE) docker_down

integration_ci: dependency docker_up
	@go test -race -coverprofile=coverage.txt -covermode=atomic -v -run Integration ./...
	@$(MAKE) docker_down


docker_up: docker_down
	@docker-compose -p integration -f docker-compose.integration.yml up -d

docker_down:
	@docker-compose -p integration -f docker-compose.integration.yml down -v


up-db:
	@docker run --name test-db -d -p 5432:5432 postgres:11.3-alpine
