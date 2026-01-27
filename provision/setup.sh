#!/bin/bash
set -e

echo "=========================================="
echo "Setting up production environment"
echo "=========================================="

# Update system packages
apt-get update
apt-get upgrade -y

# Install required packages
apt-get install -y \
    docker.io \
    docker-compose \
    nginx \
    wget \
    curl \
    git \
    make

# Start and enable Docker
systemctl start docker
systemctl enable docker

# Add vagrant user to docker group (to run docker without sudo)
usermod -aG docker vagrant

# Install Docker Compose V2 (if not already installed)
if ! command -v docker compose &> /dev/null; then
    echo "Installing Docker Compose V2..."
    apt-get install -y docker-compose-plugin
fi

# Configure Nginx
echo "Configuring Nginx..."
rm -f /etc/nginx/sites-enabled/default
cp /vagrant/provision/nginx.conf /etc/nginx/sites-available/sre-chat-api
ln -sf /etc/nginx/sites-available/sre-chat-api /etc/nginx/sites-enabled/

# Test Nginx configuration
nginx -t

# Enable and start Nginx
systemctl enable nginx
systemctl restart nginx

echo "=========================================="
echo "Setup completed successfully!"
echo "=========================================="
