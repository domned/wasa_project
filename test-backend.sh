#!/bin/bash

echo "Building backend container..."
sudo docker build -f Dockerfile.backend -t wasa-backend-test .

if [ $? -eq 0 ]; then
    echo "✅ Backend build successful"
    echo "Testing backend container..."
    
    # Run the container and see if it starts without permission errors
    sudo docker run --rm -d --name wasa-backend-test -p 3001:3000 \
        -e CFG_DB_FILENAME=/data/app.db \
        -e PORT=3000 \
        wasa-backend-test
    
    if [ $? -eq 0 ]; then
        echo "✅ Backend container started successfully"
        echo "Checking if backend is responding..."
        sleep 3
        
        # Test if backend is responding
        if curl -f http://localhost:3001/liveness 2>/dev/null; then
            echo "✅ Backend is responding to health checks"
        else
            echo "❌ Backend not responding on port 3001"
        fi
        
        # Show container logs
        echo "--- Backend container logs ---"
        sudo docker logs wasa-backend-test
        
        # Stop the test container
        sudo docker stop wasa-backend-test
    else
        echo "❌ Backend container failed to start"
    fi
else
    echo "❌ Backend build failed"
fi