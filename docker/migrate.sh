#!/bin/sh

set -e

until nc -z postgres 5432
do
    echo "Waiting for PostgreSQL..."
    sleep 2
done

echo "Running migrations..."

exec /bin/goose \
    -dir /app/migrations \
    postgres \
    "$DATABASE_DSN" \
    up