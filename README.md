# SaaS Template

A full-stack SaaS template with a Go backend and React frontend, featuring OpenAPI code generation, type-safe database access, and Docker-based local development.

## Tech Stack

**Backend:**
- Go 1.24
- Chi router
- PostgreSQL
- Jet (type-safe SQL)
- OpenAPI/oapi-codegen
- Goose migrations

**Frontend:**
- React 19
- TypeScript
- Vite
- Redux Toolkit
- React Router
- Tailwind CSS
- shadcn/ui

## Project Structure

```
.
├── backend/                    # Go backend application
│   ├── cmd/                   # Application entrypoints
│   │   ├── console/          # CLI commands
│   │   └── webserver/        # HTTP server
│   ├── config/               # Configuration management
│   │   └── provider/         # Config providers
│   ├── db/                   # Database files
│   │   └── migrations/       # SQL migration files
│   ├── generated/            # Generated code (do not edit manually)
│   │   ├── db/              # Jet-generated database models
│   │   └── oapi/            # OpenAPI-generated server code
│   ├── internal/             # Private application code
│   │   ├── app/             # Core application logic
│   │   ├── console/         # Console command handlers
│   │   └── webserver/       # HTTP handlers and middleware
│   ├── scripts/              # Build and utility scripts
│   ├── commands.sh           # Development command runner
│   ├── openapi.yaml          # OpenAPI specification (private API)
│   ├── openapi-public.yaml   # OpenAPI specification (public API)
│   └── .env.local            # Local environment variables
│
├── frontend/                  # React frontend application
│   ├── public/               # Static assets
│   ├── src/                  # Source code
│   │   ├── assets/          # Images, fonts, etc.
│   │   ├── components/      # React components
│   │   │   └── ui/         # shadcn/ui components
│   │   ├── contexts/        # React contexts
│   │   ├── hooks/           # Custom React hooks
│   │   ├── lib/             # Utility libraries
│   │   ├── pages/           # Page components
│   │   ├── services/        # API services
│   │   │   └── api/        # Generated API client from OpenAPI
│   │   ├── store/           # Redux store configuration
│   │   └── utils/           # Utility functions
│   ├── package.json
│   └── vite.config.ts
│
└── docker-compose.yml         # Docker services configuration
```

## OpenAPI Code Generation

This project uses OpenAPI specifications to ensure type-safe communication between frontend and backend.

### How It Works

1. **Define API contracts** in `backend/openapi.yaml` and `backend/openapi-public.yaml`
2. **Generate backend code** - Creates Go server stubs and request/response types
3. **Generate frontend types** - Creates TypeScript types and RTK Query API client

### Regenerate Code

After modifying OpenAPI specs, regenerate code:

```bash
cd backend
./commands.sh openapi:codegen
```

This will:
- Generate Go server code in `backend/generated/oapi/`
- Generate TypeScript types in `frontend/src/services/api/`

## Docker Setup

### Services

The Docker Compose setup includes:

- **postgres** - PostgreSQL 17 database on port 5432
- **go-backend** - Go application with hot-reloading (port 8008)

### Running with Docker

Start the database:
```bash
docker compose up postgres -d
```

Start the backend (includes database):
```bash
docker compose --profile saas-backend up -d
```

Stop all services:
```bash
docker compose down
```

### Environment Files

The backend uses different environment files for different contexts:
- `.env.local` - Local development
- `.env.docker` - Docker container
- `.env.test` - Test environment

## Database Migrations

All database operations use [Goose](https://github.com/pressly/goose) for migrations and [Jet](https://github.com/go-jet/jet) for type-safe database access.

### Create a New Migration

```bash
cd backend
./commands.sh migration:codegen <migration_name>
```

This creates a new sequential migration file in `db/migrations/`.

### Run Migrations

Apply all pending migrations:
```bash
cd backend
./commands.sh migration:up
```

This will:
1. Run migrations on the local database
2. Regenerate type-safe database models using Jet
3. Run migrations on the test database

### Rollback Migrations

Rollback the last migration:
```bash
cd backend
./commands.sh migration:down
```

Rollback all migrations:
```bash
cd backend
./commands.sh migration:reset
```

### Check Migration Status

View which migrations have been applied:
```bash
cd backend
./commands.sh migration:status
```

## Getting Started

### Prerequisites

- Go 1.24+
- Node.js 18+ (or Bun)
- Docker and Docker Compose
- PostgreSQL (or use Docker)

### Backend Setup

1. Start the database:
   ```bash
   docker compose up postgres -d
   ```

2. Navigate to backend:
   ```bash
   cd backend
   ```

3. Install Go dependencies:
   ```bash
   go mod download
   ```

4. Run migrations:
   ```bash
   ./commands.sh migration:up
   ```

5. Start the server:
   ```bash
   ./commands.sh webserver
   ```

The backend will be available at `http://localhost:8008`

### Frontend Setup

1. Navigate to frontend:
   ```bash
   cd frontend
   ```

2. Install dependencies:
   ```bash
   npm install
   # or
   bun install
   ```

3. Start the development server:
   ```bash
   npm run dev
   # or
   bun dev
   ```

The frontend will be available at the URL shown in the terminal (typically `http://localhost:5173`)

## Available Commands

### Backend Commands

All backend commands are run via `commands.sh`:

```bash
cd backend

# Development
./commands.sh webserver              # Start web server
./commands.sh console {command}      # Run console command

# Code Quality
./commands.sh lint                   # Lint codebase
./commands.sh lint:fix              # Lint and auto-fix issues
./commands.sh format                # Format code
./commands.sh test                  # Run tests

# Code Generation
./commands.sh openapi:codegen       # Generate API code from OpenAPI specs

# Database Migrations
./commands.sh migration:codegen {name}  # Create new migration
./commands.sh migration:up             # Apply migrations
./commands.sh migration:down           # Rollback last migration
./commands.sh migration:reset          # Rollback all migrations
./commands.sh migration:status         # Check migration status
```

### Frontend Commands

```bash
cd frontend

npm run dev      # Start development server
npm run build    # Build for production
npm run lint     # Lint code
npm run preview  # Preview production build
```

## Development Workflow

1. **Design API endpoints** - Update `openapi.yaml` with new endpoints
2. **Generate code** - Run `./commands.sh openapi:codegen` to generate types
3. **Create migration** - Run `./commands.sh migration:codegen <name>` for schema changes
4. **Implement backend** - Write handlers in `internal/webserver/`
5. **Implement frontend** - Use generated API client in React components
6. **Test** - Run `./commands.sh test` and verify in browser
