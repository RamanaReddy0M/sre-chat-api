# SRE Chat API - SRE Bootcamp Exercise

A REST API for a real-time chat application built with Go and Gin framework. Features PostgreSQL persistence, Server-Sent Events (SSE) for real-time updates, and a beautiful web interface.

## Web Client

A modern web interface is included - no need for curl or Postman!

![Chat Interface](chat-ui.png)

### Quick Start

1. **Start the API:**

   ```bash
   make run-api-docker
   ```

2. **Open in Browser:**

   ```text
   http://localhost:8080
   ```

3. **Join Chat:**
   - Enter your name and click "Join Chat"
   - Start chatting with real-time updates!

### Web Client Features

- ✅ Real-time message updates via SSE
- ✅ Edit and delete your own messages
- ✅ Automatic link previews for URLs
- ✅ Clickable email addresses and phone numbers
- ✅ Modern dark theme UI

### SSE Architecture

The application uses Server-Sent Events (SSE) for real-time message broadcasting:

![SSE Architecture Diagram](sse-architecture.svg)

**How it works:**

1. **Client Connection**: Web clients connect to `/api/v1/messages/stream` endpoint
2. **Registration**: API registers each client with the SSE Handler, maintaining a map of client channels
3. **Message Creation**: When a client posts a message, it's saved to PostgreSQL
4. **Broadcasting**: SSE Handler broadcasts the new message to all connected clients in real-time
5. **Keep-Alive**: Server sends ping events every 30 seconds to maintain connections

## API Features

- ✅ REST API with proper HTTP verbs (GET, POST, PUT, DELETE)
- ✅ API versioning (`/api/v1/`)
- ✅ PostgreSQL database with migrations
- ✅ Environment variable configuration
- ✅ Structured logging with zap
- ✅ Health check endpoint
- ✅ Unit tests for all endpoints
- ✅ Postman collection included
- ✅ Dockerized with multi-stage Dockerfile
- ✅ Optimized Docker image size (~65MB)
- ✅ CI/CD pipeline with GitHub Actions (self-hosted runner)

## Prerequisites

### Required Tools

The following tools must be installed on your local machine:

- **Docker** (version 20.10 or higher) - [Install Docker](https://docs.docker.com/get-docker/)
- **Docker Compose** (version 2.0 or higher) - Usually included with Docker Desktop
- **Make** - Usually pre-installed on Linux/macOS. For Windows, use [Chocolatey](https://chocolatey.org/) or [WSL](https://docs.microsoft.com/en-us/windows/wsl/)

### Optional Tools (for local development without Docker)

- Go 1.21 or higher
- PostgreSQL 12 or higher

## One-Click Local Development Setup

This setup uses Docker Compose to run both the API and PostgreSQL database with minimal configuration.

### Quick Start (Recommended)

```bash
# 1. Clone the repository
git clone git@github.com:RamanaReddy0M/sre-chat-api.git
cd sre-chat-api

# 2. Start everything (DB + API)
# This will automatically build the image if needed
make run-api
```

That's it! The API will be available at `http://localhost:8080`.

**Note:** The `run-api` target automatically builds the Docker image if it doesn't exist, so no separate build step is needed.

### Make Targets and Execution Order

The following Make targets are available for local development:

#### 1. **`make start-db`** - Start Database Container
   - Starts the PostgreSQL container
   - Waits for PostgreSQL to be ready
   - Skips if already running

#### 2. **`make migrate-db`** - Check Database Migrations
   - Ensures DB is running (calls `start-db`)
   - Checks if migrations are already applied
   - Note: Migrations run automatically when API starts (via `MIGRATION_ENABLED=true`)

#### 3. **`make build-api-image`** - Build REST API Docker Image
   - Builds the Docker image for the API
   - Uses multi-stage Dockerfile for optimal size

#### 4. **`make run-api`** - Run REST API Container (One-Click Setup)
   - **Automatically:**
     1. Starts the database container (if not running)
     2. Checks and runs migrations (if needed)
     3. Starts the API container
   - **Includes checks:**
     - Verifies DB is running before starting API
     - Checks if migrations are already applied
   - API available at: `http://localhost:8080`

#### Additional Useful Targets

- **`make stop-api`** - Stop all containers (API + DB)
- **`make logs-api`** - View API logs (follow mode)
- **`make logs-all`** - View all logs (API + DB)
- **`make restart-api`** - Restart API container only
- **`make check-db`** - Check if DB container is running
- **`make check-migrations`** - Check if migrations are applied

### Step-by-Step Manual Setup

If you prefer to run steps individually:

```bash
# 1. Start the database
make start-db

# 2. Run migrations (optional - runs automatically on API start)
make migrate-db

# 3. Build the API image (optional - run-api builds it automatically)
make build-api-image

# 4. Start the API
docker compose up -d api

# Or use the one-command approach (recommended):
make run-api
```

### Verify Setup

```bash
# Check if containers are running
docker compose ps

# Check API health
curl http://localhost:8080/api/v1/healthcheck

# View API logs
make logs-api
```

### Troubleshooting

**Database connection issues:**
- Ensure PostgreSQL container is running: `make check-db`
- Check logs: `docker compose logs postgres`

**Migrations not applied:**
- Check migration status: `make check-migrations`
- Manually run migrations: `make migrate-db`

**API not starting:**
- Check logs: `make logs-api`
- Verify DB is ready: `make check-db`
- Restart everything: `make stop-api && make run-api`

## Local Development (Without Docker)

If you prefer to run locally without Docker:

### 1. Clone and Install

```bash
git clone git@github.com:RamanaReddy0M/sre-chat-api.git
cd sre-chat-api
make deps
```

### 2. Database Setup

#### Option A: Docker Compose (PostgreSQL only)

```bash
make docker-up
```

#### Option B: Local PostgreSQL

```bash
createdb -U postgres chat_db
```

### 3. Configure Environment

```bash
cp env.example .env
# Edit .env with your database credentials
```

Or set environment variables:

```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=postgres
export DB_NAME=chat_db
export DB_SSLMODE=disable
export MIGRATION_ENABLED=true
```

### 4. Run

```bash
make build
make run
```

The server starts on `http://localhost:8080`. Migrations run automatically on startup.

## Docker

The application can be run using Docker with a multi-stage Dockerfile for optimal image size.

### Build Docker Image

```bash
# Build with version tag (semver)
VERSION=1.0.0 make docker-build

# Or manually
docker build -t sre-chat-api:1.0.0 .
```

### Run Docker Container

```bash
# Run with environment variables
VERSION=1.0.0 make docker-run

# Or manually with custom environment variables
docker run --rm -it \
  -p 8080:8080 \
  -e DB_HOST=host.docker.internal \
  -e DB_PORT=5432 \
  -e DB_USER=postgres \
  -e DB_PASSWORD=postgres \
  -e DB_NAME=chat_db \
  -e MIGRATION_ENABLED=true \
  sre-chat-api:1.0.0
```

### Docker with Docker Compose

```bash
# Start PostgreSQL
make docker-up

# Run API container (connects to PostgreSQL)
VERSION=1.0.0 make docker-run-detached

# View logs
make docker-logs-api

# Stop container
make docker-stop
```

**Note:** When running in Docker, use `host.docker.internal` (Mac/Windows) or your host IP (Linux) for `DB_HOST` to connect to PostgreSQL running on the host machine.

### Docker Image Details

- **Multi-stage build**: Separate build and runtime stages for optimal image size
- **Image size**: ~65MB (Alpine-based, statically compiled binary)
- **Security**: Runs as non-root user (`appuser`)
- **Health check**: Built-in health check endpoint monitoring
- **Semver tagging**: Uses semantic versioning (e.g., `1.0.0`, `1.0`, `1`) instead of `latest`

### Available Make Targets

```bash
# Build Docker image
VERSION=1.0.0 make docker-build

# Run container interactively
VERSION=1.0.0 make docker-run

# Run container in background
VERSION=1.0.0 make docker-run-detached

# View container logs
make docker-logs-api

# Stop container
make docker-stop

# Push to registry (set REGISTRY variable)
REGISTRY=your-dockerhub-username VERSION=1.0.0 make docker-push
```

## API Endpoints

All endpoints are versioned under `/api/v1/`:

- **GET** `/api/v1/healthcheck` - Health check
- **POST** `/api/v1/messages` - Create message
- **GET** `/api/v1/messages` - List messages
- **GET** `/api/v1/messages/:id` - Get message
- **PUT** `/api/v1/messages/:id` - Update message
- **DELETE** `/api/v1/messages/:id` - Delete message

### Example

```bash
# Create a message
curl -X POST http://localhost:8080/api/v1/messages \
  -H "Content-Type: application/json" \
  -d '{"user": "Alice", "content": "Hello!"}'

# Get all messages
curl http://localhost:8080/api/v1/messages
```

## Testing

```bash
make test              # Run unit tests
make test-coverage     # Run with coverage report
```

Import `postman_collection.json` into Postman for API testing.

## Project Structure

```text
sre-chat-api/
├── cmd/api/main.go           # Application entry point
├── internal/
│   ├── config/               # Configuration management
│   ├── database/              # Database connection & migrations
│   ├── handlers/              # API handlers & tests
│   ├── middleware/            # Logging middleware
│   ├── models/                # Data models
│   └── rest.go                # Router setup
├── migrations/                # SQL migration files
├── web/                       # Web client interface
├── go.mod                     # Dependencies
├── Makefile                   # Build commands
└── postman_collection.json    # Postman collection
```

## Makefile Commands

```bash
make build          # Build the API
make run            # Run the API
make run-api-docker # Start PostgreSQL and run API
make test           # Run tests
make deps           # Install dependencies
make clean          # Clean build artifacts
```

## Troubleshooting

### Database Connection Failed

- Check PostgreSQL is running: `make docker-ps` or `pg_isready`
- Verify credentials in environment variables

### Port Already in Use

- Change `SERVER_PORT` environment variable

### Tests Failing

- Create test database: `createdb -U postgres chat_test_db`

## Production Deployment with Vagrant

The project includes a Vagrant-based production environment setup with Nginx as a reverse proxy.

### Prerequisites

- **Vagrant** (2.2.0 or higher) - [Install Vagrant](https://www.vagrantup.com/downloads)
- **Provider:** VirtualBox on Intel/AMD; on Apple Silicon: **Parallels Desktop Pro** or **VMware Fusion** (see below)

### ⚠️ Apple Silicon (M1/M2/M3) Note

**VirtualBox does not support Apple Silicon.** Use one of these providers:

**Option A — VMware Fusion (free for personal use)**

1. Install [VMware Fusion](https://www.vmware.com/products/fusion.html) (free personal license). On Apple Silicon you typically get **VMware Fusion Tech Preview**, which installs as `VMware Fusion Tech Preview.app`. The Vagrant VMware Utility expects `VMware Fusion.app`, so create a symlink: `sudo ln -sf "/Applications/VMware Fusion Tech Preview.app" "/Applications/VMware Fusion.app"`.
2. Install the Vagrant VMware plugin:
   ```bash
   vagrant plugin install vagrant-vmware-desktop
   ```
3. Install the **Vagrant VMware Utility** (required by the plugin). Download the macOS installer from [Vagrant VMware Utility](https://www.vagrantup.com/downloads/vmware), run it, and follow the prompts. See [Vagrant VMware Utility docs](https://www.vagrantup.com/docs/providers/vmware/vagrant-vmware-utility) if needed.
4. Start the VM with VMware:
   ```bash
   VAGRANT_PROVIDER=vmware_desktop make vagrant-up
   ```
   Or: `vagrant up --provider=vmware_desktop`

**Option B — Parallels Desktop (Pro/Business edition required)**

1. Install **Parallels Desktop** (Pro, Business, or Enterprise — the free edition does not include `prlctl`).
2. Install the Vagrant Parallels plugin: `vagrant plugin install vagrant-parallels`
3. Run `make vagrant-up` (Parallels is the default on ARM64).

**Option C — Docker (no VM)**

   ```bash
   make run-api
   ```

### Quick Start

```bash
# Start the production environment
make vagrant-up

# Access the API
# - Via Nginx (reverse proxy): http://localhost:8080
# - Direct API access: http://localhost:8081
```

### Architecture

```
┌─────────────────────────────────────────┐
│         Host Machine (Your PC)          │
│                                         │
│  http://localhost:8080 (Nginx Proxy)   │
│  http://localhost:8081 (Direct API)    │
└──────────────┬──────────────────────────┘
               │
               │ Port Forwarding
               │
┌──────────────▼──────────────────────────┐
│      Vagrant VM (Production)            │
│  ┌──────────────────────────────────┐  │
│  │         Nginx (Port 80)          │  │
│  │      Reverse Proxy               │  │
│  └──────────┬───────────────────────┘  │
│             │                          │
│             │ proxy_pass               │
│             │                          │
│  ┌──────────▼───────────────────────┐  │
│  │   API Container (Port 8080)      │  │
│  │   Docker Compose                 │  │
│  └──────────┬───────────────────────┘  │
│             │                          │
│  ┌──────────▼───────────────────────┐  │
│  │  PostgreSQL (Port 5432)          │  │
│  │  Docker Container                │  │
│  └──────────────────────────────────┘  │
└─────────────────────────────────────────┘
```

### Vagrant Commands

```bash
# Start production environment
make vagrant-up

# Stop production environment
make vagrant-down

# Reload (restart) Vagrant box
make vagrant-reload

# SSH into Vagrant box
make vagrant-ssh

# Check status
make vagrant-status

# Re-provision (re-run setup scripts)
make vagrant-provision

# Deploy updates (rebuild containers)
make vagrant-deploy

# View logs
make vagrant-logs
```

### Nginx Configuration

Nginx is configured as a reverse proxy with:

- **Rate limiting**: 10 requests/second with burst of 20
- **Health check endpoint**: `/api/v1/healthcheck` (no rate limiting)
- **SSE support**: Proper configuration for Server-Sent Events
- **Security headers**: X-Frame-Options, X-Content-Type-Options, etc.
- **Static asset caching**: Cached for 1 hour
- **CORS headers**: Enabled for API endpoints

### Port Mapping

- **8080** (host) → **80** (guest): Nginx reverse proxy
- **8081** (host) → **8080** (guest): Direct API access
- **5433** (host) → **5432** (guest): PostgreSQL (for debugging)

### Accessing the Application

After starting the Vagrant box:

- **Web Client**: http://localhost:8080
- **API Health Check**: http://localhost:8080/api/v1/healthcheck
- **Direct API**: http://localhost:8081/api/v1/healthcheck

### Troubleshooting

**Vagrant box won't start:**
- **Apple Silicon users**: Use Parallels (see "Apple Silicon Note" above). VirtualBox is not supported on ARM Macs.
- **"Provider 'vmware_desktop' could not be found"**: Install the plugin first: `vagrant plugin install vagrant-vmware-desktop`, then run `VAGRANT_PROVIDER=vmware_desktop make vagrant-up`.
- **"Vagrant VMware Utility" / "vagrant-utility.client.crt" error**: The VMware plugin needs the **Vagrant VMware Utility**. Install it from [vagrantup.com/downloads/vmware](https://www.vagrantup.com/downloads/vmware) (download the macOS package, run the installer, then retry).
- **"Connection refused" to 127.0.0.1:9922 (Vagrant VMware Utility)**: The VMware Utility service is not running or is exiting. On macOS, start it with:
  ```bash
  sudo launchctl load -w /Library/LaunchDaemons/com.vagrant.vagrant-vmware-utility.plist
  ```
  Or, on newer macOS: `sudo launchctl bootstrap system /Library/LaunchDaemons/com.vagrant.vagrant-vmware-utility.plist`.  
  If the plist path doesn’t exist, reinstall the [Vagrant VMware Utility](https://www.vagrantup.com/downloads/vmware).  
  **If it still fails after a successful load:** (1) Ensure VMware Fusion is installed and has been opened at least once. (2) Check whether the job is running: `sudo launchctl list | grep vagrant`. (3) Check the utility log: `sudo cat "/Library/Application Support/vagrant-vmware-utility/service.log"` — it often explains why the process exits. (4) If the service keeps failing, use **Docker** (`make run-api`) or **Parallels Desktop Pro** instead.
- **"Failed to locate VMware installation directory!"** in `service.log` or connection refused to 9922: The Vagrant VMware Utility looks for **VMware Fusion** at `/Applications/VMware Fusion.app`.  
  **If you use VMware Fusion Tech Preview:** create a symlink, then fully restart the service:
  ```bash
  sudo ln -sf "/Applications/VMware Fusion Tech Preview.app" "/Applications/VMware Fusion.app"
  sudo launchctl bootout system/com.vagrant.vagrant-vmware-utility
  sudo launchctl bootstrap system /Library/LaunchDaemons/com.vagrant.vagrant-vmware-utility.plist
  ```
  **If it still fails:** (1) Confirm Fusion Tech Preview is installed: `ls -la "/Applications/VMware Fusion Tech Preview.app"` — the symlink only works if that app exists. (2) Reinstall the [Vagrant VMware Utility](https://www.vagrantup.com/downloads/vmware) (others report this fixing detection). (3) If you don’t have Fusion or Tech Preview, install it from [VMware](https://www.vmware.com/products/fusion.html), or use **Docker** (`make run-api`) / **Parallels Desktop Pro** instead.
- **"Load failed: 5: Input/output error"** when loading the Vagrant VMware Utility plist: This is a known issue on recent macOS. Try in order: (1) Confirm the plist exists: `ls /Library/LaunchDaemons/com.vagrant.vagrant-vmware-utility.plist`. (2) Get more detail: `sudo launchctl bootstrap system /Library/LaunchDaemons/com.vagrant.vagrant-vmware-utility.plist`. (3) If the job is already loaded, unload first: `sudo launchctl bootout system /Library/LaunchDaemons/com.vagrant.vagrant-vmware-utility.plist`, then run the installer again or run the `launchctl load` command. (4) Reinstall the [Vagrant VMware Utility](https://www.vagrantup.com/downloads/vmware). If it still fails, use **Docker for local dev** (`make run-api`) or **Parallels Desktop Pro** instead of VMware.
- Check VirtualBox is installed and running (Intel Macs only)
- Verify virtualization is enabled in BIOS (Intel Macs)
- Check available disk space (VM needs ~10GB)
- Ensure Vagrant version supports your VirtualBox version (upgrade if needed)

**Services not starting:**
- SSH into box: `make vagrant-ssh`
- Check Docker: `sudo systemctl status docker`
- Check logs: `make vagrant-logs`

**Nginx not working:**
- Check Nginx status: `sudo systemctl status nginx`
- Test config: `sudo nginx -t`
- View logs: `sudo tail -f /var/log/nginx/sre-chat-api-error.log`

**API not responding:**
- Check containers: `sudo docker compose ps`
- View API logs: `sudo docker compose logs api`
- Check health: `curl http://localhost:8080/api/v1/healthcheck`

## License

This project is part of the SRE Bootcamp exercises.
