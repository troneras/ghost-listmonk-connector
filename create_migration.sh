#!/bin/bash

# Check if a migration name was provided
if [ $# -eq 0 ]; then
    echo "No migration name provided. Usage: ./create_migration.sh <migration_name>"
    exit 1
fi

# Create the migration
migrate create -ext sql -dir database/migrations -seq $1

echo "Migration created successfully."
