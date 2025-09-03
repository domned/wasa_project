#!/bin/bash

# WASAText Deployment Script (Docker only)

set -e

echo "🚀 Starting WASAText deployment..."

# Create data directory if it doesn't exist
mkdir -p data

# Build backend image
echo "📦 Building backend container..."
docker build -f Dockerfile.backend -t wasa-backend .

# Build frontend image  
echo "📦 Building frontend container..."
docker build -f Dockerfile.frontend -t wasa-frontend .

# Create network
echo "🔗 Creating Docker network..."
docker network create wasa-network 2>/dev/null || true

# Stop and remove existing containers
echo "🧹 Cleaning up existing containers..."
docker stop wasa-backend wasa-frontend 2>/dev/null || true
docker rm wasa-backend wasa-frontend 2>/dev/null || true

# Start backend
echo "🔧 Starting backend service..."
docker run -d \
  --name wasa-backend \
  --network wasa-network \
  -p 3000:3000 \
  -v "$(pwd)/data:/data" \
  -e DATABASE_PATH=/data/app.db \
  -e PORT=3000 \
  wasa-backend

# Start frontend
echo "🔧 Starting frontend service..."
docker run -d \
  --name wasa-frontend \
  --network wasa-network \
  -p 80:80 \
  wasa-frontend

echo "⏳ Waiting for services to be ready..."
sleep 10

# Check if services are running
if docker ps | grep -q wasa-backend && docker ps | grep -q wasa-frontend; then
    echo "✅ Deployment successful!"
    echo ""
    echo "🌐 Frontend: http://localhost"
    echo "🔧 Backend API: http://localhost:3000"
    echo ""
    echo "📋 View backend logs: docker logs -f wasa-backend"
    echo "📋 View frontend logs: docker logs -f wasa-frontend"
    echo "🛑 Stop services: docker stop wasa-backend wasa-frontend"
else
    echo "❌ Deployment failed. Check logs:"
    echo "Backend: docker logs wasa-backend"
    echo "Frontend: docker logs wasa-frontend"
    exit 1
fi
