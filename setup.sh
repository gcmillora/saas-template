#!/bin/bash
set -e

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${GREEN}Setting up SaaS Template...${NC}"

# Check Docker is running
if ! docker info > /dev/null 2>&1; then
  echo -e "${RED}Error: Docker is not running. Please start Docker and try again.${NC}"
  exit 1
fi

# Copy env files if they don't exist
if [ ! -f backend/.env.local ]; then
  echo -e "${YELLOW}Creating backend/.env.local from .env.example...${NC}"
  cp backend/.env.example backend/.env.local
fi

if [ ! -f backend/.env.docker ]; then
  echo -e "${YELLOW}Creating backend/.env.docker from .env.example...${NC}"
  cp backend/.env.example backend/.env.docker
  # Override DATABASE_URL for Docker networking
  sed -i.bak 's|postgresql://postgres:postgres@localhost:5432|postgresql://postgres:postgres@postgres:5432|' backend/.env.docker
  rm -f backend/.env.docker.bak
fi

if [ ! -f backend/.env.test ]; then
  echo -e "${YELLOW}Creating backend/.env.test from .env.example...${NC}"
  cp backend/.env.example backend/.env.test
fi

# Start Postgres
echo -e "${GREEN}Starting PostgreSQL...${NC}"
docker compose up -d postgres

# Wait for Postgres to be healthy
echo -e "${YELLOW}Waiting for PostgreSQL to be ready...${NC}"
until docker compose exec postgres pg_isready -U postgres -h localhost > /dev/null 2>&1; do
  sleep 1
done
echo -e "${GREEN}PostgreSQL is ready.${NC}"

# Start backend (runs migrations + Jet codegen via entrypoint)
echo -e "${GREEN}Starting backend...${NC}"
docker compose --profile saas-backend up -d go-backend

# Start frontend
echo -e "${GREEN}Starting frontend...${NC}"
docker compose --profile frontend up -d frontend

echo ""
echo -e "${GREEN}Setup complete!${NC}"
echo -e "  Frontend: ${YELLOW}http://localhost:8009${NC}"
echo -e "  Backend:  ${YELLOW}http://localhost:8008${NC}"
echo ""
echo -e "Run ${YELLOW}docker compose logs -f${NC} to see logs."
