.PHONY: build run test clean lint deps build-all
.PHONY: dev-setup dev-start dev-stop dev-run dev-reset dev-status
.PHONY: localstack-start localstack-stop localstack-setup localstack-logs localstack-reset
.PHONY: help

# Build the application
build:
	go build -o bin/lazycloud cmd/lazycloud/main.go

# Run the application
run:
	go run cmd/lazycloud/main.go

# Run tests
test:
	go test ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Run linter (install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
lint:
	golangci-lint run

# Install dependencies
deps:
	go mod tidy
	go mod download

# Build for multiple platforms
build-all:
	GOOS=linux GOARCH=amd64 go build -o bin/lazycloud-linux-amd64 cmd/lazycloud/main.go
	GOOS=darwin GOARCH=amd64 go build -o bin/lazycloud-darwin-amd64 cmd/lazycloud/main.go
	GOOS=windows GOARCH=amd64 go build -o bin/lazycloud-windows-amd64.exe cmd/lazycloud/main.go

# =================================
# Development Environment Commands
# =================================

# Complete development setup (recommended for first time)
dev-setup:
	@echo "🚀 Setting up complete development environment..."
	@./scripts/dev-setup.sh
	@echo "⏳ Waiting for LocalStack initialization to complete..."
	@sleep 10
	@echo "✅ Development environment ready!"
	@echo "   Run 'make dev-run' to start LazyCloud with test data"

# Start development environment quickly (if already set up)
dev-start:
	@echo "🐳 Starting LocalStack..."
	@docker-compose up -d
	@echo "⏳ Waiting for LocalStack to be ready..."
	@timeout 60 sh -c 'until curl -f -s http://localhost:4566/_localstack/health > /dev/null; do echo "Waiting..."; sleep 2; done'
	@echo "✅ LocalStack is ready!"

# Stop development environment
dev-stop:
	@echo "🛑 Stopping LocalStack..."
	@docker-compose down

# Run LazyCloud in development mode with LocalStack credentials
dev-run: dev-env-check
	@echo "🚀 Starting LazyCloud in development mode..."
	@export LAZYCLOUD_LOCAL=true && \
	export AWS_ACCESS_KEY_ID=test && \
	export AWS_SECRET_ACCESS_KEY=test && \
	export AWS_DEFAULT_REGION=us-east-1 && \
	export LOCALSTACK_ENDPOINT=http://localhost:4566 && \
	go run cmd/lazycloud/main.go

# Reset development environment (clean slate)
dev-reset:
	@echo "🔄 Resetting development environment..."
	@docker-compose down -v
	@docker-compose up -d
	@echo "⏳ Waiting for LocalStack to be ready..."
	@timeout 60 sh -c 'until curl -f -s http://localhost:4566/_localstack/health > /dev/null; do echo "Waiting..."; sleep 2; done'
	@echo "⏳ Waiting for LocalStack initialization to complete..."
	@sleep 10
	@echo "✅ Development environment reset!"

# Check development environment status
dev-status:
	@echo "📊 Development Environment Status:"
	@echo "=================================="
	@docker-compose ps
	@echo ""
	@echo "LocalStack Health:"
	@curl -f -s http://localhost:4566/_localstack/health 2>/dev/null | jq '.' || echo "❌ LocalStack not accessible"
	@echo ""
	@echo "Test Resources:"
	@echo "Lambda functions:"
	@aws --endpoint-url=http://localhost:4566 lambda list-functions --query 'Functions[].FunctionName' --output table 2>/dev/null || echo "❌ Cannot list Lambda functions"
	@echo "S3 buckets:"
	@aws --endpoint-url=http://localhost:4566 s3 ls 2>/dev/null || echo "❌ Cannot list S3 buckets"

# ============================
# LocalStack Specific Commands  
# ============================

# Start just LocalStack
localstack-start:
	@docker-compose up -d localstack

# Stop LocalStack
localstack-stop:
	@docker-compose stop localstack

# Setup test data in LocalStack (run after LocalStack is started)
localstack-setup:
	@echo "🔧 Setting up test resources in LocalStack..."
	@export AWS_ACCESS_KEY_ID=test && \
	export AWS_SECRET_ACCESS_KEY=test && \
	export AWS_DEFAULT_REGION=us-east-1 && \
	./localstack/init/setup.sh

# View LocalStack logs
localstack-logs:
	@docker-compose logs -f localstack

# Reset LocalStack completely  
localstack-reset:
	@docker-compose down -v localstack
	@docker-compose up -d localstack

# Internal helper to check if development environment is running
dev-env-check:
	@if ! curl -f -s http://localhost:4566/_localstack/health > /dev/null 2>&1; then \
		echo "❌ LocalStack is not running. Run 'make dev-start' or 'make dev-setup' first."; \
		exit 1; \
	fi

# Show available commands
help:
	@echo "LazyCloud Development Commands"
	@echo "============================="
	@echo ""
	@echo "🏗️  Build Commands:"
	@echo "  make build      - Build the application"
	@echo "  make build-all  - Build for multiple platforms"
	@echo "  make clean      - Clean build artifacts"
	@echo ""
	@echo "🧪 Testing & Quality:"
	@echo "  make test       - Run tests"
	@echo "  make lint       - Run linter"
	@echo "  make deps       - Install/update dependencies"
	@echo ""
	@echo "🚀 Development Environment:"
	@echo "  make dev-setup  - Complete setup (first time)"
	@echo "  make dev-start  - Start LocalStack quickly"
	@echo "  make dev-run    - Run LazyCloud with LocalStack"
	@echo "  make dev-stop   - Stop development environment"
	@echo "  make dev-reset  - Reset to clean state"
	@echo "  make dev-status - Show environment status"
	@echo ""
	@echo "🐳 LocalStack Commands:"
	@echo "  make localstack-start  - Start LocalStack only"
	@echo "  make localstack-stop   - Stop LocalStack only"
	@echo "  make localstack-setup  - Setup test resources"
	@echo "  make localstack-logs   - View LocalStack logs"
	@echo "  make localstack-reset  - Reset LocalStack data"
	@echo ""
	@echo "🔍 Quick Start:"
	@echo "  1. make dev-setup   (first time setup)"
	@echo "  2. make dev-run     (start LazyCloud)"
	@echo ""
	@echo "💡 Environment Variables for Custom Configuration:"
	@echo "  LOCALSTACK_ENDPOINT=http://localhost:4566"
	@echo "  AWS_ACCESS_KEY_ID=test (for LocalStack)"
	@echo "  AWS_SECRET_ACCESS_KEY=test (for LocalStack)"