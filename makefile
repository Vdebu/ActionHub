.PHONY: run/api build vendor audit

run/api:
	go run ./cmd/api

build:
	go build -o=./bin/api ./cmd/api
vendor:
	@echo Tidying and verifying module dependencies...
	go mod tidy
	go mod verify
audit:
	@echo Formatting code...
	go fmt ./...