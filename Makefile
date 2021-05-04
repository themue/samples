SHELL=/bin/bash
PROTOCCMD=protoc
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOLINT=golangci-lint

GO111MODULE=on


.PHONY: download
download: ## Download module dependencies
	go mod download


.PHONY: build-protos
build-protos: ## Generate the source code by the protocol buffer definitions
	$(PROTOCCMD) -I ./protos --go_out=./cmd/weather-server/weather --python_out=./cmd/weather-client ./protos/weather.proto


.PHONY: build
build: build-protos ## Build the libraries and binaries
	$(GOBUILD) -v  ./...


.PHONY: lint
lint: ## Run the linter
	$(GOLINT) run ./...


.PHONY: test
test: ## Run all the tests
	echo 'mode: atomic' > coverage.txt && $(GOTEST) -v -race -covermode=atomic -coverprofile=coverage.txt -timeout=30s ./...

.PHONY: ci
ci: lint test ## Run all the tests and code checks


.PHONY: clean
clean: ## Clean
	$(GOCLEAN)


.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "%-30s %s\n", $$1, $$2}'

.DEFAULT_GOAL := build
