# Makefile for Go microservices application

# Variables
GOPATH=$(shell go env GOPATH)
GOBIN=$(GOPATH)/bin
BINARY_DIR=bin
SERVICES=hello user auth
GO_FILES=$(shell find . -name '*.go' -not -path './vendor/*')
PROTO_FILES=$(shell find api/proto -name '*.proto')

# Colors for pretty output
BLUE=\033[0;34m
GREEN=\033[0;32m
YELLOW=\033[0;33m
RED=\033[0;31m
NC=\033[0m # No Color

# Default target
.PHONY: help
help: ## Show this help message
	@echo "$(BLUE)Available targets:$(NC)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-20s$(NC) %s\n", $$1, $$2}'

# Installation targets
.PHONY: install
install: install-deps install-tools ## Install all dependencies and tools

.PHONY: install-deps
install-deps: ## Install Go dependencies
	@echo "$(BLUE)Installing Go dependencies...$(NC)"
	go mod download
	go mod verify

.PHONY: install-tools
install-tools: ## Install development tools
	@echo "$(BLUE)Installing development tools...$(NC)"
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Build targets
.PHONY: build
build: generate build-all ## Generate code and build all services

.PHONY: build-all
build-all: $(SERVICES) ## Build all services

.PHONY: $(SERVICES)
$(SERVICES): ## Build individual service (hello, user, auth)
	@echo "$(BLUE)Building $@ service...$(NC)"
	@mkdir -p $(BINARY_DIR)
	go build -o $(BINARY_DIR)/$@ ./cmd/$@/

.PHONY: build-linux
build-linux: generate ## Build all services for Linux
	@echo "$(BLUE)Building all services for Linux...$(NC)"
	@mkdir -p $(BINARY_DIR)
	@for service in $(SERVICES); do \
		echo "Building $$service for Linux..."; \
		GOOS=linux GOARCH=amd64 go build -o $(BINARY_DIR)/$$service-linux ./cmd/$$service/; \
	done

.PHONY: build-docker
build-docker: ## Build Docker images for all services
	@echo "$(BLUE)Building Docker images...$(NC)"
	@for service in $(SERVICES); do \
		echo "Building Docker image for $$service..."; \
		docker build -t $$service:latest -f docker/$$service.Dockerfile .; \
	done

# Code generation targets
.PHONY: generate
generate: proto swagger ## Generate all code (protobuf and swagger)

.PHONY: proto
proto: ## Generate protobuf files
	@echo "$(BLUE)Generating protobuf files...$(NC)"
	@chmod +x scripts/generate-proto.sh
	./scripts/generate-proto.sh

.PHONY: swagger
swagger: ## Generate Swagger documentation
	@echo "$(BLUE)Generating Swagger documentation...$(NC)"
	@mkdir -p docs
	@for service in $(SERVICES); do \
		echo "Generating Swagger for $$service..."; \
		mkdir -p docs/$$service; \
		swag init -g internal/$$service/handler/handler.go -o docs/$$service; \
	done

# Testing targets
.PHONY: test
test: ## Run all tests
	@echo "$(BLUE)Running tests...$(NC)"
	go test -v ./...

# Code quality targets
.PHONY: format
format: ## Format Go code
	@echo "$(BLUE)Formatting Go code...$(NC)"
	go fmt ./...
	goimports -w $(GO_FILES)

.PHONY: lint
lint: ## Lint Go code
	@echo "$(BLUE)Linting Go code...$(NC)"
	golangci-lint run

# Running targets
.PHONY: run
run: ## Run all services
	@echo "$(BLUE)Starting all services...$(NC)"
	@chmod +x scripts/run-all.sh
	./scripts/run-all.sh

.PHONY: run-hello
run-hello: ## Run hello service
	@echo "$(BLUE)Starting hello service...$(NC)"
	go run cmd/hello/main.go

.PHONY: run-user
run-user: ## Run user service
	@echo "$(BLUE)Starting user service...$(NC)"
	go run cmd/user/main.go

.PHONY: run-auth
run-auth: ## Run auth service
	@echo "$(BLUE)Starting auth service...$(NC)"
	go run cmd/auth/main.go