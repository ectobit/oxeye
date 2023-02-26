.PHONY: lint start stop test test-cov

lint:
	@golangci-lint run \
		--enable-all \
		--disable deadcode \
		--disable exhaustivestruct \
		--disable golint \
		--disable ifshort \
		--disable interfacer \
		--disable maligned \
		--disable nosnakecase \
		--disable scopelint \
		--disable structcheck \
		--disable varcheck

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
