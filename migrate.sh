#!/bin/bash

# Migration script for postgres-api
# Usage: ./migrate.sh <command> [args...]
# Examples:
#   ./migrate.sh init
#   ./migrate.sh create add_users_table
#   ./migrate.sh up
#   ./migrate.sh down
#   ./migrate.sh status

# Check if at least one argument is provided
if [ $# -eq 0 ]; then
    echo "Usage: ./migrate.sh <command> [args...]"
    echo ""
    echo "Available commands:"
    echo "  init              - Initialize the migration table"
    echo "  create <name>     - Create a new migration (use underscores for spaces)"
    echo "  up                - Run all pending migrations"
    echo "  down              - Rollback the last migration"
    echo "  status            - Check migration status"
    echo ""
    echo "Examples:"
    echo "  ./migrate.sh init"
    echo "  ./migrate.sh create add_users_table"
    echo "  ./migrate.sh up"
    echo "  ./migrate.sh down"
    echo "  ./migrate.sh status"
    exit 1
fi

# Run the migration command with all arguments
go run cmd/migrate/main.go migrate "$@"