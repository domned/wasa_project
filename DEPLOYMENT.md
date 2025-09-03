# WASAText Docker Deployment

This directory contains Docker configurations for deploying the WASAText chat application.

## Quick Start

### Option 1: Docker Compose (Recommended)

1. **Clone and navigate to the project:**

    ```bash
    git clone <repository-url>
    cd wasa_project
    ```

2. **Deploy with one command:**
    ```bash
    ./deploy.sh
    ```

### Option 2: Docker Only

If you don't have Docker Compose installed:

1. **Deploy with Docker:**

    ```bash
    ./deploy-docker.sh
    ```

2. **Stop services:**

    ```bash
    ./stop.sh
    ```

3. **Access the application:**
    - Frontend: http://localhost
    - Backend API: http://localhost:3000

## Manual Deployment

If you prefer to deploy manually:

```bash
# Build the containers
docker-compose build

# Start the services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop the services
docker-compose down
```

## Architecture

-   **Backend**: Go application running on port 3000
-   **Frontend**: Vue.js application served by Nginx on port 80
-   **Database**: SQLite database persisted in `./data` volume
-   **Reverse Proxy**: Nginx proxies `/api/*` requests to the backend

## Services

### Backend (wasa-backend)

-   Built from `Dockerfile.backend`
-   Runs the Go web API server
-   Exposes port 3000
-   Database stored in `/data/app.db`

### Frontend (wasa-frontend)

-   Built from `Dockerfile.frontend`
-   Serves the Vue.js SPA via Nginx
-   Exposes port 80
-   Proxies API calls to backend

## Configuration

### Environment Variables

-   `DATABASE_PATH`: Path to SQLite database (default: `/data/app.db`)
-   `PORT`: Backend server port (default: `3000`)

### Volumes

-   `./data:/data`: Persists the SQLite database

## Health Checks

Both services include health checks:

-   **Backend**: Checks `/liveness` endpoint
-   **Frontend**: Checks Nginx availability

## Development vs Production

### Development

-   Frontend connects directly to `http://localhost:3000`
-   Hot reload enabled
-   Separate backend/frontend processes

### Production (Docker)

-   Frontend proxies API calls through Nginx (`/api/*`)
-   Single-command deployment
-   Optimized builds with multi-stage Dockerfiles

## Troubleshooting

### View logs

```bash
docker-compose logs backend
docker-compose logs frontend
```

### Restart services

```bash
docker-compose restart
```

### Rebuild containers

```bash
docker-compose down
docker-compose build --no-cache
docker-compose up -d
```

### Database issues

The database is persisted in the `./data` directory. To reset:

```bash
docker-compose down
rm -rf data
docker-compose up -d
```

## Production Considerations

For production deployment:

1. **Use environment-specific configs:**

    ```bash
    cp .env.production .env
    ```

2. **Enable HTTPS with reverse proxy (recommended):**

    - Use Traefik, Caddy, or nginx-proxy
    - Add SSL certificates
    - Update CORS settings if needed

3. **Database backups:**

    ```bash
    # Backup
    cp data/app.db backup/app.db.$(date +%Y%m%d)

    # Restore
    cp backup/app.db.20231201 data/app.db
    docker-compose restart backend
    ```

4. **Monitor logs:**
    ```bash
    docker-compose logs -f --tail=100
    ```

## Scaling

To scale for higher load:

-   Use external database (PostgreSQL)
-   Add load balancer for multiple backend instances
-   Implement Redis for session storage
-   Use CDN for static assets
