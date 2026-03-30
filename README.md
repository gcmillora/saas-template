# SaaS Template

A production-ready SaaS starter with Go, React, and PostgreSQL. Clone it, customize it, ship your app.

**Backend:** Go 1.25 · Chi router · Jet ORM · oapi-codegen · Goose migrations
**Frontend:** React 19 · Vite · TanStack React Query · shadcn/ui · Tailwind CSS
**Infra:** PostgreSQL · Docker Compose · Nginx reverse proxy

## What's Included

- **Authentication** — Email/password signup and signin with bcrypt hashing
- **Password reset** — Secure token-based flow with email delivery (via Resend)
- **Session management** — HTTP-only secure cookies with configurable expiry
- **Audit logging** — All auth events logged with IP tracking and metadata
- **Three API tiers** — Public, authenticated, and admin with separate rate limits
- **Role-based access** — Admin and user roles with middleware enforcement
- **OpenAPI-first** — Type-safe codegen for both backend (Go) and frontend (TypeScript)
- **Type-safe database** — Jet ORM generates query builders from your schema
- **Production Docker** — Multi-stage builds, Nginx reverse proxy, security headers
- **CI pipeline** — GitHub Actions for lint, test, and build

## Getting Started

**Prerequisites:** Docker and Docker Compose

```bash
git clone https://github.com/your-org/saas-template.git
cd saas-template
./setup.sh
```

That's it. Open [http://localhost:8009](http://localhost:8009).

- Frontend runs on port **8009**
- Backend runs on port **8008**
- PostgreSQL runs on port **5432**

## Project Structure

```
├── backend/                 # Go backend
│   ├── cmd/                # Entry points (webserver, console)
│   ├── config/             # App config and dependency injection
│   ├── db/migrations/      # SQL migrations (Goose)
│   ├── generated/          # Auto-generated code (do not edit)
│   ├── internal/           # App logic, handlers, middleware
│   ├── openapi.yaml        # Authenticated API spec
│   ├── openapi-public.yaml # Public API spec
│   ├── openapi-admin.yaml  # Admin API spec
│   └── commands.sh         # Dev command runner
├── frontend/               # React frontend
│   ├── src/
│   │   ├── components/     # UI components (shadcn/ui)
│   │   ├── contexts/       # Auth context
│   │   ├── pages/          # Page components
│   │   └── services/api/   # Generated API hooks (Orval)
│   └── orval.config.ts     # API codegen config
├── nginx/                  # Production reverse proxy config
├── docker-compose.yml      # Development services
├── docker-compose.prod.yml # Production services
└── CLAUDE.md               # Codebase guide for AI-assisted development
```

## Development Commands

```bash
# Backend (from backend/)
bash commands.sh webserver              # Start dev server (hot reload)
bash commands.sh openapi:codegen        # Regenerate all API types + frontend hooks
bash commands.sh migration:codegen name # Create new migration
bash commands.sh migration:up           # Run migrations + regenerate DB models
bash commands.sh lint                   # Lint
bash commands.sh test                   # Test

# Frontend (from frontend/)
bun run dev                             # Start dev server
bun run build                           # Type check + build
bunx orval                              # Regenerate API hooks
```

## Make It Yours

After cloning, update these to match your project:

1. **Module name** — `backend/go.mod` (change `saas-template` to your project name, then update all Go imports)
2. **App name** — `frontend/src/components/app-sidebar.tsx` and `frontend/index.html`
3. **Docker images** — `docker-compose.prod.yml` (change `ghcr.io/your-org/saas-*`)
4. **Container names** — `docker-compose.yml` (change `saas-backend`, `saas-frontend`, etc.)
5. **Database name** — `docker-compose.yml` (change `POSTGRES_DB` and connection strings)
6. **Environment** — `backend/.env.example` (update secrets, URLs, API keys)

## AI-Assisted Development

This project includes a `CLAUDE.md` file that documents the codebase architecture, naming conventions, and patterns. It works with AI coding tools like Claude Code to provide context about how the project is structured.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

[MIT](LICENSE)
