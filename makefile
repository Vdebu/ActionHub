.PHONY: run/web build

run/web:
	go run ./cmd/web

build:
	go build -o=./bin/web ./cmd/web