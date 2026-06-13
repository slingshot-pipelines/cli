SHELL := bash -e

MAKEFILE_DIR := $(dir $(abspath $(firstword $(MAKEFILE_LIST))))

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: lint
lint: tidy
	go vet ./...
	go fmt ./...

.PHONY: unit-test
unit-test:
	go test ./...
