.PHONY: lint start stop test test-cov

lint:
	@golangci-lint run --fix

start:
	@docker-compose up --build

stop:
	@docker-compose down

test:
	@go test -short ./...

test-cov:
	@go test -coverprofile=coverage.out ./...
	@go tool cover -func coverage.out

update-deps:
	@go get -u all
	@go mod tidy
