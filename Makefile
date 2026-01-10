.PHONY: help build run test clean deps migrate

# Default target
help:
	@echo "Available commands:"
	@echo "  deps     - Download dependencies"
	@echo "  build    - Build the application"
	@echo "  run      - Run the application"
	@echo "  test     - Run tests"
	@echo "  migrate  - Run database migrations"
	@echo "  clean    - Clean build artifacts"

# Download dependencies
deps:
	go mod download
	go mod tidy

# Build the application
build:
	go build -o bin/e-commerce-api cmd/main.go

# Run the application
run:
	go run cmd/main.go

# Run tests
test:
	go test ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Run database migrations (requires psql)
migrate:
	psql $(GOOSE_DBSTRING) -f migrations/001_initial_schema.sql

# Development setup
dev-setup: deps migrate
	@echo "Development setup complete!"
	@echo "Copy .env.example to .env and configure your settings"
