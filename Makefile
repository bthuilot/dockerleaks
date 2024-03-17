.PHONY: build build-docker test lint

build:
	@go build -o dockerleaks .

build-docker:
	@docker build -t dockerleaks:dev .


test:
	@go test -v ./...

test-coverage:
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

lint:
	@test -z $(gofmt -l .)
	@golangci-lint run