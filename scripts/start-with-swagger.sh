#!/bin/bash

# This script starts the Docker containers and runs swagger generation

echo "Starting Docker services..."
docker compose up -d

echo "Waiting for API service to be ready..."
RETRY_COUNT=0
MAX_RETRIES=30

while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
  if wget -q -O /dev/null http://localhost:8080/healthcheck; then
    echo "API is running successfully!"
    break
  else
    echo "API not ready yet... waiting (Attempt $((RETRY_COUNT+1))/$MAX_RETRIES)"
    RETRY_COUNT=$((RETRY_COUNT+1))
    sleep 2
  fi
done

if [ $RETRY_COUNT -eq $MAX_RETRIES ]; then
  echo "API failed to start properly after multiple attempts."
  exit 1
fi

echo "API is running. Generating Swagger documentation..."
make swagger

echo "Swagger documentation generated!"
echo "You can access the Swagger UI at: http://localhost:8080/swagger/index.html"
