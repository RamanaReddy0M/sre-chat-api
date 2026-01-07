# SRE Chat API - SRE Bootcamp Exercise

A REST API for a real-time chat application built with Go and Gin framework, following the SRE Bootcamp exercise requirements. This application allows users to send messages and view chat history stored in PostgreSQL. Features real-time updates via Server-Sent Events (SSE) and a beautiful web interface. The default group "SRE Bootcamp" is automatically created on first run.

## Purpose

This repository implements a REST API for a chat application as part of the SRE Bootcamp exercises. The API allows anyone to send messages at any time and view all chat history. Messages are persisted in PostgreSQL, and a default group called "SRE Bootcamp" is automatically created.

## Features

- ✅ REST API with proper HTTP verbs (GET, POST, PUT, DELETE)
- ✅ API versioning (`/api/v1/`)
- ✅ PostgreSQL database for message persistence
- ✅ Database migrations using GORM
- ✅ Environment variable configuration (Twelve-Factor App methodology)
- ✅ Structured logging with appropriate log levels (using zap)
- ✅ Health check endpoint (`/api/v1/healthcheck`)
- ✅ Unit tests for all endpoints
- ✅ Postman collection for API testing
- ✅ Makefile for build and run automation
- ✅ Default "SRE Bootcamp" group automatically created

## Prerequisites

- Go 1.21 or higher
- PostgreSQL 12 or higher
- Make (optional, but recommended)

## Local Setup Instructions

### 1. Clone the Repository
```bash
git clone <repository-url>
cd sre-chat-api
```

### 2. Install Dependencies
```bash
make deps
# Or manually:
go mod download
```

### 3. Set Up PostgreSQL Database

#### Option A: Using Docker Compose (Recommended)

The easiest way to set up PostgreSQL is using Docker Compose:

```bash
# Start PostgreSQL in Docker
make docker-up

# Or manually:
docker-compose up -d

# Check if PostgreSQL is running
make docker-ps

# View PostgreSQL logs
make docker-logs
```

This will start PostgreSQL on port `5432` with:
- Database: `chat_db`
- User: `postgres`
- Password: `postgres`

#### Option B: Local PostgreSQL Installation

If you have PostgreSQL installed locally, create a database:
```bash
# Using psql
psql -U postgres
CREATE DATABASE chat_db;
\q

# Or using createdb command
createdb -U postgres chat_db
```

### 4. Configure Environment Variables

Copy the example environment file and update with your values:

```bash
cp env.example .env
# Edit .env with your database credentials
```

Or set environment variables directly:

```bash
export SERVER_PORT=8080
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=postgres
export DB_NAME=chat_db
export DB_SSLMODE=disable
export MIGRATION_ENABLED=true
```

**Note:** Configuration is loaded from environment variables following the Twelve-Factor App methodology. No hard-coded values in the code. See `env.example` for all available configuration options.

### 5. Run Database Migrations

Migrations run automatically when the server starts **if `MIGRATION_ENABLED=true`** (default). The default "SRE Bootcamp" group is created automatically when migrations are enabled.

**Migration Control:**
- Set `MIGRATION_ENABLED=true` to run migrations on startup (default)
- Set `MIGRATION_ENABLED=false` to skip migrations (useful in production where migrations are run separately)

Alternatively, you can run the SQL migration manually:
```bash
psql -U postgres -d chat_db -f migrations/001_create_tables.sql
```

### 6. Build and Run the API

#### Quick Start with Docker Compose

The easiest way to run everything:

```bash
# Start PostgreSQL and run the API (one command)
make run-api-docker
```

This will:
1. Start PostgreSQL in Docker
2. Wait for PostgreSQL to be ready
3. Run migrations automatically
4. Start the API server

#### Manual Steps

```bash
# Build the API
make build

# Run the API
make run

# Or build and run in one command
make run
```

**With Docker Compose PostgreSQL:**
```bash
# Start PostgreSQL
make docker-up

# Run the API (migrations will run automatically)
MIGRATION_ENABLED=true make run
```

The server will start on `http://localhost:8080` (or the port specified in `SERVER_PORT`).

You should see:
```
Database connection established successfully
Database migrations completed successfully
Default group 'SRE Bootcamp' created successfully
Starting server address=0.0.0.0:8080
```

## Web Client

A beautiful web interface is included! No need for curl or Postman.

### Quick Start

1. **Start the API:**
   ```bash
   make run-api-docker
   ```

2. **Open in Browser:**
   ```
   http://localhost:8080
   ```

3. **Join Chat:**
   - Enter your name
   - Click "Join Chat"
   - Start chatting!

### Features

- ✅ Join with your name
- ✅ Real-time message updates (polls every 2 seconds)
- ✅ See all chat history
- ✅ Beautiful, modern UI
- ✅ Your messages are highlighted


## API Endpoints

All endpoints are versioned under `/api/v1/`:

### Health Check
- **GET** `/api/v1/healthcheck` - Check API health status

### Messages
- **POST** `/api/v1/messages` - Create a new message
- **GET** `/api/v1/messages` - Get all messages (optionally filter by `?group_id=1`)
- **GET** `/api/v1/messages/:id` - Get a specific message by ID
- **PUT** `/api/v1/messages/:id` - Update a message
- **DELETE** `/api/v1/messages/:id` - Delete a message

### Example Requests

**Create a message:**
```bash
curl -X POST http://localhost:8080/api/v1/messages \
  -H "Content-Type: application/json" \
  -d '{
    "user": "Alice",
    "content": "Hello, SRE Bootcamp!"
  }'
```

**Get all messages:**
```bash
curl http://localhost:8080/api/v1/messages
```

**Get a specific message:**
```bash
curl http://localhost:8080/api/v1/messages/1
```

**Update a message:**
```bash
curl -X PUT http://localhost:8080/api/v1/messages/1 \
  -H "Content-Type: application/json" \
  -d '{
    "content": "Updated message content"
  }'
```

**Delete a message:**
```bash
curl -X DELETE http://localhost:8080/api/v1/messages/1
```

**Health check:**
```bash
curl http://localhost:8080/api/v1/healthcheck
```

## Testing

### Run Unit Tests
```bash
make test
```

### Run Tests with Coverage
```bash
make test-coverage
```

This will generate a `coverage.html` file you can open in your browser.

### Using Postman

Import the `postman_collection.json` file into Postman:
1. Open Postman
2. Click "Import"
3. Select `postman_collection.json`
4. Set the `base_url` variable to `http://localhost:8080`
5. Start testing the endpoints!

## Project Structure

```
sre-chat-api/
├── cmd/
│   └── api/
│       └── main.go              # REST API entry point
├── internal/
│   ├── config/
│   │   └── config.go            # Configuration management
│   ├── database/
│   │   └── database.go          # Database connection and migrations
│   ├── handlers/
│   │   ├── message_handler.go   # Message CRUD handlers
│   │   ├── message_handler_test.go  # Unit tests
│   │   ├── sse_handler.go      # Server-Sent Events handler
│   │   └── health_handler.go   # Health check handler
│   ├── middleware/
│   │   └── logger.go            # Logging middleware
│   └── models/
│       └── models.go            # Data models (Group, Message)
├── migrations/
│   └── 001_create_tables.sql    # SQL migration file
├── web/
│   └── index.html               # Web client interface
├── go.mod                       # Go module dependencies
├── go.sum                       # Go module checksums
├── Makefile                     # Build automation
├── postman_collection.json      # Postman API collection
└── README.md                    # This file
```

## How It Works

1. **Configuration** (`internal/config/config.go`):
   - Loads configuration from environment variables
   - No hard-coded values (Twelve-Factor App methodology)

2. **Database** (`internal/database/database.go`):
   - Establishes PostgreSQL connection using GORM
   - Runs automatic migrations on startup
   - Seeds the default "SRE Bootcamp" group

3. **Models** (`internal/models/models.go`):
   - `Group`: Represents a chat group
   - `Message`: Represents a chat message with user, content, and timestamps

4. **Handlers** (`internal/handlers/`):
   - REST API handlers for CRUD operations
   - Health check endpoint
   - Proper HTTP status codes and error handling

5. **Middleware** (`internal/middleware/logger.go`):
   - Structured logging with zap logger
   - Request/response logging with appropriate log levels
   - Panic recovery

6. **API Server** (`cmd/api/main.go`):
   - Sets up Gin router with versioned routes
   - Applies middleware
   - Starts HTTP server

## Key Features Implemented

- ✅ **API Versioning**: All endpoints under `/api/v1/`
- ✅ **Proper HTTP Verbs**: GET, POST, PUT, DELETE
- ✅ **Structured Logging**: Using zap with appropriate log levels (INFO, ERROR)
- ✅ **Health Check**: `/api/v1/healthcheck` endpoint
- ✅ **Database Migrations**: Automatic migrations using GORM
- ✅ **Environment Variables**: All configuration via env vars
- ✅ **Unit Tests**: Comprehensive test coverage
- ✅ **Postman Collection**: Ready-to-use API collection

## Makefile Commands

### Application Commands
```bash
make build          # Build the REST API
make run            # Build and run the REST API
make run-api-docker # Start PostgreSQL in Docker and run the API
make test           # Run all tests
make test-coverage  # Run tests with coverage report
make migrate        # Run database migrations (auto on startup)
make deps           # Install dependencies
make clean          # Clean build artifacts
make fmt            # Format code
```

### Docker Compose Commands
```bash
make docker-up      # Start PostgreSQL in Docker
make docker-down    # Stop and remove PostgreSQL container
make docker-logs    # View PostgreSQL logs
make docker-ps      # Check Docker Compose services status
make docker-postgres # Start PostgreSQL and wait for it to be ready
```

## SRE Bootcamp Milestone 1 Requirements

This implementation meets all requirements for Milestone 1:

✅ **Public GitHub Repository** - Ready for GitHub  
✅ **README.md** - Complete with setup instructions  
✅ **Dependency Management** - `go.mod` with explicit dependencies  
✅ **Makefile** - Build and run commands  
✅ **Database Migrations** - GORM auto-migration + SQL migration file  
✅ **Environment Variables** - All config via env vars (no hard-coding)  
✅ **Postman Collection** - `postman_collection.json` included  
✅ **API Versioning** - All endpoints under `/api/v1/`  
✅ **Proper HTTP Verbs** - GET, POST, PUT, DELETE  
✅ **Structured Logging** - zap logger with appropriate levels  
✅ **Health Check** - `/api/v1/healthcheck` endpoint  
✅ **Unit Tests** - Comprehensive test coverage

## Troubleshooting

**Database connection failed**
- Ensure PostgreSQL is running: `pg_isready` or `systemctl status postgresql`
- If using Docker Compose: `make docker-ps` to check container status
- Check database credentials in environment variables
- Verify database exists: `psql -U postgres -l`

**Docker Compose issues**
- Check logs: `make docker-logs`
- Ensure Docker is running: `docker ps`

**Port already in use**
- Change the `SERVER_PORT` environment variable (default: 8080)

**Migration errors**
- Ensure database user has CREATE TABLE permissions
- Check PostgreSQL logs for detailed error messages

**Tests failing**
- Ensure test database exists: `createdb -U postgres chat_test_db`
- Or skip tests that require database: `go test -v ./... -short`

## Development

### Adding New Endpoints

1. Add handler function in `internal/handlers/`
2. Register route in `cmd/api/main.go`
3. Add test in handler test file
4. Update Postman collection

### Database Changes

1. Modify models in `internal/models/models.go`
2. GORM will auto-migrate on next server start
3. Or create new SQL migration file in `migrations/`

## License

This project is part of the SRE Bootcamp exercises.

