GO_PATH := $(shell go env GOPATH 2> /dev/null)
MODULE := $(shell awk '/^module/ {print $$2}' go.mod)
NAMESPACE := $(shell awk -F "/" '/^module/ {print $$(NF-1)}' go.mod)
PROJECT_NAME := $(shell awk -F "/" '/^module/ {print $$(NF)}' go.mod)
PATH := $(GO_PATH)/bin:$(PATH)

help: ## Show this help menu
	@echo "Makefile targets:"
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' Makefile \
	| sed -n 's/^\(.*\): \(.*\)##\(.*\)/    \1 :: \3/p' \
	| column -t -c 1  -s '::'

full: clean lint test ## Clean project and run all checks

lint: ## Run lint checks
	@cd ; go get golang.org/x/lint/golint
	@cd ; go get golang.org/x/tools/cmd/goimports
	go get -d ./...
	go mod tidy
	gofmt -s -w .
	go vet ./...
	golint -set_exit_status=1 ./...
	goimports -w .

test: ## Run unit tests
	@mkdir -p var/
	@go test -race -cover -coverprofile  var/coverage.txt ./...
	@go tool cover -func var/coverage.txt | awk '/^total/{print $$1 " " $$3}'

docs: ## Start a godoc server
	@cd ; go get golang.org/x/tools/cmd/godoc
	@echo "Docs here: http://localhost:3232/pkg/${MODULE}"
	@godoc -http=localhost:3232 -index -index_interval 2s -play

clean: ## Remove all files listed in .gitignore
	git clean -Xdf

post-lint:
	@git diff --exit-code --quiet || (echo "There should not be any changes after the lint runs" && git status && exit 1;)

pipeline: full post-lint
