## Project Template Specification

### 1. Introduction

This document outlines the specification for a Go backend project template. The template provides a solid foundation for building multi-tenant SaaS applications with built-in authentication, session management, and a clean architecture.

### 2. Philosophy

The template follows Clean Architecture principles, separating concerns into distinct layers:
- **Presentation Layer**: HTTP handlers and middleware
- **Application Layer**: Business logic orchestration (app_service)
- **Domain Layer**: Core business rules (domain_service)
- **Data Layer**: Repository pattern with type-safe SQL
- **Infrastructure**: Configuration providers and utilities

The template is configuration-driven with lazy-loaded providers and follows a schema-first API approach using OpenAPI specifications.

### 3. Folder Structure

```
/
├── cmd/
│   ├── console/              # CLI application entry point
│   │   └── main.go
│   └── webserver/            # HTTP server entry point
│       └── main.go
├── config/
│   ├── app.go                # Central App struct with lazy-loaded providers
│   └── provider/
│       ├── cache_provider.go      # In-memory cache
│       ├── database_provider.go   # PostgreSQL connection
│       ├── env_provider.go        # Environment variable loader
│       ├── logger_provider.go     # Structured logging (slog)
│       ├── session_provider.go    # Cookie-based sessions
│       ├── supabase_provider.go   # Supabase client for auth
│       └── validation_provider.go # Request validation
├── db/
│   └── migrations/           # SQL migration files (goose)
│       └── 00001_create_users_table.sql
├── generated/
│   ├── db/                   # Jet-generated type-safe models
│   └── oapi/                 # oapi-codegen generated code
│       ├── codegen.yaml
│       └── generated.go
├── internal/
│   ├── app/                  # Business logic layer
│   │   ├── app_service/      # Use case orchestration (e.g., user/)
│   │   ├── domain_service/   # Domain logic and mapping
│   │   ├── dto/              # Data Transfer Objects
│   │   ├── errors/           # Custom error types
│   │   ├── mutation/         # Write operations (Create/Update/Delete)
│   │   ├── repository/       # Read operations and queries
│   │   │   ├── pagination.go
│   │   │   ├── tenant_repository.go
│   │   │   └── user_repository.go
│   │   └── util_service/     # Shared utilities
│   ├── console/              # CLI command implementations
│   │   └── console.go
│   └── webserver/
│       ├── handler/          # HTTP request handlers
│       │   └── handler.go
│       ├── middleware/       # HTTP middleware
│       │   └── (session, logging, auth, etc.)
│       └── webserver.go
├── .env.local                # Local development config
├── .env.test                 # Test environment config
├── .env.docker               # Docker environment config
├── commands.sh               # Development workflow scripts
├── go.mod
├── go.sum
├── openapi.yaml              # Main API specification
├── openapi-public.yaml       # Public API specification (optional)
└── README.md
```

### 4. API Implementation

- **Framework**: Chi router for HTTP routing and middleware composition
- **Schema-First**: API defined in `openapi.yaml`, code generated via `oapi-codegen`
- **Handlers**: Located in `internal/webserver/handler`, orchestrate app_services and return responses
- **Middleware**: Located in `internal/webserver/middleware` for:
  - Session management
  - Authentication (Supabase integration)
  - Logging
  - Request validation
  - Error handling

### 5. Architecture Layers

#### Repository Layer (`internal/app/repository/`)
- Read-only database queries using Jet SQL builder
- Type-safe query construction
- Pagination utilities
- Multi-tenancy support (tenant-scoped queries)
- Example: `user_repository.go` shows GetUserByID, GetUserByAuthID, etc.

#### Mutation Layer (`internal/app/mutation/`)
- Write operations (Create, Update, Delete)
- Database transactions
- Timestamp management
- Example: `user_mutation.go` shows CreateUser, DeleteUser

#### App Service Layer (`internal/app/app_service/`)
- Business logic orchestration
- Combines repository, mutation, and domain services
- Session and auth context handling
- Example: `user/get_user.go` demonstrates retrieving authenticated user data

#### Domain Service Layer (`internal/app/domain_service/`)
- Domain-specific business rules
- Response mapping (DTO to API models)
- Business validations

#### DTOs (`internal/app/dto/`)
- Custom data structures for internal use
- Bridge between database models and API responses

### 6. Database

- **Migrations**: Plain SQL files in `db/migrations/`, applied with Goose
- **Type-Safe Queries**: Jet generates models from schema in `generated/db/`
- **Multi-Tenancy**: Built-in tenant_id support in repositories
- **Connection**: PostgreSQL via `lib/pq` driver

### 7. Configuration & Providers

The `config/App` struct uses lazy-loaded providers with getter methods:
- `app.DB()` - Database connection
- `app.EnvVars()` - Environment variables
- `app.Logger()` - Structured logger (slog)
- `app.Cache()` - In-memory cache
- `app.Session()` - Cookie store
- `app.Supabase()` - Auth client

**Environment Files**:
- `.env.local` - Local development
- `.env.test` - Test environment
- `.env.docker` - Docker environment

### 8. Key Dependencies

**Core Framework**:
- `github.com/go-chi/chi/v5` - HTTP router and middleware
- `github.com/oapi-codegen/runtime` - OpenAPI runtime utilities

**Database**:
- `github.com/lib/pq` - PostgreSQL driver
- `github.com/go-jet/jet/v2` - Type-safe SQL builder
- `github.com/pressly/goose/v3` - Database migrations

**Authentication & Sessions**:
- `github.com/supabase-community/supabase-go` - Supabase client
- `github.com/gorilla/sessions` - Cookie-based sessions

**Utilities**:
- `github.com/google/uuid` - UUID generation
- `github.com/gookit/validate` - Request validation
- `github.com/patrickmn/go-cache` - In-memory cache
- `github.com/joho/godotenv` - Environment variable loading

**Development Tools** (via `go tool`):
- `github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen` - OpenAPI code generation
- `github.com/go-jet/jet/v2/cmd/jet` - Database model generation
- `github.com/pressly/goose/v3/cmd/goose` - Migration CLI

### 9. Development Workflow (`commands.sh`)

The `commands.sh` script provides a unified interface for all development tasks:

**Server & Console**:
- `./commands.sh webserver` - Start HTTP server (auto-loads `.env.local`)
- `./commands.sh console {command}` - Run CLI commands

**Code Quality**:
- `./commands.sh lint` - Run golangci-lint
- `./commands.sh lint:fix` - Auto-fix linting issues
- `./commands.sh format` - Format code
- `./commands.sh test` - Run test suite

**Code Generation**:
- `./commands.sh openapi:codegen` - Generate server code from OpenAPI specs
  - Generates backend Go code
  - Generates frontend TypeScript types (RTK Query)
  - Supports both `openapi.yaml` and `openapi-public.yaml`

**Database Migrations**:
- `./commands.sh migration:codegen {name}` - Create new SQL migration file
- `./commands.sh migration:up` - Apply all pending migrations + regenerate Jet models
- `./commands.sh migration:down` - Rollback last migration + regenerate models
- `./commands.sh migration:reset` - Rollback all migrations
- `./commands.sh migration:status` - View migration status

**Environment Management**:
The script automatically detects and uses the appropriate `.env` file based on `APP_ENV` or defaults to `.env.local`.

### 10. Getting Started

To create a new project from this template:

1. **Setup Project**
   ```bash
   # Copy template to new directory
   cp -r backend my-new-project
   cd my-new-project
   
   # Update go.mod with your module name
   go mod edit -module github.com/your-org/your-project
   ```

2. **Configure Environment**
   ```bash
   # Copy and edit .env files
   cp .env.local.example .env.local  # Configure DATABASE_URL, SUPABASE_URL, etc.
   cp .env.test.example .env.test
   ```

3. **Define Your API**
   - Edit `openapi.yaml` to define your endpoints
   - Run `./commands.sh openapi:codegen` to generate code

4. **Setup Database**
   ```bash
   # Create migrations
   ./commands.sh migration:codegen create_initial_tables
   
   # Edit the generated SQL file in db/migrations/
   # Then apply migrations
   ./commands.sh migration:up
   ```

5. **Implement Business Logic**
   - **Repositories** (read operations): `internal/app/repository/`
   - **Mutations** (write operations): `internal/app/mutation/`
   - **App Services** (orchestration): `internal/app/app_service/`
   - **Handlers** (HTTP layer): `internal/webserver/handler/`
   
   Refer to the user service example for patterns.

6. **Run & Test**
   ```bash
   ./commands.sh webserver    # Start server
   ./commands.sh test         # Run tests
   ./commands.sh lint         # Check code quality
   ```

### 11. Reference Implementation

The template includes a complete **User Service** example demonstrating all architectural layers:

- `internal/app/repository/user_repository.go` - Query patterns with Jet
- `internal/app/mutation/user_mutation.go` - Create/Delete operations
- `internal/app/app_service/user/get_user.go` - Session handling & orchestration
- `openapi.yaml` - API endpoint definition (`/user`)

Follow this pattern when implementing new features.

### 12. Best Practices

**Repository Pattern**:
- Keep repositories read-only
- Use Jet for type-safe queries
- Always include tenant_id in WHERE clauses for multi-tenancy

**Mutation Pattern**:
- Handle timestamps (created_at, updated_at) automatically
- Use transactions for multi-table operations
- Return inserted/updated records

**App Service Pattern**:
- Extract session data from context
- Call repositories for reads, mutations for writes
- Map database models to API responses
- Handle business logic validation

**Error Handling**:
- Use custom errors from `internal/app/errors/`
- Return appropriate HTTP status codes
- Log errors with structured logging via `app.Logger()`
