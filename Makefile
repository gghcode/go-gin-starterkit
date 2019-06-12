dependency:
	@go get -v ./...

unit: dependency
	@go test -v -short ./...

integration: dependency docker-up
	@TEST_POSTGRES_DRIVER=postgres \
	 TEST_POSTGRES_HOST=localhost \
	 TEST_POSTGRES_PORT=5431 \
	 TEST_POSTGRES_USER=postgres \
	 TEST_POSTGRES_NAME=postgres \
	 TEST_POSTGRES_PASSWORD=postgres \
	 go test -v -run Integration ./...
	@make docker-down

docker-up: docker-down
	@docker-compose -f test-compose.yml up -d

docker-down:
	@docker-compose -f test-compose.yml down

fmt:
	@go fmt ./...
