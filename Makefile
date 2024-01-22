.PHONY: help

help: ## Prints help for targets with comments
	@grep -E '^[a-zA-Z0-9.\ _-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

update: ## Installs all dependencies
	@go mod download

build: ## Build this software
	@go build -o backup_to_remote_storage cmd/main.go || echo "go build failed"

unittests: ## Executes all unit tests
	@go test ./cmd/...
