.EXPORT_ALL_VARIABLES:
TEST_POSTGRES_DRIVER=postgres
TEST_POSTGRES_HOST=127.0.0.1
TEST_POSTGRES_PORT=5431
TEST_POSTGRES_USER=postgres
TEST_POSTGRES_NAME=postgres
TEST_POSTGRES_PASSWORD=postgres


dependency:
	@go get -v ./...


live:
	@gin -b go-gin-starterkit -p 8081 -a 8080 run main.go


unit: dependency
	@go test -race -v -short ./...

unit_ci:
	@go test -race -coverprofile=coverage.txt -covermode=atomic -v -short ./...


integration: dependency docker-up
	@go test -race -v -run Integration ./...
	@make docker-down

integration_ci: dependency docker-up
	@go test -race -coverprofile=coverage.txt -covermode=atomic -v -run Integration ./...
	@make docker-down


docker-up: docker-down
	@docker-compose -f docker-compose.integration.yml up -d

docker-down:
	@docker-compose -f docker-compose.integration.yml down -v


fmt:
	@go fmt ./...
