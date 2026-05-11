#!/usr/bin/env bash

set -e

DB_HOST="localhost"
DB_PORT="5432"
DB_URL="postgres://postgres:pass@localhost:5432/user_profiles?sslmode=disable"

echo "Starting Postgres..."

docker compose up -d postgres

echo "Waiting for Postgres to be ready..."

# Проверка готовности БД
until docker exec postgres_profiles pg_isready -U postgres >/dev/null 2>&1; do
  sleep 1
done

echo "Postgres is ready ✔"

echo "1. Dropping all tables..."
migrate -path ./migrations -database "$DB_URL" drop -f

echo "2. Running migrations..."
go run cmd/migrate/auto.go -up

echo "3. Creating roles..."
go run ./cmd/seeds -mode=roles -roles="admin,moderator,user"

echo "4. Creating users..."

go run ./cmd/seeds -mode=users -role=admin -count=3
go run ./cmd/seeds -mode=users -role=moderator -count=2
go run ./cmd/seeds -mode=users -role=user

echo "5. Creating posts..."

go run ./cmd/seeds -mode=posts -role=user -count=15
go run ./cmd/seeds -mode=posts -role=moderator -count=5
go run ./cmd/seeds -mode=posts -role=admin -count=1

echo "Done ✔"