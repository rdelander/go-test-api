.PHONY: help test test-unit test-integration test-all test-coverage build run watch clean docker-build docker-up docker-up-db docker-down migrate-up migrate-down sqlc-generate fmt lint dev

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

test: test-unit ## Run unit tests (alias for test-unit)

test-unit: ## Run unit tests only
	@echo "Running unit tests..."
	@go test -tags=unit -v ./...

test-unit-ci: ## Run unit tests with race detection and coverage for CI
	@echo "Running unit tests with race detection and coverage..."
	@go test -tags=unit -v -race -coverprofile=coverage.out ./...

test-integration: ## Run integration tests only
	@echo "Running integration tests..."
	@./scripts/test-integration.sh

test-integration-ci: ## Run integration tests for CI (no docker setup)
	@echo "Running integration tests..."
	@go test -tags=integration -v ./...

test-all: test-unit test-integration ## Run all tests (unit + integration)

test-coverage: ## Run unit tests with coverage report
	@echo "Running tests with coverage..."
	@go test -tags=unit -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

build: ## Build the server binary
	@echo "Building server..."
	@go build -o bin/server ./cmd/server

run: ## Run the server locally
	@go run cmd/server/main.go

watch: ## Run the server with hot reload using air
	@air

clean: ## Remove build artifacts
	@rm -rf bin/
	@rm -f server
	@rm -f coverage.out coverage.html

docker-build: ## Build Docker image
	@docker build -t go-test-api .

docker-up: ## Start all services with docker-compose
	@docker-compose up -d

docker-up-db: ## Start only database services (postgres)
	@docker-compose up -d postgres

docker-down: ## Stop all services
	@docker-compose down

migrate-up: ## Run database migrations up
	@./migrate.sh up

migrate-down: ## Run database migrations down
	@./migrate.sh down

sqlc-generate: ## Generate sqlc code from SQL queries
	@sqlc generate

fmt: ## Format Go code
	@go fmt ./...

fmt-check: ## Check if code is formatted
	@if [ "$$(gofmt -s -l . | wc -l)" -gt 0 ]; then \
		echo "The following files need formatting:"; \
		gofmt -s -l .; \
		exit 1; \
	fi

vet: ## Run go vet
	@go vet ./...

lint: ## Run linter (requires golangci-lint)
	@golangci-lint run

lint-ci: fmt-check vet lint ## Run all linting checks for CI

dev: docker-up-db migrate-up ## Start development environment (DB + migrations)
	@echo "Development environment ready!"
	@echo "Run 'make watch' to start the server with hot reload"
