.PHONY: encdec-bench lint start stop test test-cov

encdec-bench:
	cd encdec && go test -bench=. -benchmem

lint:
	@golangci-lint run

start:
	@docker-compose up --build

stop:
	@docker-compose down

test:
	@go test -short ./...

test-cov:
	@go test -coverprofile=coverage.out ./...
	@go tool cover -func coverage.out
