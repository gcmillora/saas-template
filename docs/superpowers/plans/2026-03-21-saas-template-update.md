# SaaS Template Update Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Rewrite saas-template to match batch project's architecture, adding password reset, role/onboarding fields, dashboard page, and full deployment pipeline.

**Architecture:** Go 1.25 backend with Chi router, oapi-codegen strict server, Jet SQL builder, Gorilla sessions. React 19 frontend with Vite, Tailwind v4, shadcn/ui, Orval-generated React Query hooks. Docker-based deployment with Nginx reverse proxy, GitHub Actions CI/CD.

**Tech Stack:** Go 1.25, PostgreSQL 18, Chi v5, Jet v2, oapi-codegen, React 19, Vite 7, Tailwind CSS v4, React Query 5, Orval 8, Docker, GitHub Actions

**Spec:** `docs/superpowers/specs/2026-03-21-saas-template-update-design.md`

**Reference project:** `/Users/greg/dev/projects/batch` — all patterns should match this project exactly unless the spec explicitly deviates.

---

### Task 1: Backend — Config & Providers

Rewrite the config layer to match batch exactly.

**Files:**
- Modify: `backend/go.mod`
- Modify: `backend/config/app.go`
- Modify: `backend/config/provider/env_provider.go`
- Modify: `backend/config/provider/database_provider.go`
- Modify: `backend/config/provider/session_provider.go`
- Modify: `backend/config/provider/logger_provider.go`
- Modify: `backend/config/provider/cache_provider.go`
- Modify: `backend/config/provider/validation_provider.go`
- Create: `backend/config/provider/supabase_provider.go`
- Create: `backend/config/provider/resend_provider.go`
- Modify: `backend/.env.example` (create if not exists)
- Modify: `backend/.env.local`
- Modify: `backend/.env.docker`
- Modify: `backend/.env.test`

- [ ] **Step 1: Update go.mod**

Update module name to `saas-template`, set Go 1.25, and add all required dependencies. Copy batch's `go.mod` as reference, removing batch-specific deps (`go-pdf/fpdf`). Run `go mod tidy` after. The module name should be `saas-template`.

- [ ] **Step 2: Rewrite config/app.go**

Copy batch's `config/app.go` exactly, replacing `batch/` import path with `saas-template/`. The App struct has fields: env, db, rootDir, logger, cache, session, storage, resend. All getters use lazy init pattern. `NewApp()` eagerly initializes env, db, validation, and session.

- [ ] **Step 3: Rewrite all provider files**

Copy each provider file from batch, replacing import paths:
- `env_provider.go` — all env vars including SUPABASE_URL, SUPABASE_KEY, RESEND_API_KEY, RESEND_FROM_EMAIL
- `database_provider.go` — DbProvider with WithTransaction
- `session_provider.go` — CookieStore from SessionSecret
- `logger_provider.go` — slog with production/debug levels
- `cache_provider.go` — go-cache with 5s default, 1m cleanup
- `validation_provider.go` — gookit/validate config
- `supabase_provider.go` — storage-go client (NEW)
- `resend_provider.go` — resend client (NEW)

- [ ] **Step 4: Create/update .env files**

Create `.env.example` with all vars documented. Update `.env.local`, `.env.docker`, `.env.test` to include SUPABASE_URL, SUPABASE_KEY (can be placeholder values for local dev).

- [ ] **Step 5: Run `go mod tidy` and verify compilation**

```bash
cd backend && go mod tidy && go build ./...
```

- [ ] **Step 6: Commit**

```bash
git add backend/go.mod backend/go.sum backend/config/ backend/.env*
git commit -m "feat: update config and providers to match batch patterns"
```

---

### Task 2: Backend — Database Migrations

Replace existing migration with the new schema.

**Files:**
- Modify: `backend/db/migrations/00001_init.sql`
- Create: `backend/db/migrations/00002_password_reset.sql`

- [ ] **Step 1: Rewrite 00001_init.sql**

Replace contents with the spec's migration SQL:

```sql
-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE tenant_tbl (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE user_tbl (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255),
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    auth_provider VARCHAR(50) NOT NULL DEFAULT 'email',
    auth_provider_id VARCHAR(255),
    role VARCHAR(50) NOT NULL DEFAULT 'user',
    onboarding_completed BOOLEAN NOT NULL DEFAULT FALSE,
    tenant_id UUID NOT NULL REFERENCES tenant_tbl(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_user_tbl_email ON user_tbl(email);
CREATE INDEX idx_user_tbl_tenant_id ON user_tbl(tenant_id);

-- +goose Down
DROP TABLE IF EXISTS user_tbl;
DROP TABLE IF EXISTS tenant_tbl;
DROP EXTENSION IF EXISTS "uuid-ossp";
```

- [ ] **Step 2: Create 00002_password_reset.sql**

```sql
-- +goose Up
CREATE TABLE password_reset_tbl (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES user_tbl(id) ON DELETE CASCADE,
    token_hash VARCHAR(64) UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    used_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_password_reset_tbl_token_hash ON password_reset_tbl(token_hash);
CREATE INDEX idx_password_reset_tbl_user_id ON password_reset_tbl(user_id);

-- +goose Down
DROP TABLE IF EXISTS password_reset_tbl;
```

- [ ] **Step 3: Run migrations and generate Jet models**

Start PostgreSQL if not running, then run migrations and generate Jet models. This is required because Tasks 5-9 depend on the generated Jet model code.

```bash
docker compose up postgres -d
cd backend && bash commands.sh migration:up
```

This runs goose migrations, then runs `jet` to generate type-safe Go models in `generated/db/database/public/`.

- [ ] **Step 4: Commit**

```bash
git add backend/db/migrations/ backend/generated/db/
git commit -m "feat: update migrations with role, onboarding, password reset"
```

---

### Task 3: Backend — OpenAPI Specs & Code Generation

Update the OpenAPI specs and codegen config to match the spec.

**Files:**
- Modify: `backend/openapi.yaml`
- Modify: `backend/openapi-public.yaml`
- Modify: `backend/generated/oapi/codegen.yaml`
- Modify: `backend/generated/oapi/public/codegen.yaml`

- [ ] **Step 1: Rewrite openapi.yaml**

Authenticated API with two endpoints:
- `GET /user` (operationId: `get-api-v1-user`) — returns BaseUser
- `PATCH /user/onboarding` (operationId: `patch-api-v1-user-onboarding`) — accepts `{ onboardingCompleted: boolean }`, returns BaseUser

BaseUser schema: id (uuid), email (email), firstName (string, nullable), lastName (string, nullable), role (enum: admin, user), onboardingCompleted (boolean). All required: id, email, role, onboardingCompleted.

UpdateOnboardingRequest schema: `{ onboardingCompleted: boolean }` required.

Error schema: `{ message: string }`.

Match batch's openapi.yaml structure exactly (minus project/placeholder/generate endpoints).

- [ ] **Step 2: Rewrite openapi-public.yaml**

Public API with these endpoints:
- `POST /signin` (operationId: `post-api-v1-signin`)
- `POST /signup` (operationId: `post-api-v1-signup`)
- `POST /signout` (operationId: `post-api-v1-signout`)
- `POST /forgot-password` (operationId: `post-api-v1-forgot-password`)
- `POST /reset-password` (operationId: `post-api-v1-reset-password`)

Copy batch's `openapi-public.yaml`, then remove the `GET /health` endpoint definition (health will be registered directly on the router instead). Keep all other endpoints — batch already has all the listed endpoints with the correct schemas (MessageResponse, Error, ValidationError).

- [ ] **Step 3: Update codegen.yaml files**

`generated/oapi/codegen.yaml`:
```yaml
package: oapi
generate:
  chi-server: true
  models: true
  client: true
  strict-server: true
output: generated/oapi/generated.go
output-options:
  skip-prune: true
```

`generated/oapi/public/codegen.yaml`:
```yaml
package: oapi
generate:
  chi-server: true
  models: true
  client: true
  strict-server: true
output: generated/oapi/public/generated.go
output-options:
  skip-prune: true
```

- [ ] **Step 4: Run code generation**

```bash
cd backend && go tool oapi-codegen -config ./generated/oapi/codegen.yaml ./openapi.yaml && go tool oapi-codegen -config ./generated/oapi/public/codegen.yaml ./openapi-public.yaml
```

- [ ] **Step 5: Verify generated code compiles**

```bash
cd backend && go build ./...
```

- [ ] **Step 6: Commit**

```bash
git add backend/openapi.yaml backend/openapi-public.yaml backend/generated/
git commit -m "feat: update OpenAPI specs and regenerate server code"
```

---

### Task 4: Backend — Middleware

Update middleware to match batch exactly.

**Files:**
- Modify: `backend/internal/webserver/middleware/auth_middleware.go`
- Modify: `backend/internal/webserver/middleware/context.go`
- Modify: `backend/internal/webserver/middleware/session.go`
- Modify: `backend/internal/webserver/middleware/logger.go`
- Modify: `backend/internal/webserver/middleware/rate_limiter.go`

- [ ] **Step 1: Copy all middleware files from batch**

Copy each middleware file from batch verbatim, replacing `batch/` import paths with `saas-template/`:
- `auth_middleware.go` — checks session user_id, returns 401
- `context.go` — contextKey type, ResponseWriterKey/RequestKey, NewContextInjectorMiddleware, GetResponseWriter, GetRequest
- `session.go` — SessionData struct (UserID, TenantID, Role), GetSessionData helper
- `logger.go` — NewLoggerMiddleware + HandleErrorWithLog
- `rate_limiter.go` — visitor-based per-IP rate limiter

- [ ] **Step 2: Verify compilation**

```bash
cd backend && go build ./...
```

- [ ] **Step 3: Commit**

```bash
git add backend/internal/webserver/middleware/
git commit -m "feat: update middleware to match batch patterns"
```

---

### Task 5: Backend — Repository & Mutation Layer

Create/update all data access functions.

**Files:**
- Modify: `backend/internal/app/repository/user_repository.go`
- Create: `backend/internal/app/repository/password_reset_repository.go`
- Modify: `backend/internal/app/repository/tenant_repository.go`
- Modify: `backend/internal/app/repository/pagination.go`
- Modify: `backend/internal/app/mutation/user_mutation.go`
- Create: `backend/internal/app/mutation/password_reset_mutation.go`
- Create: `backend/internal/app/mutation/tenant_mutation.go`

- [ ] **Step 1: Copy repository files from batch**

Copy verbatim (fixing import paths to `saas-template/`):
- `user_repository.go` — GetUserByID (scoped by tenantID), GetUserByEmail, GetUsers
- `password_reset_repository.go` — GetPasswordResetByTokenHash (checks UsedAt IS NULL, ExpiresAt > now)
- `tenant_repository.go` — GetTenantByID, GetTenants
- `pagination.go` — PaginationParams, PaginatedResponse, ApplyPagination

- [ ] **Step 2: Copy mutation files from batch**

Copy verbatim (fixing import paths):
- `user_mutation.go` — CreateUser, DeleteUser, UpdateUserPassword, UpdateUserOnboarding. For UpdateUserOnboarding, rename the field from `OnboardingDismissed` to `OnboardingCompleted` to match our schema. Remove `SetUserHasGenerated` (batch-specific).
- `password_reset_mutation.go` — CreatePasswordReset, InvalidateAllUserTokens
- `tenant_mutation.go` — CreateTenant

Note: The Jet-generated models won't exist yet (need DB running + migrations + jet codegen), so these files won't compile until Task 2's migrations are run and Jet models are generated. That's expected — just ensure the Go code itself is correct.

- [ ] **Step 3: Commit**

```bash
git add backend/internal/app/repository/ backend/internal/app/mutation/
git commit -m "feat: update repository and mutation layer"
```

---

### Task 6: Backend — Util Services

Add password validation and email sending utilities.

**Files:**
- Create: `backend/internal/app/util_service/password/validate.go`
- Create: `backend/internal/app/util_service/password/validate_test.go`
- Create: `backend/internal/app/util_service/email/send_reset_email.go`
- Keep: `backend/internal/app/util_service/datetime_util_service.go`
- Keep: `backend/internal/app/util_service/db_util_service.go`
- Keep: `backend/internal/app/util_service/simple_util_service.go`

- [ ] **Step 1: Create password validation**

Copy `password/validate.go` from batch exactly. It validates: min 8 chars, 1 uppercase, 1 lowercase, 1 digit, 1 special character. Returns `[]string` of error messages.

- [ ] **Step 2: Create password validation tests**

Copy `password/validate_test.go` from batch exactly. Table-driven tests covering: empty, too short, missing uppercase/lowercase/digit/special, valid, edge cases.

- [ ] **Step 3: Run password tests**

```bash
cd backend && go test ./internal/app/util_service/password/...
```

Expected: all tests pass.

- [ ] **Step 4: Create email utility**

Copy `email/send_reset_email.go` from batch. Uses Resend client to send password reset HTML email with reset URL.

- [ ] **Step 5: Commit**

```bash
git add backend/internal/app/util_service/
git commit -m "feat: add password validation and email utilities"
```

---

### Task 7: Backend — App Services

Create business logic layer matching batch patterns.

**Files:**
- Create: `backend/internal/app/app_service/authentication/signin.go`
- Modify: `backend/internal/app/app_service/authentication/signup.go`
- Create: `backend/internal/app/app_service/authentication/signout.go`
- Create: `backend/internal/app/app_service/authentication/forgot_password.go`
- Create: `backend/internal/app/app_service/authentication/reset_password.go`
- Modify: `backend/internal/app/app_service/user/get_user.go`
- Create: `backend/internal/app/app_service/user/update_onboarding.go`

- [ ] **Step 1: Create signin.go**

Copy from batch's `authentication/signin.go`, removing audit.Log calls (no audit logging in template). Fix import paths (`batch/` → `saas-template/`). Function: `SignIn(ctx, app, email, password) (*model.UserTbl, error)` — gets user by email, compares bcrypt hash, returns user or error.

- [ ] **Step 2: Rewrite signup.go**

Rewrite existing `signup.go`. Copy from batch's `authentication/signup.go`, removing audit.Log calls. Fix import paths. Uses `PostSignupBody` struct. Validates passwords match, complexity, email uniqueness. Creates tenant, hashes password, creates user with role "user".

- [ ] **Step 3: Create signout.go**

Simple function matching batch's `authentication/signout.go` but without audit.Log. Just a no-op function (session clearing happens in handler).

```go
package authentication

import (
	"context"
	"saas-template/config"

	"github.com/google/uuid"
)

func SignOut(ctx context.Context, app *config.App, userID uuid.UUID, tenantID uuid.UUID) {
	_ = ctx
	_ = app
}
```

- [ ] **Step 4: Create forgot_password.go**

Copy from batch's `authentication/forgot_password.go`, removing audit.Log calls. Fix import paths. Generates random token, SHA256 hashes it, stores hash in DB, sends email with raw token in URL. Silently returns if user not found (no email enumeration).

- [ ] **Step 5: Create reset_password.go**

Copy from batch's `authentication/reset_password.go`, removing audit.Log calls. Fix import paths. Validates password match + complexity, hashes token to find DB record, updates password, invalidates all user tokens.

- [ ] **Step 6: Rewrite get_user.go**

Copy from batch's `user/get_user.go`, fixing import paths. Adapt the response to use `OnboardingCompleted` instead of `OnboardingDismissed` and remove `HasGenerated`. Returns `GetApiV1User200JSONResponse`.

- [ ] **Step 7: Create update_onboarding.go**

Copy from batch's `user/update_onboarding.go`, fixing import paths. Adapt to use `OnboardingCompleted` instead of `OnboardingDismissed` and remove `HasGenerated`.

- [ ] **Step 8: Commit**

```bash
git add backend/internal/app/app_service/
git commit -m "feat: add app services for auth, user, and password reset"
```

---

### Task 8: Backend — Handlers

Create all HTTP handlers matching batch's thin handler pattern.

**Files:**
- Modify: `backend/internal/webserver/handler/baseHandler.go`
- Modify: `backend/internal/webserver/handler/post_api_v1_signin.go`
- Modify: `backend/internal/webserver/handler/post_api_v1_signup.go`
- Modify: `backend/internal/webserver/handler/post_api_v1_signout.go`
- Create: `backend/internal/webserver/handler/post_api_v1_forgot_password.go`
- Create: `backend/internal/webserver/handler/post_api_v1_reset_password.go`
- Modify: `backend/internal/webserver/handler/get_api_v1_me_user.go` (rename to `get_api_v1_user.go`)
- Create: `backend/internal/webserver/handler/patch_api_v1_user_onboarding.go`
- Delete: `backend/internal/webserver/handler/get_api_v1_health.go` (health moves to router)
- Delete: `backend/internal/webserver/handler/get_api_v1_health_test.go`
- Delete: `backend/internal/webserver/handler/handler_template.go.example` (if exists)

- [ ] **Step 1: Rewrite baseHandler.go**

Copy from batch, removing admin API interface. The CombinedStrictServerInterface should combine only oapi.StrictServerInterface and oapi_public.StrictServerInterface (no admin). Fix import paths.

- [ ] **Step 2: Rewrite post_api_v1_signin.go**

Copy from batch verbatim, fixing import paths. Gets HTTP primitives from context, calls authentication.SignIn, sets session values (user_id, tenant_id, role), returns MessageResponse.

- [ ] **Step 3: Rewrite post_api_v1_signup.go**

Copy from batch, fixing import paths. Creates PostSignupBody from request, calls authentication.SignUp, sets session (user_id, tenant_id), returns MessageResponse.

- [ ] **Step 4: Rewrite post_api_v1_signout.go**

Copy from batch, fixing import paths. Remove audit.Log call (template has no audit). Gets session, calls authentication.SignOut, sets MaxAge = -1, saves session.

- [ ] **Step 5: Create post_api_v1_forgot_password.go**

Copy from batch verbatim, fixing import paths. Calls authentication.ForgotPassword, always returns 200 MessageResponse.

- [ ] **Step 6: Create post_api_v1_reset_password.go**

Copy from batch, fixing import paths. Calls authentication.ResetPassword, handles ResetPasswordResult with validation errors or error.

- [ ] **Step 7: Rename and rewrite get_api_v1_user.go**

Delete `get_api_v1_me_user.go`, create `get_api_v1_user.go`. Copy from batch's `get_api_v1_me_user.go`, fixing import paths. Delegates to user.GetUser.

- [ ] **Step 8: Create patch_api_v1_user_onboarding.go**

Copy from batch, fixing import paths. Adapt to use `OnboardingCompleted` field name instead of `OnboardingDismissed`. Delegates to user.UpdateOnboarding.

- [ ] **Step 9: Delete old files**

Remove `get_api_v1_health.go`, `get_api_v1_health_test.go`, and `handler_template.go.example` if they exist.

- [ ] **Step 10: Commit**

```bash
git add backend/internal/webserver/handler/
git commit -m "feat: update handlers to match batch patterns"
```

---

### Task 9: Backend — Webserver & Entry Points

Update the router setup and main entry point.

**Files:**
- Modify: `backend/internal/webserver/webserver.go`
- Modify: `backend/cmd/webserver/main.go`
- Modify: `backend/cmd/console/main.go`
- Modify: `backend/internal/app/errors/custom_error.go`
- Modify: `backend/internal/app/dto/types.go`

- [ ] **Step 1: Rewrite webserver.go**

Copy from batch's webserver.go, removing the admin API route group and admin imports. Keep:
1. Global middleware: Logger, Recoverer, CORS
2. Health endpoint: `r.Get("/health", ...)`
3. Public route group: RateLimiter + ContextInjector + oapi_public handlers
4. Authenticated route group: ContextInjector + AuthMiddleware + oapi handlers
5. Graceful shutdown

Fix import paths (`batch/` → `saas-template/`). Remove admin imports and admin route group.

**CORS deviation from batch:** Add `PATCH` to AllowedMethods (batch only has GET, POST, PUT, DELETE, OPTIONS). This is needed because the template has a `PATCH /user/onboarding` endpoint.

- [ ] **Step 2: Rewrite cmd/webserver/main.go**

Copy from batch, fixing import paths. Simple: create App, create Webserver, handle "routes:list" command, else start.

- [ ] **Step 2b: Copy cmd/console/main.go**

Copy from batch, fixing import paths. This is the CLI commands entry point.

- [ ] **Step 3: Copy errors and DTO**

Copy `errors/custom_error.go` from batch (BadRequestError, UnauthorizedError, ForbiddenError, ValidationError, NotFoundError). Copy `dto/types.go` (empty package).

- [ ] **Step 4: Verify compilation**

```bash
cd backend && go build ./...
```

Note: This may not fully compile until Jet models exist, but handler/middleware/webserver code should be syntactically correct.

- [ ] **Step 5: Commit**

```bash
git add backend/internal/webserver/ backend/cmd/ backend/internal/app/errors/ backend/internal/app/dto/
git commit -m "feat: update webserver router and entry points"
```

---

### Task 10: Backend — Build & Dev Tooling

Update Dockerfile, commands.sh, Air config, linting config.

**Files:**
- Modify: `backend/Dockerfile_webserver` (rename to `backend/Dockerfile`)
- Modify: `backend/commands.sh`
- Modify: `backend/.air.toml`
- Modify: `backend/.golangci.yml`

- [ ] **Step 1: Create backend/Dockerfile**

Copy from batch's Dockerfile exactly. Multi-stage: Go 1.25 Alpine builder, build server + goose, Alpine 3.21 runtime with nonroot user. Delete old `Dockerfile_webserver`.

- [ ] **Step 2: Rewrite commands.sh**

Copy from batch's commands.sh exactly. Includes: webserver, console, lint, lint:fix, format, test, openapi:codegen (both specs + orval), migration:codegen, migration:up (runs migrations + jet + test DB), migration:down, migration:reset, migration:status.

Remove the admin codegen line from openapi:codegen since there's no admin API.

- [ ] **Step 3: Copy .air.toml**

Copy from batch exactly. Watches go/yaml/html files, excludes generated/migrations dirs.

- [ ] **Step 4: Copy .golangci.yml**

Copy from batch exactly. Version 2 config with errchkjson, gocritic, importas linters.

- [ ] **Step 5: Commit**

```bash
git add backend/Dockerfile backend/commands.sh backend/.air.toml backend/.golangci.yml
git rm backend/Dockerfile_webserver 2>/dev/null || true
git commit -m "feat: update backend build and dev tooling"
```

---

### Task 11: Frontend — Core Setup & Config

Update frontend config files, dependencies, and build setup.

**Files:**
- Modify: `frontend/package.json`
- Modify: `frontend/vite.config.ts`
- Modify: `frontend/tsconfig.json`
- Modify: `frontend/tsconfig.app.json`
- Modify: `frontend/components.json`
- Modify: `frontend/eslint.config.js`
- Modify: `frontend/index.html`
- Modify: `frontend/orval.config.ts`
- Create: `frontend/Dockerfile`
- Create: `frontend/nginx.conf`

- [ ] **Step 1: Update package.json**

Copy batch's package.json, removing batch-specific deps: konva, react-konva, papaparse, recharts, @atlaskit/pragmatic-drag-and-drop and their @types. Keep all other deps. Update project name to `saas-template`.

- [ ] **Step 2: Copy config files**

Copy from batch exactly:
- `vite.config.ts` — Vite with tailwindcss plugin, react-swc, path alias, proxy to 8008
- `tsconfig.json` — references to app/node configs, path alias
- `tsconfig.app.json` — strict TS config, path mapping
- `tsconfig.node.json` — node-specific TS config
- `components.json` — shadcn/ui new-york style config
- `eslint.config.js` — ts-eslint with react-hooks, ignore generated dirs
- `index.html` — update title to "SaaS Template"

- [ ] **Step 3: Update orval.config.ts**

Copy from batch, removing the v1Admin target. Keep only v1 (from openapi.yaml) and v1Public (from openapi-public.yaml).

- [ ] **Step 4: Create frontend Dockerfile**

Copy from batch. Multi-stage: Node 22 Alpine with Bun, build, serve with Nginx Alpine.

- [ ] **Step 5: Create frontend nginx.conf**

Copy from batch. SPA fallback, cache static assets.

- [ ] **Step 6: Install dependencies**

```bash
cd frontend && bun install
```

- [ ] **Step 7: Commit**

```bash
git add frontend/package.json frontend/bun.lock frontend/vite.config.ts frontend/tsconfig.json frontend/tsconfig.app.json frontend/components.json frontend/eslint.config.js frontend/index.html frontend/orval.config.ts frontend/Dockerfile frontend/nginx.conf
git commit -m "feat: update frontend config and dependencies"
```

---

### Task 12: Frontend — Styling & Utilities

Update CSS theme and utility functions.

**Files:**
- Modify: `frontend/src/index.css`
- Modify: `frontend/src/lib/utils.ts`
- Modify: `frontend/src/App.tsx`
- Modify: `frontend/src/main.tsx`

- [ ] **Step 1: Copy index.css from batch**

Copy exactly — includes Inter font import, Tailwind imports, theme inline vars, :root color tokens (oklch values), base layer styles, component layer styles (card shadow, button glow).

- [ ] **Step 2: Verify utils.ts matches batch**

Should already be: `cn()` function with clsx + tailwind-merge. Copy from batch if different.

- [ ] **Step 3: Copy App.tsx and main.tsx from batch**

App.tsx wraps ErrorBoundary > QueryProvider > RouterProvider. main.tsx renders App in StrictMode.

- [ ] **Step 4: Commit**

```bash
git add frontend/src/index.css frontend/src/lib/utils.ts frontend/src/App.tsx frontend/src/main.tsx
git commit -m "feat: update frontend styling and app shell"
```

---

### Task 13: Frontend — API Layer & Auth Context

Update fetch wrappers, generate API hooks, update auth context.

**Files:**
- Modify: `frontend/src/services/api/axios-v1.ts`
- Modify: `frontend/src/services/api/axios-v1-public.ts`
- Modify: `frontend/src/contexts/authContext.tsx`
- Modify: `frontend/src/providers/QueryProvider.tsx`

- [ ] **Step 1: Copy fetch wrappers from batch**

`axios-v1.ts` — customInstance with `/api/v1` prefix, credentials include, throw on !ok, return json.
`axios-v1-public.ts` — same but `/api/public/v1` prefix.

Both should match batch exactly.

- [ ] **Step 2: Generate API hooks with Orval**

```bash
cd frontend && bunx orval
```

This generates `v1.ts` and `v1-public.ts` from the OpenAPI specs.

- [ ] **Step 3: Copy auth context from batch**

Copy `contexts/authContext.tsx`. Uses `useGetApiV1User` with `retry: false`. Redirects to `/signin` on error. Provides user object via context.

- [ ] **Step 4: Verify QueryProvider.tsx matches batch**

Should export queryClient and wrap children with QueryClientProvider.

- [ ] **Step 5: Commit**

```bash
git add frontend/src/services/api/ frontend/src/contexts/ frontend/src/providers/
git commit -m "feat: update API layer and auth context"
```

---

### Task 14: Frontend — shadcn/ui Components

Ensure all required shadcn components are installed.

**Files:**
- Modify/Create: `frontend/src/components/ui/` (multiple files)

- [ ] **Step 1: Install required shadcn components**

Need: button, card, input, label, separator, sidebar, tooltip, dropdown-menu, avatar, sheet. Some already exist. Install missing ones:

```bash
cd frontend && bunx shadcn@latest add dropdown-menu avatar
```

(button, card, input, label, separator, sidebar, tooltip, sheet already exist)

- [ ] **Step 2: Verify all components exist**

Check that all of these files exist in `frontend/src/components/ui/`:
- button.tsx, card.tsx, input.tsx, label.tsx, separator.tsx, sidebar.tsx, tooltip.tsx, dropdown-menu.tsx, avatar.tsx, sheet.tsx

- [ ] **Step 3: Commit**

```bash
git add frontend/src/components/ui/
git commit -m "feat: add missing shadcn/ui components"
```

---

### Task 15: Frontend — Layout Components

Create app layout, sidebar, error boundary, and guest route.

**Files:**
- Modify: `frontend/src/components/AppLayout.tsx`
- Modify: `frontend/src/components/app-sidebar.tsx`
- Modify: `frontend/src/components/ErrorBoundary.tsx`
- Create: `frontend/src/components/GuestRoute.tsx`
- Create: `frontend/src/components/PasswordRequirements.tsx`

- [ ] **Step 1: Copy AppLayout.tsx from batch**

Copy exactly. Uses SidebarProvider (defaultOpen false), AppSidebar, SidebarInset with header (SidebarTrigger) and main (Outlet).

- [ ] **Step 2: Rewrite app-sidebar.tsx**

Adapt from batch's sidebar. Simplify significantly:
- Remove project-related items (FolderOpen, recent projects, FileText)
- Remove admin analytics link
- Remove StamplLogo, replace with generic app name
- Keep: sidebar structure, user avatar/dropdown in footer, sign out button
- Add: single "Dashboard" nav item with LayoutDashboard icon (from lucide-react)
- Keep: user initials, display name, email in footer
- Use `usePostApiV1Signout` from v1-public hooks

The sidebar should have:
1. Header with app name
2. Content with single "Dashboard" menu item
3. Footer with user dropdown (sign out option)

- [ ] **Step 3: Copy ErrorBoundary.tsx**

Copy from batch, removing StamplLogo reference. Use a simple text fallback instead.

- [ ] **Step 4: Create GuestRoute.tsx**

Copy from batch's GuestRoute.tsx, changing the redirect destination from `/projects` to `/dashboard`.

- [ ] **Step 5: Create PasswordRequirements.tsx**

Copy from batch exactly. Shows password requirement checklist with green checkmarks.

- [ ] **Step 6: Commit**

```bash
git add frontend/src/components/
git commit -m "feat: add layout, sidebar, and route guard components"
```

---

### Task 16: Frontend — Pages

Create all page components.

**Files:**
- Modify: `frontend/src/pages/SignIn.tsx`
- Modify: `frontend/src/pages/SignUp.tsx`
- Create: `frontend/src/pages/ForgotPassword.tsx`
- Create: `frontend/src/pages/ResetPassword.tsx`
- Create: `frontend/src/pages/Dashboard.tsx`
- Create: `frontend/src/pages/ErrorPage.tsx`
- Create: `frontend/src/pages/NotFound.tsx`

- [ ] **Step 1: Rewrite SignIn.tsx**

Copy from batch, removing AuthShowcasePanel and StamplLogo references. Simplify to a centered card layout (no split-screen showcase). Keep: email/password fields, error display, forgot password link, signup link, loading state with Loader2.

- [ ] **Step 2: Rewrite SignUp.tsx**

Copy from batch, removing AuthShowcasePanel and StamplLogo. Simplify to centered card layout. Keep: email, password with PasswordRequirements, confirm password, error display, signin link, loading state.

- [ ] **Step 3: Create ForgotPassword.tsx**

Copy from batch exactly. Card-based layout with email input, success message after submit.

- [ ] **Step 4: Create ResetPassword.tsx**

Copy from batch exactly. Card-based layout with password + confirm password, token from URL params, success redirect to signin.

- [ ] **Step 5: Create Dashboard.tsx**

Simple welcome page:

```tsx
import { useContext } from "react";
import { AuthContext } from "@/contexts/authContext";

export function Dashboard() {
  const { user } = useContext(AuthContext);
  const greeting = user?.firstName ? `Welcome, ${user.firstName}` : "Welcome";

  return (
    <div>
      <h1 className="text-2xl font-bold tracking-tight">{greeting}</h1>
    </div>
  );
}
```

- [ ] **Step 6: Create ErrorPage.tsx**

Copy from batch, removing StamplLogo. Use simple text for the 404 display instead of the logo-in-404 treatment.

- [ ] **Step 7: Create NotFound.tsx**

Copy from batch, removing StamplLogo. Simple text-based 404 page.

- [ ] **Step 8: Commit**

```bash
git add frontend/src/pages/
git commit -m "feat: add all page components"
```

---

### Task 17: Frontend — Routing

Update router and route definitions.

**Files:**
- Modify: `frontend/src/router.tsx`
- Modify: `frontend/src/Routes.tsx`

- [ ] **Step 1: Rewrite router.tsx**

Adapt from batch's router.tsx:

```tsx
import { createBrowserRouter, Navigate } from "react-router";
import { AppLayout } from "./components/AppLayout";
import { AuthContextProvider } from "./contexts/authContext";
import { GuestRoute } from "./components/GuestRoute";
import { SignIn } from "./pages/SignIn";
import { SignUp } from "./pages/SignUp";
import { ForgotPassword } from "./pages/ForgotPassword";
import { ResetPassword } from "./pages/ResetPassword";
import { Dashboard } from "./pages/Dashboard";
import { ErrorPage } from "./pages/ErrorPage";
import { NotFound } from "./pages/NotFound";

function AuthenticatedLayout() {
  return (
    <AuthContextProvider>
      <AppLayout />
    </AuthContextProvider>
  );
}

export const router = createBrowserRouter([
  {
    errorElement: <ErrorPage />,
    children: [
      {
        element: <AuthenticatedLayout />,
        children: [
          {
            index: true,
            element: <Navigate to="/dashboard" replace />,
          },
          {
            path: "dashboard",
            element: <Dashboard />,
          },
        ],
      },
      {
        path: "/signin",
        element: <GuestRoute><SignIn /></GuestRoute>,
      },
      {
        path: "/signup",
        element: <GuestRoute><SignUp /></GuestRoute>,
      },
      {
        path: "/forgot-password",
        element: <GuestRoute><ForgotPassword /></GuestRoute>,
      },
      {
        path: "/reset-password",
        element: <GuestRoute><ResetPassword /></GuestRoute>,
      },
      {
        path: "*",
        element: <NotFound />,
      },
    ],
  },
]);
```

- [ ] **Step 2: Update Routes.tsx**

```tsx
export const Routes = {
  home: "/",
  dashboard: "/dashboard",
} as const;

export default Routes;
```

- [ ] **Step 3: Verify frontend builds**

```bash
cd frontend && bun run build
```

- [ ] **Step 4: Commit**

```bash
git add frontend/src/router.tsx frontend/src/Routes.tsx
git commit -m "feat: update routing with dashboard and auth pages"
```

---

### Task 18: Infrastructure — Docker Compose

Create dev and prod Docker Compose files.

**Files:**
- Modify: `docker-compose.yml`
- Create: `docker-compose.prod.yml`
- Create: `nginx/nginx.conf`

- [ ] **Step 1: Rewrite docker-compose.yml**

Adapt from batch. Change container names and profile names from `batch-*` to `saas-*`. Keep same structure: postgres (18), go backend (1.25 Alpine + Air), frontend (Node 22 + Bun). Database name: `saas_template`. Update all references.

- [ ] **Step 2: Create docker-compose.prod.yml**

Adapt from batch. Change GHCR image names from `batch-*` to `saas-template-*`. Change container names. Keep: backend (with migration entrypoint + health check), frontend (Nginx), nginx (reverse proxy), dozzle (logging). Update server_name in health check URL.

- [ ] **Step 3: Create nginx/nginx.conf**

Copy from batch. Update server_name. Keep: upstream backend (8008), upstream frontend (80), proxy /api/ to backend, proxy / to frontend.

- [ ] **Step 4: Commit**

```bash
git add docker-compose.yml docker-compose.prod.yml nginx/
git commit -m "feat: add dev and production Docker Compose setup"
```

---

### Task 19: Infrastructure — CI/CD

Create GitHub Actions workflows.

**Files:**
- Create: `.github/workflows/ci.yml`
- Create: `.github/workflows/deploy.yml`

- [ ] **Step 1: Create ci.yml**

Copy from batch. Three jobs: backend-lint, backend-test, frontend. Uses Go from go.mod, golangci-lint-action, Bun for frontend. No changes needed except it works with the template structure.

- [ ] **Step 2: Create deploy.yml**

Copy from batch. Update image names from `batch-*` to `saas-template-*`. Update deploy directory from `~/batch` to `~/saas-template`. Keep: workflow_run trigger on CI success, matrix build for backend/frontend, GHCR push, SCP + SSH deploy, health check loop.

- [ ] **Step 3: Commit**

```bash
git add .github/workflows/
git commit -m "feat: add CI/CD pipelines"
```

---

### Task 20: Project Files — CLAUDE.md, .env.example, .gitignore, README

Create/update project-level documentation and config.

**Files:**
- Create: `CLAUDE.md`
- Create: `backend/.env.example`
- Modify: `.gitignore`

- [ ] **Step 1: Create CLAUDE.md**

Adapt from batch's CLAUDE.md. Remove all batch-specific references (projects, placeholders, PDF, canvas, admin API). Keep: commands, code generation workflow, backend practices (handler pattern, app service pattern, repository/mutation pattern, adding new endpoint flow, adding new table flow), frontend practices (component conventions, page conventions, styling, API layer, file naming). Replace "Stampl" with generic "SaaS Template".

- [ ] **Step 2: Create backend/.env.example**

```
APP_ENV=local
APP_SECRET=change-me-to-a-random-secret
APP_BASE_URL=http://localhost:8009
SERVER_PORT=8008
DATABASE_URL=postgresql://postgres:postgres@localhost:5432/saas_template?sslmode=disable
SESSION_SECRET=change-me-to-a-random-secret
LOG_LEVEL=debug
SUPABASE_URL=https://your-project.supabase.co
SUPABASE_KEY=your-supabase-key
RESEND_API_KEY=
RESEND_FROM_EMAIL=noreply@example.com
GOOGLE_CLIENT_ID=
GOOGLE_CLIENT_SECRET=
GITHUB_CLIENT_ID=
GITHUB_CLIENT_SECRET=
```

- [ ] **Step 3: Update .gitignore**

Ensure it includes: `.env.local`, `.env.docker`, `.env.test`, `tmp/`, `generated/db/`, `node_modules/`, `dist/`, `.worktrees/`.

- [ ] **Step 4: Clean up old files**

Remove files that shouldn't be in the template:

```bash
git rm backend/mise.toml 2>/dev/null || true
git rm backend/spec.md 2>/dev/null || true
git rm -r .worktrees/ 2>/dev/null || true
rm -rf frontend/dist/
```

Add `.worktrees/` and `frontend/dist/` to `.gitignore` if not already there.

- [ ] **Step 5: Commit**

```bash
git add CLAUDE.md backend/.env.example .gitignore
git commit -m "feat: add project documentation and config"
```

---

### Task 21: Verification — Full Build Check

Verify everything compiles and builds correctly.

- [ ] **Step 1: Start PostgreSQL and run migrations**

```bash
docker compose up postgres -d
cd backend && bash commands.sh migration:up
```

This runs migrations + generates Jet models.

- [ ] **Step 2: Verify backend compiles**

```bash
cd backend && go build ./...
```

- [ ] **Step 3: Run backend tests**

```bash
cd backend && go test ./...
```

- [ ] **Step 4: Run backend lint**

```bash
cd backend && bash commands.sh lint
```

- [ ] **Step 5: Verify frontend builds**

```bash
cd frontend && bun run build
```

- [ ] **Step 6: Run frontend lint**

```bash
cd frontend && bun run lint
```

- [ ] **Step 7: Fix any issues found**

Address compilation errors, lint warnings, or test failures.

- [ ] **Step 8: Final commit**

```bash
git add -A
git commit -m "fix: resolve build issues from template update"
```
