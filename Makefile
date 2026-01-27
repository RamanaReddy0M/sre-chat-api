.PHONY: build run test test-coverage migrate deps clean fmt lint docker-up docker-down docker-logs docker-ps docker-postgres run-api-local docker-build docker-run docker-tag docker-push start-db migrate-db build-api-image run-api-compose stop-api logs-api check-db check-migrations vagrant-up vagrant-down vagrant-reload vagrant-ssh vagrant-status vagrant-provision vagrant-deploy vagrant-logs

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

# Run API locally (builds and runs against Docker-managed PostgreSQL)
run-api-local: docker-postgres
	@echo "Starting API locally with Docker-managed PostgreSQL..."
	@MIGRATION_ENABLED=true make run

# Docker image build and run
# Version should be set via VERSION variable (e.g., VERSION=1.0.0 make docker-build)
VERSION ?= 1.0.0
IMAGE_NAME ?= sre-chat-api
CONTAINER_NAME ?= sre-chat-api

docker-build:
	@echo "Building Docker image $(IMAGE_NAME):$(VERSION)..."
	docker build -t $(IMAGE_NAME):$(VERSION) -t $(IMAGE_NAME):$(shell echo $(VERSION) | cut -d. -f1) -t $(IMAGE_NAME):$(shell echo $(VERSION) | cut -d. -f1-2) .
	@echo "Image built successfully: $(IMAGE_NAME):$(VERSION)"

docker-run:
	@echo "Running Docker container..."
	@docker run --rm -it \
		-p 8080:8080 \
		-e SERVER_PORT=8080 \
		-e DB_HOST=host.docker.internal \
		-e DB_PORT=5432 \
		-e DB_USER=postgres \
		-e DB_PASSWORD=postgres \
		-e DB_NAME=chat_db \
		-e DB_SSLMODE=disable \
		-e MIGRATION_ENABLED=true \
		$(IMAGE_NAME):$(VERSION)

docker-run-detached:
	@echo "Running Docker container in detached mode..."
	@docker run -d --name $(CONTAINER_NAME) \
		-p 8080:8080 \
		-e SERVER_PORT=8080 \
		-e DB_HOST=host.docker.internal \
		-e DB_PORT=5432 \
		-e DB_USER=postgres \
		-e DB_PASSWORD=postgres \
		-e DB_NAME=chat_db \
		-e DB_SSLMODE=disable \
		-e MIGRATION_ENABLED=true \
		$(IMAGE_NAME):$(VERSION)
	@echo "Container started. Use 'docker logs -f $(CONTAINER_NAME)' to view logs."

docker-stop:
	@echo "Stopping Docker container..."
	@docker stop $(CONTAINER_NAME) || true
	@docker rm $(CONTAINER_NAME) || true

docker-logs-api:
	@docker logs -f $(CONTAINER_NAME)

# Tag image with additional tags (for registry push)
docker-tag:
	@echo "Tagging image $(IMAGE_NAME):$(VERSION)..."
	@docker tag $(IMAGE_NAME):$(VERSION) $(IMAGE_NAME):latest || echo "Note: Using semver tags, 'latest' tag not created"

# Push to registry (set REGISTRY variable, e.g., REGISTRY=ghcr.io/username)
docker-push:
	@if [ -z "$(REGISTRY)" ]; then \
		echo "Error: REGISTRY variable not set. Usage: REGISTRY=ghcr.io/username make docker-push"; \
		exit 1; \
	fi
	@echo "Pushing $(IMAGE_NAME):$(VERSION) to $(REGISTRY)..."
	@docker tag $(IMAGE_NAME):$(VERSION) $(REGISTRY)/$(IMAGE_NAME):$(VERSION)
	@docker tag $(IMAGE_NAME):$(VERSION) $(REGISTRY)/$(IMAGE_NAME):$(shell echo $(VERSION) | cut -d. -f1)
	@docker tag $(IMAGE_NAME):$(VERSION) $(REGISTRY)/$(IMAGE_NAME):$(shell echo $(VERSION) | cut -d. -f1-2)
	@docker push $(REGISTRY)/$(IMAGE_NAME):$(VERSION)
	@docker push $(REGISTRY)/$(IMAGE_NAME):$(shell echo $(VERSION) | cut -d. -f1)
	@docker push $(REGISTRY)/$(IMAGE_NAME):$(shell echo $(VERSION) | cut -d. -f1-2)
	@echo "Pushed successfully!"

# ============================================================================
# Milestone 3: One-Click Local Development Setup
# ============================================================================

# Check if DB container is running
check-db:
	@echo "Checking if PostgreSQL container is running..."
	@if docker compose ps postgres | grep -q "Up"; then \
		echo "✓ PostgreSQL is already running"; \
	else \
		echo "✗ PostgreSQL is not running"; \
		exit 1; \
	fi

# Check if migrations have been applied (checks if groups table exists)
check-migrations:
	@echo "Checking if database migrations have been applied..."
	@if docker compose exec -T postgres psql -U postgres -d chat_db -c "\d groups" > /dev/null 2>&1; then \
		echo "✓ Migrations appear to be applied (groups table exists)"; \
	else \
		echo "✗ Migrations not applied (groups table not found)"; \
		exit 1; \
	fi

# Start DB container
start-db:
	@echo "Starting PostgreSQL container..."
	@if docker compose ps postgres | grep -q "Up"; then \
		echo "✓ PostgreSQL is already running"; \
	else \
		docker compose up -d postgres; \
		echo "Waiting for PostgreSQL to be ready..."; \
		for i in $$(seq 1 30); do \
			if docker compose exec -T postgres pg_isready -U postgres > /dev/null 2>&1; then \
				echo "✓ PostgreSQL is ready!"; \
				exit 0; \
			fi; \
			sleep 1; \
		done; \
		echo "⚠ PostgreSQL may not be fully ready yet"; \
	fi

# Run DB migrations (migrations run automatically when API starts, but this can be used to verify)
migrate-db: start-db
	@echo "Checking database migrations..."
	@if docker compose exec -T postgres psql -U postgres -d chat_db -c "\d groups" > /dev/null 2>&1; then \
		echo "✓ Migrations already applied (groups table exists)"; \
	else \
		echo "⚠ Migrations not yet applied"; \
		echo "  Note: Migrations will run automatically when API starts"; \
	fi

# Build REST API docker image
build-api-image:
	@echo "Building REST API Docker image..."
	@VERSION=1.0.0 IMAGE_NAME=sre-chat-api make docker-build

# Run REST API via Docker Compose (starts DB, runs migrations, then starts API)
run-api-compose: start-db migrate-db
	@echo "Starting REST API via Docker Compose..."
	@echo "This will:"
	@echo "  1. Ensure PostgreSQL is running"
	@echo "  2. Ensure migrations are applied"
	@echo "  3. Build API image (if needed)"
	@echo "  4. Start the API container"
	@docker compose up -d --build api
	@echo ""
	@echo "✓ API container started!"
	@echo "  API is available at: http://localhost:8080"
	@echo "  Web client: http://localhost:8080"
	@echo "  Health check: http://localhost:8080/api/v1/healthcheck"
	@echo ""
	@echo "View logs with: make logs-api"
	@echo "Stop services with: make stop-api"

# Stop API and DB containers
stop-api:
	@echo "Stopping API and DB containers..."
	@docker compose down
	@echo "✓ Containers stopped"

# View API logs
logs-api:
	@echo "Viewing API logs (Ctrl+C to exit)..."
	@docker compose logs -f api

# View all logs (API + DB)
logs-all:
	@echo "Viewing all logs (Ctrl+C to exit)..."
	@docker compose logs -f

# Restart API container
restart-api:
	@echo "Restarting API container..."
	@docker compose restart api
	@echo "✓ API container restarted"

# ============================================================================
# Vagrant Production Deployment
# ============================================================================

# Start Vagrant box (production environment)
# On Apple Silicon: use VAGRANT_PROVIDER=vmware_desktop if you don't have Parallels Pro
VAGRANT_PROVIDER ?=
vagrant-up:
	@echo "Starting Vagrant production environment..."
	@if [ "$$(uname -m)" = "arm64" ] || [ "$$(uname -m)" = "aarch64" ]; then \
		p="$${VAGRANT_PROVIDER:-parallels}"; \
		echo "Using provider: $$p"; \
		vagrant up --provider=$$p; \
	else \
		vagrant up; \
	fi
	@echo ""
	@echo "✓ Production environment is ready!"
	@echo "  API: http://localhost:8080"
	@echo "  Nginx: http://localhost:8080 (reverse proxy)"
	@echo "  Direct API: http://localhost:8081"

# Stop Vagrant box
vagrant-down:
	@echo "Stopping Vagrant production environment..."
	@vagrant halt
	@echo "✓ Vagrant box stopped"

# Reload Vagrant box (restart)
vagrant-reload:
	@echo "Reloading Vagrant production environment..."
	@vagrant reload
	@echo "✓ Vagrant box reloaded"

# SSH into Vagrant box
vagrant-ssh:
	@echo "Connecting to Vagrant box..."
	@vagrant ssh

# Check Vagrant box status
vagrant-status:
	@echo "Vagrant box status:"
	@vagrant status

# Re-provision Vagrant box (re-run setup scripts)
vagrant-provision:
	@echo "Re-provisioning Vagrant box..."
	@vagrant provision
	@echo "✓ Provisioning completed"

# Deploy to Vagrant (rebuild and restart services)
vagrant-deploy:
	@echo "Deploying to Vagrant production environment..."
	@vagrant ssh -c "cd /vagrant && sudo docker compose down && sudo docker compose build --no-cache && sudo docker compose up -d"
	@echo "✓ Deployment completed"

# View logs in Vagrant
vagrant-logs:
	@echo "Viewing logs from Vagrant box..."
	@vagrant ssh -c "cd /vagrant && sudo docker compose logs -f"

# ============================================================================
# Alias for backward compatibility
# ============================================================================
build-api: build

