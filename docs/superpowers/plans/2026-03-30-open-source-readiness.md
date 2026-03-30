# Open-Source Readiness Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Prepare saas-template for public GitHub release with proper licensing, documentation, and one-command setup.

**Architecture:** No code changes. This is documentation, config, and cleanup only. Add LICENSE, CONTRIBUTING.md, setup.sh, rewrite README, gitignore internal docs, remove misleading OpenAPI server block.

**Tech Stack:** Bash, Markdown, Git

**Spec:** `docs/superpowers/specs/2026-03-30-open-source-readiness-design.md`

---

## File Structure

### New Files
- `LICENSE` — MIT license
- `CONTRIBUTING.md` — Contributor guidelines
- `setup.sh` — One-command setup script

### Modified Files
- `.gitignore` — Add `docs/superpowers/`
- `README.md` — Full rewrite for external audience
- `backend/openapi.yaml` — Remove `servers` block

### Deleted
- `docs/superpowers/` — Internal design specs and plans (after gitignore)

---

### Task 1: Add LICENSE and CONTRIBUTING.md

**Files:**
- Create: `LICENSE`
- Create: `CONTRIBUTING.md`

- [ ] **Step 1: Create MIT LICENSE file**

Create `LICENSE`:

```
MIT License

Copyright (c) 2026 Greg

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

- [ ] **Step 2: Create CONTRIBUTING.md**

Create `CONTRIBUTING.md`:

```markdown
# Contributing

Thanks for your interest in contributing to this project!

## Getting Started

1. Fork the repo
2. Clone your fork
3. Run `./setup.sh` to get the dev environment running
4. Make your changes

## Bug Reports

Open an issue with:
- What you expected to happen
- What actually happened
- Steps to reproduce

## Feature Requests

Open an issue to discuss before building. This keeps scope manageable and avoids wasted effort.

## Pull Requests

- Bug fixes are always welcome
- Keep changes small and focused
- Follow existing patterns in the codebase
- Make sure `bash commands.sh lint` and `bash commands.sh test` pass in the backend
- Make sure `bun run build` passes in the frontend

I review PRs when I can — please be patient.

## Development Commands

```bash
# One-command setup
./setup.sh

# Backend (from backend/)
bash commands.sh webserver          # Start dev server
bash commands.sh lint               # Lint
bash commands.sh test               # Test
bash commands.sh openapi:codegen    # Regenerate API types

# Frontend (from frontend/)
bun run dev                         # Start dev server
bun run build                       # Type check + build
bun run lint                        # Lint
```
```

- [ ] **Step 3: Commit**

```bash
git add LICENSE CONTRIBUTING.md
git commit -m "chore: add MIT license and contributing guide"
```

---

### Task 2: Create setup.sh

**Files:**
- Create: `setup.sh`

- [ ] **Step 1: Create the setup script**

Create `setup.sh`:

```bash
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
```

- [ ] **Step 2: Make it executable**

```bash
chmod +x setup.sh
```

- [ ] **Step 3: Commit**

```bash
git add setup.sh
git commit -m "chore: add one-command setup script"
```

---

### Task 3: Gitignore internal docs and clean up OpenAPI

**Files:**
- Modify: `.gitignore`
- Modify: `backend/openapi.yaml`
- Delete: `docs/superpowers/`

- [ ] **Step 1: Add docs/superpowers/ to .gitignore**

In `.gitignore`, add at the top (before "# OS files"):

```
# Internal docs
docs/superpowers/
```

- [ ] **Step 2: Remove the docs/superpowers/ directory from git tracking**

```bash
git rm -r --cached docs/superpowers/
```

Note: This removes the files from git tracking but leaves them on disk (since they're now gitignored). Then delete the actual directory:

```bash
rm -rf docs/superpowers/
```

- [ ] **Step 3: Remove misleading servers block from openapi.yaml**

In `backend/openapi.yaml`, remove lines 5-6:

```yaml
servers:
  - url: "http://localhost:4000"
```

So the file starts with:

```yaml
openapi: 3.0.0
info:
  title: "SaaS Template API"
  version: "1.0.0"
paths:
```

- [ ] **Step 4: Commit**

```bash
git add .gitignore backend/openapi.yaml
git commit -m "chore: gitignore internal docs, remove misleading OpenAPI server URL"
```

---

### Task 4: Rewrite README.md

**Files:**
- Modify: `README.md`

- [ ] **Step 1: Replace entire README.md**

Replace the entire content of `README.md` with:

```markdown
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
```

- [ ] **Step 2: Commit**

```bash
git add README.md
git commit -m "docs: rewrite README for open-source release"
```

---

### Task 5: Final check

- [ ] **Step 1: Verify no sensitive files are tracked**

```bash
git status
git ls-files | grep -i '\.env' | grep -v example
```

Expected: No `.env.local`, `.env.docker`, `.env.test` files in output. Only `.env.example` files.

- [ ] **Step 2: Verify docs/superpowers is gone from git**

```bash
git ls-files | grep superpowers
```

Expected: No output.

- [ ] **Step 3: Verify all new files exist**

```bash
ls -la LICENSE CONTRIBUTING.md setup.sh
```

Expected: All three files exist, `setup.sh` has execute permission.

- [ ] **Step 4: Review the commit log**

```bash
git log --oneline -5
```

Expected: Clean commit history with descriptive messages.
