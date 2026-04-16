#!/bin/sh

# Wait for database to be ready
echo "Waiting for postgres to be ready..."
until pg_isready -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER"; do
  sleep 2
done

echo "Running migrations..."
# Using goose for migrations (as it's already in go.mod)
# Assuming migrations are in /app/migrations and binaries are not available, 
# we can use go run or build a small helper. 
# Better to build goose in a separate stage or use its binary.
# Given time constraints, I'll use a simple go run approach if possible, 
# or assume the user has goose installed in the image.

# Actually, let's just use the calendar binary if it can run migrations, 
# but the task says "one-shot script". 
# I'll use a specialized migration image like 'pressly/goose'.

# Since I don't want to rely on external tools not in the image, 
# I'll create a small migration helper in Go.

goose -dir /app/migrations postgres "host=$DB_HOST port=$DB_PORT user=$DB_USER password=$DB_PASSWORD dbname=$DB_NAME sslmode=disable" up
