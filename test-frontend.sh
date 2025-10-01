#!/bin/bash

echo "Building frontend container..."
sudo docker build -f Dockerfile.frontend -t wasa-frontend-test .

if [ $? -eq 0 ]; then
    echo "✅ Frontend build successful"
    echo "Testing frontend container..."
    
    # Run the container
    sudo docker run --rm -d --name wasa-frontend-test -p 8080:80 wasa-frontend-test
    
    if [ $? -eq 0 ]; then
        echo "✅ Frontend container started successfully"
        echo "Checking if frontend is responding..."
        sleep 3
        
        # Test if frontend is responding
        if curl -f http://localhost:8080 2>/dev/null | head -20; then
            echo "✅ Frontend is responding on port 8080"
        else
            echo "❌ Frontend not responding on port 8080"
        fi
        
        # Show container logs
        echo "--- Frontend container logs ---"
        sudo docker logs wasa-frontend-test
        
        # Stop the test container
        sudo docker stop wasa-frontend-test
    else
        echo "❌ Frontend container failed to start"
    fi
else
    echo "❌ Frontend build failed"
fi