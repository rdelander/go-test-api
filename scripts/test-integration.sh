#!/bin/bash
set -e

echo "Running integration tests..."

# Ensure database is running
echo "Checking if database is available..."
if ! docker ps | grep -q go-test-api-db; then
    echo "Starting database..."
    docker-compose up -d postgres
    echo "Waiting for database to be ready..."
    sleep 3
fi

# Run migrations
echo "Running migrations..."
./migrate.sh up

# Run integration tests
echo "Running tests..."
go test -tags=integration -v ./...

echo "Integration tests complete!"
