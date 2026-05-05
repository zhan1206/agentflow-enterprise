.PHONY: all build test clean docker docs help

# Variables
GO := go
GOFLAGS := -v
DOCKER := docker
DOCKER_COMPOSE := docker-compose

# Directories
BACKEND_DIR := backend
FRONTEND_DIR := frontend
DOCS_DIR := docs

# Binary
BINARY := agentflow

# Version
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "v0.1.0-dev")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"

all: clean build

## build: Build all components
build: build-backend build-frontend

## build-backend: Build backend binary
build-backend:
	cd $(BACKEND_DIR) && $(GO) build $(LDFLAGS) -o $(BINARY) ./cmd/server

## build-frontend: Build frontend
build-frontend:
	cd $(FRONTEND_DIR) && npm install && npm run build

## test: Run all tests
test: test-backend test-frontend

## test-backend: Run backend tests
test-backend:
	cd $(BACKEND_DIR) && $(GO) test -v -race -coverprofile=coverage.out ./...

## test-frontend: Run frontend tests
test-frontend:
	cd $(FRONTEND_DIR) && npm test

## lint: Run linters
lint:
	cd $(BACKEND_DIR) && $(GO) vet ./...
	cd $(FRONTEND_DIR) && npm run lint

## docker-build: Build Docker images
docker-build:
	$(DOCKER) build -t agentflow-enterprise/api:$(VERSION) -f $(BACKEND_DIR)/Dockerfile $(BACKEND_DIR)
	$(DOCKER) build -t agentflow-enterprise/frontend:$(VERSION) -f $(FRONTEND_DIR)/Dockerfile $(FRONTEND_DIR)

## docker-up: Start services with Docker Compose
docker-up:
	$(DOCKER_COMPOSE) up -d

## docker-down: Stop Docker Compose services
docker-down:
	$(DOCKER_COMPOSE) down

## docker-logs: View Docker Compose logs
docker-logs:
	$(DOCKER_COMPOSE) logs -f

## dev: Start development environment
dev:
	$(DOCKER_COMPOSE) up -d kafka redis elasticsearch jaeger
	cd $(BACKEND_DIR) && $(GO) run ./cmd/server --config ../deploy/configs/development.yaml

## clean: Clean build artifacts
clean:
	rm -rf $(BACKEND_DIR)/$(BINARY)
	rm -rf $(FRONTEND_DIR)/dist
	rm -rf $(FRONTEND_DIR)/node_modules
	$(GO) clean -cache

## docs: Generate documentation
docs:
	cd $(DOCS_DIR) && make html

## fmt: Format code
fmt:
	cd $(BACKEND_DIR) && $(GO) fmt ./...
	cd $(FRONTEND_DIR) && npm run format

## deps: Install dependencies
deps:
	cd $(BACKEND_DIR) && $(GO) mod download
	cd $(FRONTEND_DIR) && npm install

## help: Show this help message
help:
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@sed -n 's/^## //p' $(MAKEFILE_LIST) | column -t -s ':'
