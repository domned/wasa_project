#!/bin/bash

# WASAText Deployment Script

set -e

echo "🚀 Starting WASAText deployment..."

# Check if Docker is running
if ! docker info >/dev/null 2>&1; then
    echo "❌ Docker is not running or not installed"
    echo ""
    echo "📋 Please install and start Docker:"
    echo "   1. Install Docker Desktop: https://www.docker.com/products/docker-desktop/"
    echo "   2. Start Docker Desktop"
    echo "   3. Wait for the Docker whale icon to appear in your menu bar"
    echo ""
    echo "🔄 Alternative: Use local development setup:"
    echo "   Terminal 1: go run ./cmd/webapi"
    echo "   Terminal 2: cd webui && npm run dev"
    echo ""
    echo "📖 See DOCKER_SETUP.md for detailed instructions"
    exit 1
fi

# Check if docker-compose is available
if ! command -v docker-compose &> /dev/null; then
    echo "⚠️  docker-compose not found, falling back to Docker-only deployment..."
    echo "🔄 Running ./deploy-docker.sh instead..."
    exec ./deploy-docker.sh
fi

# Create data directory if it doesn't exist
mkdir -p data

# Build and start the containers
echo "📦 Building Docker containers..."
docker-compose build

echo "🔧 Starting services..."
docker-compose up -d

echo "⏳ Waiting for services to be ready..."
sleep 10

# Check if services are running
if docker-compose ps | grep -q "Up"; then
    echo "✅ Deployment successful!"
    echo ""
    echo "🌐 Frontend: http://localhost"
    echo "🔧 Backend API: http://localhost:3000"
    echo ""
    echo "📋 View logs with: docker-compose logs -f"
    echo "🛑 Stop services with: docker-compose down"
else
    echo "❌ Deployment failed. Check logs with: docker-compose logs"
    exit 1
fi
