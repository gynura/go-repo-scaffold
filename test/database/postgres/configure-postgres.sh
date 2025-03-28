#!/usr/bin/env bash
set -e

CONTAINER_NAME="postgres_10_12"

docker cp create-databases.sql "${CONTAINER_NAME}":/
docker exec -e PGPASSWORD=postgres -it "${CONTAINER_NAME}" psql -U postgres -a -f /create-databases.sql
