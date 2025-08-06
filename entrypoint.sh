#!/bin/bash
set -e

# Wait for PostgreSQL to be ready
echo "Waiting for PostgreSQL..."
/wait-for-it.sh db:5432 --timeout=60 --strict -- echo "PostgreSQL is up"

# Run database migrations
echo "Running migrations..."
migrate -path /app/migrations -database "$DATABASE_URL" up

# Run the application
echo "Starting application..."
exec "$@"
