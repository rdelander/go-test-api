#!/bin/bash
set -e

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${GREEN}Starting E2E tests...${NC}"

# Configuration
NETWORK_NAME="e2e-test-network"
DB_CONTAINER="e2e-postgres"
API_CONTAINER="e2e-api"
IMAGE_TAG="go-test-api:e2e-test"
API_PORT="8080"

# Cleanup function
cleanup() {
    echo -e "${YELLOW}Cleaning up...${NC}"
    docker stop $API_CONTAINER $DB_CONTAINER 2>/dev/null || true
    docker rm $API_CONTAINER $DB_CONTAINER 2>/dev/null || true
    docker network rm $NETWORK_NAME 2>/dev/null || true
}

# Trap to ensure cleanup on exit
trap cleanup EXIT

# Clean up any existing resources
cleanup

# Build Docker image
echo -e "${GREEN}Building Docker image...${NC}"
docker build -t $IMAGE_TAG .

# Create network
echo -e "${GREEN}Creating network...${NC}"
docker network create $NETWORK_NAME

# Start PostgreSQL
echo -e "${GREEN}Starting PostgreSQL...${NC}"
docker run -d \
    --name $DB_CONTAINER \
    --network $NETWORK_NAME \
    -e POSTGRES_USER=gouser \
    -e POSTGRES_PASSWORD=gopassword \
    -e POSTGRES_DB=gotestdb \
    postgres:17-alpine

# Wait for PostgreSQL to be ready
echo -e "${YELLOW}Waiting for PostgreSQL...${NC}"
for i in {1..30}; do
    if docker exec $DB_CONTAINER pg_isready -U gouser > /dev/null 2>&1; then
        echo -e "${GREEN}PostgreSQL is ready${NC}"
        break
    fi
    if [ $i -eq 30 ]; then
        echo -e "${RED}PostgreSQL failed to start${NC}"
        docker logs $DB_CONTAINER
        exit 1
    fi
    sleep 1
done

# Start API server
echo -e "${GREEN}Starting API server...${NC}"
docker run -d \
    --name $API_CONTAINER \
    --network $NETWORK_NAME \
    -p $API_PORT:8080 \
    -e ENV=production \
    -e PORT=8080 \
    -e DB_HOST=$DB_CONTAINER \
    -e DB_PORT=5432 \
    -e DB_USER=gouser \
    -e DB_PASSWORD=gopassword \
    -e DB_NAME=gotestdb \
    -e DB_SSLMODE=disable \
    -e JWT_SECRET=test-secret-for-e2e \
    $IMAGE_TAG

# Give API a moment to start
sleep 2

# Show API logs in background
echo -e "${YELLOW}API logs:${NC}"
docker logs -f $API_CONTAINER &
LOGS_PID=$!

# Wait a bit for logs to start showing
sleep 1

# Run E2E tests
echo -e "${GREEN}Running E2E tests...${NC}"
export API_BASE_URL="http://localhost:$API_PORT"
go test -tags=e2e -v ./test/e2e/...

# Stop background logs
kill $LOGS_PID 2>/dev/null || true

echo -e "${GREEN}E2E tests complete!${NC}"
