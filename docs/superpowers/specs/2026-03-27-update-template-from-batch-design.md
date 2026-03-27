# Update SaaS Template from Batch App

**Date:** 2026-03-27
**Status:** Approved

## Overview

Port general infrastructure improvements from the batch app (`/Users/greg/dev/projects/batch/`) to the saas-template. Strip all batch-specific features (PDF generation, projects, placeholders, editor, Supabase storage). The goal is a clean, production-ready SaaS template with proper security, audit logging, admin tier, and deployment infra.

## Decisions

- **Migration strategy:** Add new migrations on top of existing 2 (no rewrite)
- **Auth page design:** Split layout with subtle gradient/pattern right panel (no content to customize)
- **Admin tier:** Include full plumbing (middleware, route group, OpenAPI spec)
- **Landing page:** No landing page; `/` redirects to `/signin` with GuestRoute guard in place
- **Audit logging:** Include full audit service and audit_log_tbl
- **Production infra:** Include docker-compose.prod.yml, nginx reverse proxy, frontend Dockerfile

---

## Backend Changes

### New Migrations

**`00003_audit_log.sql`**
- Create `audit_log_tbl` with columns: id (UUID PK), action (VARCHAR), actor_id (UUID FK nullable), tenant_id (UUID FK nullable), ip_address (VARCHAR), metadata (JSONB), created_at (TIMESTAMPTZ)
- Indexes on: action, actor_id, created_at, (tenant_id, created_at)
- Down migration drops the table

**`00004_user_role.sql`**
- Add `role` column to `user_tbl`: VARCHAR(50) NOT NULL DEFAULT 'user'
- Add CHECK constraint: role IN ('admin', 'user')
- Down migration drops the column

### Security Fixes

**`config/provider/session_provider.go`**
- Add cookie options: Path="/", MaxAge=86400*30 (30 days), HttpOnly=true, Secure=true when env != local, SameSite=Lax

**`internal/webserver/middleware/logger.go`**
- Change error response from `err.Error()` to generic `"Internal server error"` to prevent leaking internals

**`internal/app/app_service/authentication/signup.go`**
- Change duplicate email error from "a user with this email already exists" to "unable to create account"

### New Files

**`internal/app/app_service/audit/log.go`**
- `Log(ctx, app, action, actorID, tenantID, metadata)` function
- Extract IP address from X-Forwarded-For header, fallback to RemoteAddr
- Get request from context middleware
- Insert audit log record via mutation

**`internal/app/mutation/audit_log_mutation.go`**
- `CreateAuditLog(ctx, db, model)` — Insert into audit_log_tbl using Jet

**`internal/webserver/middleware/admin_middleware.go`**
- Extract session, check `role` value equals "admin"
- Return 403 Forbidden JSON response if not admin

**`openapi-admin.yaml`**
- Admin OpenAPI spec with a placeholder health check endpoint
- Minimal spec for codegen to work

**`generated/oapi/admin/codegen.yaml`**
- oapi-codegen config for admin API (chi-server, models, strict-server)

**`backend/.dockerignore`**
- Exclude .git, generated docs, test fixtures from Docker context

### Updated Files

**`internal/webserver/webserver.go`**
- Add third route group: `/api/admin/v1/*` with admin middleware and rate limiter (5 req/sec, burst 10)
- Add `http.MethodPatch` to CORS AllowedMethods

**`commands.sh`**
- Add admin OpenAPI codegen step in the openapi:codegen command

**`internal/app/app_service/authentication/signup.go`**
- Add `audit.Log(ctx, app, "signup", &created.ID, &created.TenantID, nil)` after successful signup

**`internal/app/app_service/authentication/signin.go`**
- Add audit log for failed signin (email in metadata)
- Add audit log for successful signin

**`internal/app/app_service/authentication/signout.go`**
- Replace no-op stubs with `audit.Log(ctx, app, "signout", &userID, &tenantID, nil)`

**`internal/app/app_service/authentication/forgot_password.go`**
- Add `audit.Log(ctx, app, "password_reset_request", &user.ID, &user.TenantID, nil)`

**`internal/app/app_service/authentication/reset_password.go`**
- Add `mutation.InvalidateAllUserTokens(ctx, db, userID)` after successful password update
- Add `audit.Log(ctx, app, "password_reset_complete", &reset.UserID, nil, nil)`

**`internal/app/mutation/password_reset_mutation.go`**
- Add `InvalidateAllUserTokens(ctx, db, userID)` function — sets `used_at = now()` on all active tokens for user

**Handler rename:** `get_api_v1_user.go` → `get_api_v1_me_user.go`
- Update operation ID to match `/me/user` REST pattern
- Update `openapi.yaml` path accordingly

---

## Frontend Changes

### New Files

**`components/GuestRoute.tsx`**
- Checks auth state; if user is authenticated, redirects to home route
- Wraps guest-only pages (signin, signup, forgot-password, reset-password)

**`components/PasswordRequirements.tsx`**
- Visual checklist showing password strength requirements
- Requirements: 8+ chars, uppercase, lowercase, digit, special character
- Real-time validation as user types

**`pages/ForgotPassword.tsx`**
- Email input form
- Calls `usePostApiV1ForgotPassword()` mutation
- Success state: "Check your email" message
- Link back to sign in

**`pages/ResetPassword.tsx`**
- Reads token from URL query params
- New password + confirm password fields
- Uses PasswordRequirements component
- Calls `usePostApiV1ResetPassword()` mutation
- Success: auto-redirect to signin after brief delay

**`pages/ErrorPage.tsx`**
- Generic error/404 page
- Styled with template branding
- Link back to home

**`services/api/axios-v1-admin.ts`**
- Custom fetch wrapper for `/api/admin/v1` prefix
- Same pattern as axios-v1.ts with credentials: "include"

**`services/api/v1-admin.ts`**
- Generated by Orval from openapi-admin.yaml

**`frontend/Dockerfile`**
- Multi-stage: Node 22 Alpine with Bun for build, Nginx Alpine for serve
- Copies built dist to nginx html dir
- Copies nginx.conf for SPA routing

**`frontend/nginx.conf`**
- SPA fallback: all routes → index.html
- 1-year cache with immutable for /assets/

### Updated Files

**`router.tsx`**
- Add GuestRoute wrapper for signin, signup, forgot-password, reset-password
- Add error element with ErrorPage
- `/` redirects to `/signin` for unauthenticated users

**`Routes.tsx`**
- Add route constants: forgotPassword, resetPassword

**`pages/SignIn.tsx`**
- Split layout: form on left, subtle gradient/pattern panel on right
- Add "Forgot password?" link
- Better input sizing (h-11)
- Improved error display

**`pages/SignUp.tsx`**
- Split layout with gradient panel
- Add PasswordRequirements component below password field
- Better input sizing

**`components/app-sidebar.tsx`**
- Use AuthContext to get user data
- Show user avatar (initials), name, email in footer
- Dropdown menu for logout
- Call `queryClient.clear()` on signout before redirect
- Remove static placeholder menu items

**`components/AppLayout.tsx`**
- Reduce header height to h-12
- Remove "Dashboard" title text, keep minimal

**`providers/QueryProvider.tsx`**
- Export `queryClient` instance so sidebar can clear cache on signout

**`index.css`**
- Import Inter font
- Add fade-in-up animation and animation delay utilities
- Enhanced card shadow/hover effects
- Richer CSS custom properties

**`orval.config.ts`**
- Add v1Admin API config pointing to openapi-admin.yaml

**`index.html`**
- Title: "SaaS Template"

**`package.json`**
- Add `@types/luxon` to devDependencies if missing

---

## Infrastructure Changes

### New Files

**`docker-compose.prod.yml`**
- **backend service:** image ref (genericized), healthcheck on `/api/public/v1/health`, migration entrypoint (`/goose -dir /migrations postgres "$DATABASE_URL" up` then `exec /server`), log rotation (10m max, 3 files), environment variables via `${VAR}` substitution
- **frontend service:** Nginx-based image, log rotation
- **nginx service:** Reverse proxy container, depends on backend (healthy) + frontend (started), port 80, mounts nginx config
- **dozzle service:** Log viewer on localhost:8888 only, mounts docker socket read-only
- **network:** IPv6-enabled with defined subnets
- Container names: saas-backend, saas-frontend, saas-nginx, saas-dozzle
- Image refs use placeholder format: `ghcr.io/your-org/saas-backend:latest`

**`nginx/nginx.conf`**
- Upstream blocks for backend (:8008) and frontend (:80)
- Server block on port 80, server_name localhost
- client_max_body_size 20m
- Security headers: X-Content-Type-Options nosniff, X-Frame-Options DENY, Referrer-Policy strict-origin-when-cross-origin, Permissions-Policy (deny camera/mic/geo)
- `/api/` location: proxy to backend with X-Real-IP, X-Forwarded-For, X-Forwarded-Proto headers, connect timeout 10s, read timeout 60s
- `/` location: proxy to frontend with same headers

---

## What Is NOT Ported

- Projects, placeholders, PDF generation (`fpdf`), Supabase storage provider
- Konva canvas, drag-and-drop (`@atlaskit/*`), CSV upload/preview, editor components
- Onboarding flow (WelcomeModal, OnboardingChecklist, CoachMark)
- Landing page, product mockup components
- `has_generated` field, `onboarding_dismissed` field
- Batch-specific migrations (00002_projects, 00006_pdf_generation, 00007_onboarding)
- `papaparse`, `konva`, `react-konva`, `recharts`, `motion` dependencies
- Analytics pages and analytics repository
- StamplLogo component (template uses generic branding)
