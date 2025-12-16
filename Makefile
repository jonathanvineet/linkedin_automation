.PHONY: help build run dev clean test install

# Default target
help:
	@echo "LinkedIn Automation - Build Commands"
	@echo ""
	@echo "Usage:"
	@echo "  make install    - Install all dependencies"
	@echo "  make build      - Build backend and frontend"
	@echo "  make run        - Run the application"
	@echo "  make dev        - Run in development mode"
	@echo "  make clean      - Clean build artifacts"
	@echo "  make test       - Run tests"
	@echo ""

# Install dependencies
install:
	@echo "📦 Installing Go dependencies..."
	go mod download
	@echo "📦 Installing Node dependencies..."
	npm install
	@echo "✅ Dependencies installed"

# Build backend and frontend
build:
	@echo "🔨 Building Go backend..."
	@mkdir -p bin
	CGO_ENABLED=1 go build -o bin/automation ./cmd/app
	@echo "🔨 Building React frontend..."
	npm run build
	@echo "✅ Build complete"

# Run the application
run: build
	@echo "🚀 Starting automation..."
	./bin/automation

# Development mode
dev:
	@echo "🔧 Starting in development mode..."
	@./start.sh

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	rm -rf bin/
	rm -rf dist/
	rm -rf node_modules/
	rm -rf data/
	rm -rf logs/
	rm -f *.log
	@echo "✅ Clean complete"

# Run tests
test:
	@echo "🧪 Running Go tests..."
	go test -v ./...
	@echo "✅ Tests complete"

# Format code
fmt:
	@echo "✨ Formatting code..."
	go fmt ./...
	@echo "✅ Format complete"

# Lint code
lint:
	@echo "🔍 Linting code..."
	golangci-lint run
	@echo "✅ Lint complete"
