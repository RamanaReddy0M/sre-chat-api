#!/bin/bash
set -e

echo "=========================================="
echo "Deploying SRE Chat API"
echo "=========================================="

cd /vagrant

# Ensure Docker is running
sudo systemctl start docker || true

# Build and start services using Docker Compose
echo "Building and starting services..."
sudo docker compose down || true
sudo docker compose build --no-cache
sudo docker compose up -d

# Wait for services to be ready
echo "Waiting for services to be ready..."
sleep 10

# Check if API is healthy
max_attempts=30
attempt=0
while [ $attempt -lt $max_attempts ]; do
    if curl -f http://localhost:8080/api/v1/healthcheck > /dev/null 2>&1; then
        echo "✓ API is healthy!"
        break
    fi
    attempt=$((attempt + 1))
    echo "Waiting for API... ($attempt/$max_attempts)"
    sleep 2
done

if [ $attempt -eq $max_attempts ]; then
    echo "⚠ API health check failed after $max_attempts attempts"
    echo "Checking logs..."
    sudo docker compose logs api
    exit 1
fi

# Reload Nginx to ensure it's using the latest configuration
sudo systemctl reload nginx

echo "=========================================="
echo "Deployment completed!"
echo "=========================================="
echo "API is available at: http://localhost:8080"
echo "Nginx reverse proxy: http://localhost:8080 (via port 80)"
echo "=========================================="
