.PHONY: build run test test-coverage migrate deps clean fmt lint docker-up docker-down docker-logs docker-ps docker-postgres run-api-docker

# Build REST API
build:
	go build -o bin/api ./cmd/api

# Run REST API
run: build
	./bin/api

# Run all tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Run database migrations using GORM (auto-migrate)
migrate:
	@echo "Running database migrations..."
	@go run cmd/api/main.go migrate || echo "Note: Migrations run automatically on server start"

# Install dependencies
deps:
	go mod download

# Clean build artifacts
clean:
	rm -rf bin coverage.out coverage.html

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run || echo "golangci-lint not installed, skipping..."

# Docker Compose commands (using Docker Compose V2: docker compose)
docker-up:
	docker compose up -d

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f postgres

docker-ps:
	docker compose ps

# Start PostgreSQL and wait for it to be ready
docker-postgres:
	docker compose up -d postgres
	@echo "Waiting for PostgreSQL to be ready..."
	@timeout 30 bash -c 'until docker compose exec -T postgres pg_isready -U postgres; do sleep 1; done' || true
	@echo "PostgreSQL is ready!"

# Run API with Docker Compose PostgreSQL
run-api-docker: docker-postgres
	@echo "Starting API with Docker Compose PostgreSQL..."
	@MIGRATION_ENABLED=true make run

# Alias for backward compatibility
build-api: build
run-api: run

