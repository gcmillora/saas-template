#!/bin/bash

# ANSI Color Codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

show_divider() {
  echo "----------------------------------------"
}

use_env() {
  echo "Using .env.${1} as ENV"
  set -o allexport
  source ".env.${1}"
  set +o allexport
}

check_app_env() {
  echo "Checking APP_ENV. Current value: '${APP_ENV:-not set}'"
  if [ -z "${APP_ENV}" ]; then
    echo "APP_ENV not set, defaulting to 'local'"
    use_env "local"
  else
    echo "Using existing APP_ENV: $APP_ENV"
  fi
}

# Function to display help message
display_help() {
  echo "Usage: $0 {command}"
  echo
  echo "Available commands:"
  echo -e "  ${YELLOW}webserver${NC}                 - Starts the web server using the .env.local configuration."
  echo -e "  ${YELLOW}console {command}${NC}          - Runs a console command."
  echo -e "  ${YELLOW}lint${NC}                      - Lints the codebase using golangci-lint."
  echo -e "  ${YELLOW}lint:fix${NC}                  - Lints the codebase and applies automated fixes."
  echo -e "  ${YELLOW}format${NC}                    - Formats the codebase using golangci-lint fmt."
  echo -e "  ${YELLOW}test${NC}                      - Runs the test suite."
  echo -e "  ${YELLOW}openapi:codegen${NC}           - Generates server and model code from openapi.yaml."
  echo -e "  ${YELLOW}migration:codegen {name}${NC}  - Creates a new empty SQL migration file."
  echo -e "  ${YELLOW}migration:up${NC}              - Runs all pending database migrations and regenerates the type-safe database models."
  echo -e "  ${YELLOW}migration:down${NC}            - Rolls back the last database migration and regenerates the models."
  echo -e "  ${YELLOW}migration:reset${NC}             - Rolls back all database migrations."
  echo -e "  ${YELLOW}migration:status${NC}          - Shows the status of all migrations."
}

# Main script logic
case "${1}" in
"webserver")
  check_app_env
  echo -e "${GREEN}Starting webserver...${NC}"
  go run ./cmd/webserver
  exit 0
  ;;
"console")
  shift
  echo -e "${GREEN}Running console command: $@${NC}"
  go run cmd/console/main.go "$@"
  ;;
"lint")
  echo -e "${GREEN}Linting codebase...${NC}"
  golangci-lint run
  ;;
"lint:fix")
  echo -e "${GREEN}Linting and fixing codebase...${NC}"
  golangci-lint run --fix
  ;;
"format")
  echo -e "${GREEN}Formatting codebase...${NC}"
  golangci-lint fmt
  ;;
"test")
  echo -e "${GREEN}Running tests...${NC}"
  go test ./...
  ;;
"openapi:codegen")
  echo -e "${GREEN}Generating OpenAPI code...${NC}"
  go tool oapi-codegen -config ./generated/oapi/codegen.yaml ./openapi.yaml

  echo -e "${GREEN}Generating OpenAPI public code...${NC}"
  go tool oapi-codegen -config ./generated/oapi/public/codegen.yaml ./openapi-public.yaml

  echo -e "${GREEN}Finished generating OpenAPI code...${NC}"
  echo -e "${GREEN}Generating frontend types with Orval..."
  cd ../frontend && bunx orval && cd ../backend
  echo -e "${GREEN}Finished generating frontend types"
  exit 0
  ;;

"migration:codegen")
  mkdir -p db/migrations
  # create a stub migration files
  # -s means to create a sequential migraiton instead of timestamped
  go tool goose -s -dir db/migrations create "${2}" sql
  exit 0
  ;;

"migration:status")
  use_env "docker"
  echo "DATABASE_URL: $DATABASE_URL"
  go tool goose -dir db/migrations postgres "$DATABASE_URL" status
  exit 0
  ;;

"migration:up")
  # run migrations for local db
  use_env "local"
  echo 'Running migrations from db/migrations'

  echo "DATABASE_URL: $DATABASE_URL"
  echo "DATABASE_NAME: $DATABASE_NAME"

  if [ -z "$DATABASE_URL" ]; then
    echo "Error: DATABASE_URL is not set"
    exit 1
  fi

  show_divider
  echo 'Migrating local DB'
  use_env "local"
  go tool goose -dir db/migrations postgres "$DATABASE_URL" up

  show_divider

  echo 'Autogenerating DB models'
  use_env "local"
  go tool jet -dsn="$DATABASE_URL" -schema=public -path=./generated/db

  show_divider

  echo 'Migrating local DB for tests'
  use_env "test"
  go tool goose -dir db/migrations postgres "$DATABASE_URL" up
  echo 'Finished running migrations'
  exit 0
  ;;

"migration:down")
  echo 'Migrating local DB'
  use_env "local"
  go tool goose -dir db/migrations postgres "$DATABASE_URL" down

  show_divider

  echo 'Autogenerating DB models'
  use_env "local"
  go tool jet -dsn="$DATABASE_URL" -schema=public -path=./generated/db

  show_divider

  echo 'Migrating local DB for tests'
  use_env "test"
  go tool goose -dir db/migrations postgres "$DATABASE_URL" down
  echo 'Finished running migrations'
  exit 0
  ;;

"migration:reset")
  echo 'Migrating local DB'
  use_env "local"
  go tool goose -dir db/migrations postgres "$DATABASE_URL" reset

  echo 'Migrating local DB for tests'
  use_env "test"
  go tool goose -dir db/migrations postgres "$DATABASE_URL" reset
  echo 'Finished running migrations'
  exit 0
  ;;

esac

