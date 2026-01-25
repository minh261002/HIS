.PHONY: help run build test clean docker-up docker-down migrate-up migrate-down

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

run: ## Run the application with hot reload
	@echo "Starting application with hot reload..."
	@air

build: ## Build the application
	@echo "Building application..."
	@go build -o bin/api cmd/api/main.go

test: ## Run tests
	@echo "Running tests..."
	@go test -v -cover ./...

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf bin/ tmp/

docker-up: ## Start Docker services
	@echo "Starting Docker services..."
	@cd docker && docker-compose up -d

docker-down: ## Stop Docker services
	@echo "Stopping Docker services..."
	@cd docker && docker-compose down

migrate-up: ## Run database migrations up
	@echo "Running migrations..."
	@migrate -path migrations -database "mysql://his_user:his_password@tcp(localhost:3306)/hospital_db" up

migrate-down: ## Rollback database migrations
	@echo "Rolling back migrations..."
	@migrate -path migrations -database "mysql://his_user:his_password@tcp(localhost:3306)/hospital_db" down

migrate-create: ## Create a new migration (usage: make migrate-create name=migration_name)
	@echo "Creating migration: $(name)"
	@migrate create -ext sql -dir migrations -seq $(name)

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download

tidy: ## Tidy dependencies
	@echo "Tidying dependencies..."
	@go mod tidy

lint: ## Run linter
	@echo "Running linter..."
	@golangci-lint run

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...
