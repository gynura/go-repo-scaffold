#!/usr/bin/env bash
set -e

CONTAINER_NAME="postgres_10_12"

echo "Waiting for Postgres"
until docker exec "${CONTAINER_NAME}" pg_isready -U postgres > /dev/null 2>&1; do
  echo "Postgres is not ready yet. Retrying in 2 seconds..."
  sleep 2
done

docker cp create-databases.sql "${CONTAINER_NAME}":/
docker exec -e PGPASSWORD=postgres -it "${CONTAINER_NAME}" psql -U postgres -a -f /create-databases.sql

echo "Migrations done!"
