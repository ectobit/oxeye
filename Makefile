.PHONY: lint start stop test test-cov

lint:
	@golangci-lint run --exclude-use-default=false --enable-all \
		--disable golint \
		--disable interfacer \
		--disable maligned \
		--disable scopelint \
	    --disable exhaustivestruct

start:
	@docker-compose up --build

stop:
	@docker-compose down

test:
	@go test -short ./...

test-cov:
	@go test -coverprofile=coverage.out ./...
	@go tool cover -func coverage.out
