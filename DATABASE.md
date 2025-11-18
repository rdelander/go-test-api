# Database Setup Guide

## Prerequisites

- Docker and Docker Compose installed
- golang-migrate CLI (optional, for manual migrations)

## Install golang-migrate (optional)

```bash
# macOS
brew install golang-migrate

# Or download binary from https://github.com/golang-migrate/migrate/releases
```

## Quick Start

### 1. Start Database

```bash
docker-compose up -d
```

This starts:

- **PostgreSQL** on `localhost:5432`
  - User: `gouser`
  - Password: `gopassword`
  - Database: `gotestdb`
- **pgAdmin** on `http://localhost:5050`
  - Email: `admin@example.com`
  - Password: `admin`

### 2. Run Migrations

Using the provided script:

```bash
./migrate.sh up
```

Or manually with golang-migrate:

```bash
migrate -database "postgresql://gouser:gopassword@localhost:5432/gotestdb?sslmode=disable" \
  -path migrations up
```

### 3. Start the API

```bash
go run cmd/server/main.go
```

## Migration Commands

```bash
# Run all migrations
./migrate.sh up

# Rollback last migration
./migrate.sh down

# Create new migration
./migrate.sh create add_user_avatar

# Check migration version
./migrate.sh version
```

## Database Access

### psql CLI

```bash
docker exec -it go-test-api-db psql -U gouser -d gotestdb
```

### pgAdmin Web UI

1. Open http://localhost:5050
2. Login with admin@example.com / admin
3. Add server:
   - Host: postgres (use Docker service name)
   - Port: 5432
   - Database: gotestdb
   - Username: gouser
   - Password: gopassword

## Stop Database

```bash
# Stop but keep data
docker-compose stop

# Stop and remove containers (keeps volumes/data)
docker-compose down

# Stop and remove everything including data
docker-compose down -v
```

## Troubleshooting

### Connection refused

Make sure PostgreSQL is running:

```bash
docker-compose ps
```

### Reset database

```bash
docker-compose down -v
docker-compose up -d
./migrate.sh up
```
