dependency:
	@go get -v ./...

unit: dependency
	@go test -v -short ./...

integration: dependency docker-up
	@TEST_POSTGRES_DRIVER=postgres \
	 TEST_POSTGRES_HOST=127.0.0.1 \
	 TEST_POSTGRES_PORT=5431 \
	 TEST_POSTGRES_USER=postgres \
	 TEST_POSTGRES_NAME=postgres \
	 TEST_POSTGRES_PASSWORD=postgres \
	 go test -v -run Integration ./...
	@make docker-down

docker-up:
	@docker-compose -f docker-compose.integration.yml up -d

docker-down:
	@docker-compose -f docker-compose.integration.yml down -v

fmt:
	@go fmt ./...
