# Open-Source the SaaS Template

**Date:** 2026-03-30
**Status:** Approved

## Overview

Prepare saas-template for public GitHub release. Target audience: solo developers and beginners who want a production-ready Go + React SaaS starter. Light maintenance model — bug fixes and questions, not actively seeking contributors.

## Decisions

- **License:** MIT
- **Git history:** No rewrite. People clone from HEAD.
- **Internal docs:** Gitignore `docs/superpowers/`, delete the directory. Not useful to template users.
- **CLAUDE.md:** Keep. Useful for template users and their AI tools.
- **No changelog, releases, demo site, or badges wall.** Keep it simple.

---

## New Files

### `LICENSE`
MIT license. Standard for templates — low friction, maximum adoption.

### `CONTRIBUTING.md`
Short and honest. Contents:
- How to run locally (points to setup.sh)
- Bug reports: open an issue with reproduction steps
- Feature requests: open an issue first to discuss before building
- PRs: bug fixes welcome, keep scope small, follow existing patterns
- Set expectations: "I review PRs when I can — please be patient"

### `setup.sh`
One-command setup script. Behavior:
1. Check Docker is running (exit with helpful message if not)
2. Copy `backend/.env.example` to `backend/.env.local`, `backend/.env.docker`, `backend/.env.test` if they don't already exist
3. Run `docker compose up -d postgres` and wait for healthy
4. Run `docker compose up -d go-backend` (entrypoint handles migrations + Jet codegen)
5. Run `docker compose up -d frontend`
6. Print "Ready at http://localhost:8009"

Make executable (`chmod +x`).

---

## Modified Files

### `.gitignore`
Add `docs/superpowers/` to prevent internal specs/plans from being tracked.

### `README.md`
Full rewrite for external audience. Structure:

1. **One-line description** — "A production-ready SaaS starter template with Go, React, and PostgreSQL"
2. **Tech stack** — Go 1.25, Chi, Jet ORM, React 19, Vite, shadcn/ui, PostgreSQL, Docker
3. **What's included** — Feature list:
   - Email/password auth with bcrypt
   - Password reset flow with secure tokens
   - Session-based auth with secure cookies
   - Audit logging (all auth events with IP tracking)
   - Three API tiers (public, authenticated, admin)
   - Role-based access control (admin/user)
   - Rate limiting per route group
   - OpenAPI-first with type-safe codegen (backend + frontend)
   - Type-safe database queries (Jet ORM)
   - Production Docker setup with Nginx reverse proxy
   - CI pipeline (lint + test + build)
4. **Getting started** — `git clone` → `./setup.sh` → open http://localhost:8009
5. **Project structure** — Brief directory overview (backend/, frontend/, configs/, nginx/)
6. **Make it yours** — Checklist:
   - Change module name in `backend/go.mod`
   - Update app name in `frontend/src/components/app-sidebar.tsx` and `frontend/index.html`
   - Update image refs in `docker-compose.prod.yml`
   - Update `docker-compose.yml` container names and database name
   - Replace `backend/.env.example` values with your own
7. **AI-assisted development** — "This project includes a CLAUDE.md that documents the codebase architecture and patterns for AI-assisted development."
8. **Contributing** — Link to CONTRIBUTING.md
9. **License** — MIT

### `backend/openapi.yaml`
Remove the `servers` block (`url: "http://localhost:4000"`). It's misleading — the app runs on 8008 and this field isn't used by the codegen.

---

## Deleted

### `docs/superpowers/`
Delete entire directory. Contains internal design specs and implementation plans that aren't useful to template users.

---

## Not Changed

- **Git history** — No squash or rewrite
- **CI pipeline** — Already sufficient (lint + test + build)
- **CLAUDE.md** — Keep as-is
- **No CODE_OF_CONDUCT** — Implies active moderation; not signing up for that with light maintenance
- **No changelog/releases** — People clone, they don't install
- **No demo site** — Maintenance burden; screenshots in README suffice
- **No monorepo tooling** — Flat structure is simple and that's the point
