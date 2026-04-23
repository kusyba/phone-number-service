#!/bin/sh

echo "Waiting for PostgreSQL..."
until pg_isready -h ${DB_HOST:-postgres} -p ${DB_PORT:-5432} -U ${DB_USER:-postgres}; do
  sleep 1
done

echo "Running migrations..."
goose -dir migrations postgres "host=${DB_HOST:-postgres} port=${DB_PORT:-5432} user=${DB_USER:-postgres} password=${DB_PASSWORD:-postgres} dbname=${DB_NAME:-phoneservice} sslmode=disable" up

echo "Generating sqlc code..."
sqlc generate

echo "Building application..."
go build -o /app/server ./cmd/server

echo "Starting application..."
/app/server
