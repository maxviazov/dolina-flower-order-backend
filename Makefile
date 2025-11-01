# Dolina Flower Order Backend Makefile

.PHONY: help build run test clean deps fmt lint vet mod-tidy docker-build docker-run

# Variables
APP_NAME=dolina-flower-backend
BIN_DIR=bin
MAIN_FILE=cmd/server/main.go
DOCKER_IMAGE=$(APP_NAME):latest

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Development
build: ## Build the application
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BIN_DIR)
	@go build -o $(BIN_DIR)/$(APP_NAME) $(MAIN_FILE)
	@echo "Build complete: $(BIN_DIR)/$(APP_NAME)"

run: ## Run the application
	@echo "Starting $(APP_NAME)..."
	@go run $(MAIN_FILE)

dev: ## Run in development mode
	@echo "Starting development server..."
	@go mod tidy && go run $(MAIN_FILE)

# Testing
test: ## Run tests
	@echo "Running tests..."
	@go test -v ./...

test-cover: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test -v -cover ./...

bench: ## Run benchmarks
	@echo "Running benchmarks..."
	@go test -bench=. ./...

# Code quality
fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...

vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...

# Dependencies
deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download

mod-tidy: ## Tidy go modules
	@echo "Tidying go modules..."
	@go mod tidy

# Cleanup
clean: ## Clean build artifacts
	@echo "Cleaning up..."
	@rm -rf $(BIN_DIR)
	@go clean

# AWS Lambda
lambda-build: ## Build for AWS Lambda
	@echo "Building for AWS Lambda..."
	@mkdir -p $(BIN_DIR)
	@GOOS=linux GOARCH=amd64 go build -o $(BIN_DIR)/bootstrap $(MAIN_FILE)
	@zip -j $(BIN_DIR)/lambda.zip $(BIN_DIR)/bootstrap
	@cp $(BIN_DIR)/lambda.zip terraform/
	@echo "Lambda deployment package created"

# Deploy to AWS
deploy: lambda-build ## Deploy to AWS
	@echo "Deploying to AWS..."
	@cd terraform && terraform init && terraform apply -auto-approve

# Destroy AWS resources
destroy: ## Destroy AWS resources
	@echo "Destroying AWS resources..."
	@cd terraform && terraform destroy -auto-approve
	@rm -f terraform/lambda.zip

# Setup
setup: deps ## Setup development environment
	@echo "Development environment setup complete!"
	@echo "Copy .env.example to .env and configure your settings"

check: fmt vet test ## Run all checks
	@echo "All checks passed!"