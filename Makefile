GO_PATH := $(shell go env GOPATH 2> /dev/null)
MODULE := $(shell awk '/^module/ {print $$2}' go.mod)
NAMESPACE := $(shell awk -F "/" '/^module/ {print $$(NF-1)}' go.mod)
PROJECT_NAME := $(shell awk -F "/" '/^module/ {print $$(NF)}' go.mod)
PATH := $(GO_PATH)/bin:$(PATH)

help:
	@echo "Makefile targets:"
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' Makefile \
	| sed -n 's/^\(.*\): \(.*\)##\(.*\)/    \1 :: \3/p' \
	| column -t -c 1  -s '::'

full: clean lint test

lint:
	@cd ; go get golang.org/x/lint/golint
	@cd ; go get golang.org/x/tools/cmd/goimports
	go get -d ./...
	gofmt -s -w .
	go vet ./...
	golint -set_exit_status=1 ./...
	goimports -w .

test:
	@mkdir -p var/
	@go test -race -cover -coverprofile  var/coverage.txt ./...
	@go tool cover -func var/coverage.txt | awk '/^total/{print $$1 " " $$3}'

clean:
	git clean -Xdf