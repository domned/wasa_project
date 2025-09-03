# Docker Setup and Deployment Guide

## Docker Installation Required

It looks like Docker is not installed or running on your system. Here's how to get it set up:

### Install Docker Desktop (Recommended)

1. **Download Docker Desktop:**

    - Go to https://www.docker.com/products/docker-desktop/
    - Download Docker Desktop for Mac
    - Install the .dmg file

2. **Start Docker Desktop:**
    - Open Docker Desktop from Applications
    - Wait for it to start (you'll see the Docker icon in your menu bar)
    - You may need to complete the initial setup

### Alternative: Install via Homebrew

```bash
# Install Docker
brew install --cask docker

# Or install Docker CLI tools only
brew install docker docker-compose
```

### Verify Installation

```bash
# Check Docker is running
docker --version
docker ps

# Check Docker Compose (if installed)
docker-compose --version
```

## Deployment Options

### Option 1: Docker Desktop + Docker Compose

Once Docker Desktop is installed and running:

```bash
./deploy.sh
```

### Option 2: Docker Only

If you only have Docker CLI (without Compose):

```bash
./deploy-docker.sh
```

### Option 3: Local Development (No Docker)

Continue with your current setup:

```bash
# Terminal 1: Backend
go run ./cmd/webapi

# Terminal 2: Frontend
cd webui && npm run dev
```

## Troubleshooting

### "Cannot connect to Docker daemon"

-   Make sure Docker Desktop is running
-   You should see the Docker whale icon in your menu bar
-   Try restarting Docker Desktop if needed

### "docker-compose: command not found"

-   Use `./deploy-docker.sh` instead of `./deploy.sh`
-   Or install Docker Desktop which includes Docker Compose

### Port conflicts

If ports 80 or 3000 are already in use:

```bash
# Check what's using the ports
lsof -i :80
lsof -i :3000

# Stop conflicting services or modify the scripts to use different ports
```

## Current Status

Your application is working fine in development mode. Docker deployment is optional but provides:

-   Production-like environment
-   Easy deployment and scaling
-   Consistent environment across different machines
-   Single-command deployment

You can continue using the local development setup while Docker gets configured.
