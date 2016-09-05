# Makefile configuration
.DEFAULT_GOAL := help
.PHONY: help fmt vet test deps cyclo lint travis

ok: fmt vet cyclo lint test ## Prepares codebase (fmt+vet+test)

fmt: ## Golang code formatting tool
	@echo "Running formatting tool"
	@gofmt -s -w .

vet: ## Check code against common errors
	@echo "Running code inspection tools"
	@go vet ./...

cyclo: ## Check cyclomatic complexity
	@echo "Running cyclomatic complexity test"
	@${GOPATH}/bin/gocyclo -over 15 .

test: ## Run tests
	@echo "Running unit tests"
	@go test ./...

lint: ## Code linting
	@echo "Running code linting"
	@${GOPATH}/bin/golint ./...

deps: ## Download required dependencies
	go get github.com/stretchr/testify/assert

travis: deps vet test

help:
	@grep --extended-regexp '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-10s\033[0m %s\n", $$1, $$2}'
