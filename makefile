.PHONY: run/web build vendor audit

run/web:
	go run ./cmd/web

build:
	go build -o=./bin/web ./cmd/web
vendor:
	@echo Tidying and verifying module dependencies...
	go mod tidy
	go mod verify
audit:
	@echo Formatting code...
	go fmt ./...