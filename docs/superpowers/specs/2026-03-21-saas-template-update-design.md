# SaaS Template Update — Design Spec

Update saas-template to match batch project's architecture, patterns, and practices. Add password reset, role field, onboarding flag, and full deployment pipeline. Keep feature set minimal — auth, user management, and a basic dashboard shell.

## Backend

### Go Version & Dependencies

Upgrade to Go 1.25. Match batch's dependency set:
- chi/v5, go-jet/jet/v2, oapi-codegen, gorilla/sessions
- go-cache, supabase-community/storage-go, resend-go/v2
- gookit/validate, golang.org/x/crypto, golang.org/x/time
- pressly/goose/v3

### Config & Providers

Rewrite `config/app.go` to match batch's lazy-loading App struct pattern. Each provider initializes on first access via getter methods.

Providers (one file each in `config/provider/`):
- `app.go` — App struct definition with all provider fields
- `env_provider.go` — environment variables (APP_ENV, APP_BASE_URL, APP_SECRET, SESSION_SECRET, DATABASE_URL, SERVER_PORT, LOG_LEVEL, SUPABASE_URL, SUPABASE_KEY, RESEND_API_KEY, RESEND_FROM_EMAIL, GOOGLE_CLIENT_ID, GOOGLE_CLIENT_SECRET, GITHUB_CLIENT_ID, GITHUB_CLIENT_SECRET)
- `database_provider.go` — PostgreSQL connection
- `session_provider.go` — Gorilla cookie store
- `logger_provider.go` — slog.Logger
- `cache_provider.go` — go-cache in-memory
- `validation_provider.go` — gookit/validate
- `supabase_provider.go` — Supabase Storage client
- `resend_provider.go` — Resend email client

### Database Migrations

Three Goose migrations in `db/migrations/`:

`00001_init.sql`:
```sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE tenant_tbl (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE TABLE user_tbl (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255),
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    auth_provider VARCHAR(50) DEFAULT 'email',
    auth_provider_id VARCHAR(255),
    role VARCHAR(50) NOT NULL DEFAULT 'user',
    tenant_id UUID NOT NULL REFERENCES tenant_tbl(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE INDEX idx_user_tbl_email ON user_tbl(email);
CREATE INDEX idx_user_tbl_tenant_id ON user_tbl(tenant_id);
```

`00002_password_reset.sql`:
```sql
CREATE TABLE password_reset_tbl (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES user_tbl(id) ON DELETE CASCADE,
    token VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    used BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE INDEX idx_password_reset_tbl_token ON password_reset_tbl(token);
CREATE INDEX idx_password_reset_tbl_user_id ON password_reset_tbl(user_id);
```

`00003_onboarding.sql`:
```sql
ALTER TABLE user_tbl ADD COLUMN onboarding_completed BOOLEAN NOT NULL DEFAULT FALSE;
```

### OpenAPI Specs

Two specs, matching batch's structure.

`openapi-public.yaml` (base: `/api/public/v1`, no auth):
- `POST /signin` — email + password, returns user or 401
- `POST /signup` — email, password, confirmPassword, firstName, lastName, returns user or 400
- `POST /forgot-password` — email, returns 200 (always succeeds for security)
- `POST /reset-password` — token + newPassword + confirmPassword, returns 200 or 400
- `GET /health` — returns 200

`openapi.yaml` (base: `/api/v1`, requires auth):
- `GET /user` — returns authenticated user
- `PATCH /user/onboarding` — marks onboarding as completed
- `POST /signout` — clears session

Schemas: BaseUser (id, email, firstName, lastName, role, onboardingCompleted), Error (message), SigninRequest, SignupRequest, ForgotPasswordRequest, ResetPasswordRequest, OnboardingRequest.

Code generation config files:
- `generated/oapi/codegen.yaml` — strict-server for openapi.yaml
- `generated/oapi/public/codegen.yaml` — strict-server for openapi-public.yaml

### Handlers

File naming: `{method}_api_v1_{resource}.go` for authenticated, `{method}_api_public_v1_{resource}.go` for public.

Each handler has a comment `// METHOD (/route)` above the function. All handlers are thin — delegate to app_service, return result. Session management (set/clear cookies) stays in handlers since it needs HTTP primitives.

Files:
- `post_api_public_v1_signin.go` — call authentication.SignIn, set session
- `post_api_public_v1_signup.go` — call authentication.SignUp
- `post_api_public_v1_forgot_password.go` — call authentication.ForgotPassword
- `post_api_public_v1_reset_password.go` — call authentication.ResetPassword
- `get_api_public_v1_health.go` — return 200
- `get_api_v1_user.go` — call user.GetUser
- `patch_api_v1_user_onboarding.go` — call user.UpdateOnboarding
- `post_api_v1_signout.go` — clear session, return 200

`baseHandler.go` — Handler struct with `app *config.App` field and constructor.

### App Services

One file per operation, organized by domain in `internal/app/app_service/`.

`authentication/`:
- `sign_in.go` — validate credentials, return user
- `sign_up.go` — validate input, create tenant, hash password, create user
- `forgot_password.go` — find user by email, generate token, send email via Resend
- `reset_password.go` — validate token, hash new password, update user

`user/`:
- `get_user.go` — get user by ID from session
- `update_onboarding.go` — set onboarding_completed = true

`tenant/`:
- `create_tenant.go` — create tenant record (called from signup)

`util_service/`:
- `email.go` — send email via Resend
- `password.go` — password validation rules

All functions follow signature: `func Name(ctx context.Context, app *config.App, ...params) (ResponseType, error)`. Extract sessionData (userID, tenantID) from session at start. All queries scoped by tenantID.

### Repository & Mutation

`repository/`:
- `user_repository.go` — GetUserByEmail, GetUserByID
- `password_reset_repository.go` — GetPasswordResetByToken
- `tenant_repository.go` — GetTenantByID
- `pagination.go` — shared pagination helper

`mutation/`:
- `user_mutation.go` — CreateUser, UpdateUser
- `password_reset_mutation.go` — CreatePasswordReset, UpdatePasswordReset
- `tenant_mutation.go` — CreateTenant

All use Jet SQL builder. First param `ctx context.Context`, second `qrm.DB`. Mutations set CreatedAt/UpdatedAt. Use `RETURNING(ctbl.AllColumns)`.

### Middleware

`internal/webserver/middleware/`:
- `auth_middleware.go` — check session has valid user_id, return 401 if not
- `context.go` — inject http.Request and http.ResponseWriter into context
- `session.go` — SessionData struct and GetSessionData helper
- `logger.go` — request/response logging
- `rate_limiter.go` — per-IP rate limiting (1 req/sec, burst 5)

### Webserver Setup

`internal/webserver/webserver.go` — Chi router setup matching batch:
1. Global: Logger, Recoverer, CORS
2. Public routes group: ContextInjector, RateLimiter, public handlers
3. Authenticated routes group: ContextInjector, AuthMiddleware, authenticated handlers
4. Graceful shutdown

### Entry Points

`cmd/webserver/main.go` — initialize App, start webserver
`cmd/console/main.go` — CLI commands entry point

### Error Handling

`internal/app/errors/` — error types matching batch. ResponseErrorHandlerFunc on the server for consistent error responses.

### DTO

`internal/app/dto/` — data transfer objects for mapping between layers.

## Frontend

### Dependencies

Match batch's versions:
- react@19, react-dom@19, react-router@7
- @tanstack/react-query@5, vite@7, tailwindcss@4
- shadcn/ui components, lucide-react, luxon@3
- tailwind-merge, clsx

Build tooling: Bun, Vite with Rolldown.

### Pages

`pages/`:
- `SignIn.tsx` — email/password form, uses usePostApiPublicV1Signin
- `SignUp.tsx` — registration form with password confirmation
- `ForgotPassword.tsx` — email form, uses usePostApiPublicV1ForgotPassword
- `ResetPassword.tsx` — new password form with token from URL
- `Dashboard.tsx` — simple welcome page showing user's name from auth context
- `ErrorPage.tsx` — error boundary page
- `NotFound.tsx` — 404 page

### Components

`components/`:
- `AppLayout.tsx` — authenticated layout wrapper with sidebar, matching batch
- `app-sidebar.tsx` — sidebar with Dashboard nav link + sign out button
- `GuestRoute.tsx` — redirects authenticated users away from auth pages
- `ErrorBoundary.tsx` — React error boundary
- `PasswordRequirements.tsx` — password strength indicator

`components/ui/` — shadcn components needed:
- button, card, input, label, separator, sidebar, tooltip, dropdown-menu, avatar, sheet

### Routing

`router.tsx` matching batch's pattern:
```
createBrowserRouter([
  {
    errorElement: <ErrorPage />,
    children: [
      {
        element: <AuthenticatedLayout />,  // wraps with AuthContextProvider
        children: [
          { path: "/", redirect to "/dashboard" },
          { path: "/dashboard", element: <Dashboard /> },
        ]
      },
      // Guest routes
      { path: "/signin", element: <GuestRoute><SignIn /></GuestRoute> },
      { path: "/signup", element: <GuestRoute><SignUp /></GuestRoute> },
      { path: "/forgot-password", element: <GuestRoute><ForgotPassword /></GuestRoute> },
      { path: "/reset-password", element: <GuestRoute><ResetPassword /></GuestRoute> },
      { path: "*", element: <NotFound /> },
    ]
  }
])
```

`Routes.tsx` — enum for route paths.

### Auth Context

`contexts/authContext.tsx` — matches batch. Calls useGetApiV1User on mount. Provides user object. Redirects to /signin on 401.

### API Layer

`services/api/`:
- `v1.ts` — generated by Orval from openapi.yaml
- `v1-public.ts` — generated by Orval from openapi-public.yaml
- `axios-v1.ts` — fetch wrapper, base `/api/v1`, credentials include
- `axios-v1-public.ts` — fetch wrapper, base `/api/public/v1`, credentials include

`orval.config.ts` — two targets matching batch's config structure.

### Styling

`index.css` — Tailwind CSS v4 with semantic color tokens matching batch (primary, secondary, muted, destructive, etc.).

`lib/utils.ts` — `cn()` function (clsx + tailwind-merge).

### Other Frontend Files

- `App.tsx` — root component
- `main.tsx` — entry point
- `vite.config.ts` — Vite config with Tailwind plugin, React SWC, path alias
- `tsconfig.json`, `tsconfig.app.json`, `tsconfig.node.json`
- `eslint.config.js`
- `components.json` — shadcn/ui config
- `index.html`
- `Dockerfile` — multi-stage (Bun build + Nginx)
- `nginx.conf` — frontend nginx config

## Infrastructure

### Docker Compose — Development

`docker-compose.yml`:
- `postgres` — PostgreSQL 18, port 5432, health check, persistent volume
- `saas-backend` (profile: `saas-backend`) — Go 1.25 Alpine, Air hot-reload, runs migrations + Jet on start, port 8008
- `saas-fe` (profile: `saas-fe`) — Node 22 Alpine + Bun, Vite dev server, port 8009

### Docker Compose — Production

`docker-compose.prod.yml`:
- `backend` — GHCR image, runs migrations on start, health check
- `frontend` — GHCR image, static Nginx
- `nginx` — reverse proxy, port 80
- `dozzle` — container log viewer, port 8888 (localhost only)

### Nginx

`nginx/nginx.conf` — reverse proxy matching batch. Routes `/api/` to backend, everything else to frontend.

### Dockerfiles

Backend `Dockerfile` — multi-stage:
- Builder: Go 1.25 Alpine, build binary + Goose
- Runtime: Alpine 3.21, nonroot user

Frontend `Dockerfile` — multi-stage:
- Builder: Node 22 Alpine, Bun install + build
- Runtime: Nginx Alpine, serve dist/

### GitHub Actions

`.github/workflows/ci.yml`:
- Trigger: push to any branch
- Backend job: Go setup, lint (golangci-lint), test
- Frontend job: Bun setup, install, lint, build

`.github/workflows/deploy.yml`:
- Trigger: CI success on main
- Concurrency: cancel previous deploys
- Build + push backend and frontend images to GHCR
- SCP config files to VPS
- SSH: create .env, docker compose pull, up -d --force-recreate, health check, prune

### Project Files

`CLAUDE.md` — project guidelines generalized from batch (commands, code generation workflow, backend/frontend practices).

`.env.example` — all environment variables with descriptions.

`commands.sh` — matching batch's command set (webserver, test, lint, lint:fix, format, openapi:codegen, migration:codegen, migration:up, migration:down).

`.gitignore` — Go + Node ignores, .env files, generated directories.

`README.md` — project setup instructions.

## What's Excluded

- Admin API and admin middleware
- Projects, placeholders, PDF generation, canvas editor
- CSV upload/parsing
- Analytics
- Audit logging
- Onboarding UI (modal, checklist) — only the backend flag
- Konva, papaparse, recharts, go-pdf/fpdf dependencies
- Application-specific business logic
