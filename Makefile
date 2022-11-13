.PHONY: lint start stop test test-cov

lint:
	@golangci-lint run --exclude-use-default=false --enable-all \
		--disable interfacer \
		--disable scopelint \
		--disable ifshort \
		--disable maligned \
		--disable nosnakecase \
		--disable golint \
	    --disable exhaustivestruct \
		--disable deadcode \
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
