#!/bin/bash

# Migration helper script
# Usage: ./migrate.sh up|down|create|version

set -e

DB_URL="postgresql://gouser:gopassword@localhost:5432/gotestdb?sslmode=disable"

case "$1" in
  up)
    echo "Running migrations..."
    migrate -database "${DB_URL}" -path migrations up
    ;;
  down)
    echo "Rolling back one migration..."
    migrate -database "${DB_URL}" -path migrations down 1
    ;;
  create)
    if [ -z "$2" ]; then
      echo "Usage: ./migrate.sh create <migration_name>"
      exit 1
    fi
    echo "Creating migration: $2"
    migrate create -ext sql -dir migrations -seq "$2"
    ;;
  version)
    migrate -database "${DB_URL}" -path migrations version
    ;;
  *)
    echo "Usage: ./migrate.sh {up|down|create|version}"
    exit 1
    ;;
esac
