.PHONY: build build-docker test lint

build:
	@go build -o dockerleaks .

build-docker:
	@docker build -t dockerleaks:dev .


test:
	@go test -v ./...

lint:
	@test -z $(gofmt -l .)
	@golangci-lint run