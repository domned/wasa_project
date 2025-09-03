#!/bin/bash

# WASAText Stop Script

echo "🛑 Stopping WASAText services..."

# Stop containers
docker stop wasa-backend wasa-frontend 2>/dev/null || true

# Remove containers
docker rm wasa-backend wasa-frontend 2>/dev/null || true

# Optional: Remove network (uncomment if desired)
# docker network rm wasa-network 2>/dev/null || true

echo "✅ Services stopped successfully!"
echo ""
echo "💡 To also remove the Docker images:"
echo "   docker rmi wasa-backend wasa-frontend"
echo ""
echo "💡 To remove the database (WARNING: this will delete all data):"
echo "   rm -rf data/"
