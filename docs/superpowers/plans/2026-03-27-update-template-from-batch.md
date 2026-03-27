# Update SaaS Template from Batch App — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Port general infrastructure improvements (security, audit logging, admin tier, production infra) from the batch app to the saas-template.

**Architecture:** Backend changes add an audit logging layer, admin middleware/routes, and security hardening to the existing layered architecture. Frontend changes improve auth pages, add password reset flows, and update the sidebar with user context. Infrastructure adds production Docker Compose with Nginx reverse proxy.

**Tech Stack:** Go 1.25, Chi v5, Jet ORM, oapi-codegen, React 19, TanStack React Query, shadcn/ui, Orval, Docker, Nginx

**Spec:** `docs/superpowers/specs/2026-03-27-update-template-from-batch-design.md`

---

## File Structure

### Backend — New Files
- `backend/db/migrations/00003_audit_log.sql` — Audit log table
- `backend/db/migrations/00004_user_role_constraint.sql` — CHECK constraint on role column
- `backend/internal/app/app_service/audit/log.go` — Audit logging service
- `backend/internal/app/mutation/audit_log_mutation.go` — Audit log DB insert
- `backend/internal/webserver/middleware/admin_middleware.go` — Admin role guard
- `backend/openapi-admin.yaml` — Admin OpenAPI spec (placeholder)
- `backend/generated/oapi/admin/codegen.yaml` — Admin codegen config
- `backend/.dockerignore` — Docker build exclusions

### Backend — Modified Files
- `backend/config/provider/session_provider.go` — Add cookie security options
- `backend/internal/webserver/middleware/logger.go` — Generic error message
- `backend/internal/webserver/webserver.go` — Add admin route group + rate limiter
- `backend/internal/app/app_service/authentication/signup.go` — Secure error + audit log
- `backend/internal/app/app_service/authentication/signin.go` — Audit log calls
- `backend/internal/app/app_service/authentication/signout.go` — Audit log call
- `backend/internal/app/app_service/authentication/forgot_password.go` — Audit log call
- `backend/internal/app/app_service/authentication/reset_password.go` — Audit log + token invalidation
- `backend/internal/webserver/handler/baseHandler.go` — Add admin interface
- `backend/commands.sh` — Add admin codegen step
- `backend/.golangci.yml` — Already matches batch (no change needed)

### Frontend — New Files
- `frontend/src/components/GuestRoute.tsx` — Auth guard for guest pages
- `frontend/src/components/PasswordRequirements.tsx` — Password strength indicator
- `frontend/src/pages/ForgotPassword.tsx` — Forgot password page
- `frontend/src/pages/ResetPassword.tsx` — Reset password page
- `frontend/src/pages/ErrorPage.tsx` — Error/404 page
- `frontend/src/services/api/axios-v1-admin.ts` — Admin API fetch wrapper
- `frontend/Dockerfile` — Multi-stage production build
- `frontend/nginx.conf` — SPA fallback routing

### Frontend — Modified Files
- `frontend/src/router.tsx` — Complete rewrite with guest routes, error handling
- `frontend/src/Routes.tsx` — Add all route constants
- `frontend/src/pages/SignIn.tsx` — Split layout with gradient panel
- `frontend/src/pages/SignUp.tsx` — Split layout with gradient panel + password requirements
- `frontend/src/components/app-sidebar.tsx` — User context, avatar, dropdown logout
- `frontend/src/components/AppLayout.tsx` — Minimal h-12 header
- `frontend/src/providers/QueryProvider.tsx` — Export queryClient
- `frontend/src/index.css` — Enhanced theming
- `frontend/src/contexts/authContext.tsx` — Simplified auth check
- `frontend/orval.config.ts` — Add admin API config
- `frontend/index.html` — Update title
- `frontend/package.json` — Add missing deps

### Infrastructure — New Files
- `docker-compose.prod.yml` — Production compose
- `nginx/nginx.conf` — Reverse proxy config

---

### Task 1: Backend — Add Audit Log Migration and Mutation

**Files:**
- Create: `backend/db/migrations/00003_audit_log.sql`
- Create: `backend/internal/app/mutation/audit_log_mutation.go`

- [ ] **Step 1: Create audit_log migration**

Create `backend/db/migrations/00003_audit_log.sql`:

```sql
-- +goose Up
-- +goose StatementBegin
CREATE TABLE audit_log_tbl (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    action VARCHAR(50) NOT NULL,
    actor_id UUID,
    tenant_id UUID,
    ip_address VARCHAR(45),
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_audit_log_action ON audit_log_tbl(action);
CREATE INDEX idx_audit_log_actor ON audit_log_tbl(actor_id);
CREATE INDEX idx_audit_log_created ON audit_log_tbl(created_at);
CREATE INDEX idx_audit_log_tenant_created ON audit_log_tbl(tenant_id, created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS audit_log_tbl;
-- +goose StatementEnd
```

- [ ] **Step 2: Create role CHECK constraint migration**

Create `backend/db/migrations/00004_user_role_constraint.sql`:

```sql
-- +goose Up
ALTER TABLE user_tbl ADD CONSTRAINT chk_user_role CHECK (role IN ('admin', 'user'));

-- +goose Down
ALTER TABLE user_tbl DROP CONSTRAINT chk_user_role;
```

Note: The template already has the `role` column in its init migration. This just adds the constraint.

- [ ] **Step 3: Create audit_log_mutation.go**

Create `backend/internal/app/mutation/audit_log_mutation.go`:

```go
package mutation

import (
	"context"
	"saas-template/generated/db/database/public/model"
	"saas-template/generated/db/database/public/table"
	"time"

	"github.com/go-jet/jet/v2/qrm"
)

func CreateAuditLog(
	ctx context.Context,
	db qrm.DB,
	auditLog model.AuditLogTbl,
) error {
	ctbl := table.AuditLogTbl

	auditLog.CreatedAt = time.Now()

	stmt := ctbl.INSERT(ctbl.MutableColumns).
		MODEL(auditLog)

	_, err := stmt.ExecContext(ctx, db)
	return err
}
```

- [ ] **Step 4: Commit**

```bash
git add backend/db/migrations/00003_audit_log.sql backend/db/migrations/00004_user_role_constraint.sql backend/internal/app/mutation/audit_log_mutation.go
git commit -m "feat: add audit log migration and mutation"
```

---

### Task 2: Backend — Add Audit Service

**Files:**
- Create: `backend/internal/app/app_service/audit/log.go`

- [ ] **Step 1: Create audit log service**

Create `backend/internal/app/app_service/audit/log.go`:

```go
package audit

import (
	"context"
	"encoding/json"
	"log/slog"
	"net"
	"saas-template/config"
	"saas-template/generated/db/database/public/model"
	"saas-template/internal/app/mutation"
	"saas-template/internal/webserver/middleware"
	"strings"

	"github.com/google/uuid"
)

func Log(
	ctx context.Context,
	app *config.App,
	action string,
	actorID *uuid.UUID,
	tenantID *uuid.UUID,
	metadata map[string]string,
) {
	ip := extractIP(ctx)

	auditLog := model.AuditLogTbl{
		Action:    action,
		ActorID:   actorID,
		TenantID:  tenantID,
		IPAddress: ip,
	}

	if metadata != nil {
		b, err := json.Marshal(metadata)
		if err == nil {
			s := string(b)
			auditLog.Metadata = &s
		}
	}

	err := mutation.CreateAuditLog(ctx, app.DB(), auditLog)
	if err != nil {
		slog.Default().
			ErrorContext(ctx, "failed to create audit log", "error", err, "action", action)
	}
}

func extractIP(ctx context.Context) *string {
	r := middleware.GetRequest(ctx)
	if r == nil {
		return nil
	}

	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ip := strings.TrimSpace(strings.Split(xff, ",")[0])
		return &ip
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		addr := r.RemoteAddr
		return &addr
	}
	return &ip
}
```

- [ ] **Step 2: Commit**

```bash
git add backend/internal/app/app_service/audit/log.go
git commit -m "feat: add audit logging service"
```

---

### Task 3: Backend — Security Hardening

**Files:**
- Modify: `backend/config/provider/session_provider.go`
- Modify: `backend/internal/webserver/middleware/logger.go`
- Modify: `backend/internal/app/app_service/authentication/signup.go`

- [ ] **Step 1: Add cookie security to session provider**

Replace the entire content of `backend/config/provider/session_provider.go`:

```go
package provider

import (
	"net/http"

	"github.com/gorilla/sessions"
)

func NewSessionProvider(env *EnvProvider) *sessions.CookieStore {
	store := sessions.NewCookieStore([]byte(env.SessionSecret()))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 30, // 30 days
		HttpOnly: true,
		Secure:   env.AppEnv() != AppEnvLocal,
		SameSite: http.SameSiteLaxMode,
	}
	return store
}
```

- [ ] **Step 2: Fix error leaking in logger**

In `backend/internal/webserver/middleware/logger.go`, in the `HandleErrorWithLog` function, change:

```go
http.Error(w, err.Error(), http.StatusInternalServerError)
```

to:

```go
http.Error(w, "Internal server error", http.StatusInternalServerError)
```

- [ ] **Step 3: Fix signup error message**

In `backend/internal/app/app_service/authentication/signup.go`, change:

```go
return nil, errors.New("a user with this email already exists")
```

to:

```go
return nil, errors.New("unable to create account")
```

- [ ] **Step 4: Commit**

```bash
git add backend/config/provider/session_provider.go backend/internal/webserver/middleware/logger.go backend/internal/app/app_service/authentication/signup.go
git commit -m "fix: harden session cookies, prevent error/email leaking"
```

---

### Task 4: Backend — Add Audit Logging to Auth Flows

**Files:**
- Modify: `backend/internal/app/app_service/authentication/signup.go`
- Modify: `backend/internal/app/app_service/authentication/signin.go`
- Modify: `backend/internal/app/app_service/authentication/signout.go`
- Modify: `backend/internal/app/app_service/authentication/forgot_password.go`
- Modify: `backend/internal/app/app_service/authentication/reset_password.go`

- [ ] **Step 1: Add audit logging to signup**

In `backend/internal/app/app_service/authentication/signup.go`:

Add to imports:
```go
"saas-template/internal/app/app_service/audit"
```

After the `created, err := mutation.CreateUser(...)` block and before `return created, nil`, add:

```go
audit.Log(ctx, app, "signup", &created.ID, &created.TenantID, nil)
```

- [ ] **Step 2: Add audit logging to signin**

Replace entire content of `backend/internal/app/app_service/authentication/signin.go`:

```go
package authentication

import (
	"context"
	"errors"
	"saas-template/config"
	"saas-template/generated/db/database/public/model"
	"saas-template/internal/app/app_service/audit"
	"saas-template/internal/app/repository"

	"golang.org/x/crypto/bcrypt"
)

func SignIn(
	ctx context.Context,
	app *config.App,
	email string,
	passwordInput string,
) (*model.UserTbl, error) {
	user, err := repository.GetUserByEmail(ctx, app.DB(), email)
	if err != nil {
		audit.Log(ctx, app, "signin_failed", nil, nil, map[string]string{"email": email})
		return nil, errors.New("invalid email or password")
	}

	if user.PasswordHash == nil {
		audit.Log(
			ctx,
			app,
			"signin_failed",
			&user.ID,
			&user.TenantID,
			map[string]string{"email": email},
		)
		return nil, errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(*user.PasswordHash),
		[]byte(passwordInput),
	); err != nil {
		audit.Log(
			ctx,
			app,
			"signin_failed",
			&user.ID,
			&user.TenantID,
			map[string]string{"email": email},
		)
		return nil, errors.New("invalid email or password")
	}

	audit.Log(ctx, app, "signin_success", &user.ID, &user.TenantID, nil)
	return user, nil
}
```

- [ ] **Step 3: Add audit logging to signout**

Replace entire content of `backend/internal/app/app_service/authentication/signout.go`:

```go
package authentication

import (
	"context"
	"saas-template/config"
	"saas-template/internal/app/app_service/audit"

	"github.com/google/uuid"
)

func SignOut(ctx context.Context, app *config.App, userID uuid.UUID, tenantID uuid.UUID) {
	audit.Log(ctx, app, "signout", &userID, &tenantID, nil)
}
```

- [ ] **Step 4: Add audit logging to forgot_password**

In `backend/internal/app/app_service/authentication/forgot_password.go`:

Add to imports:
```go
"saas-template/internal/app/app_service/audit"
```

At the end of the function, just before the closing `}`, add:

```go
audit.Log(ctx, app, "password_reset_request", &user.ID, &user.TenantID, nil)
```

- [ ] **Step 5: Add audit logging to reset_password**

In `backend/internal/app/app_service/authentication/reset_password.go`:

Add to imports:
```go
"saas-template/internal/app/app_service/audit"
```

After the `InvalidateAllUserTokens` call and before `return nil, nil`, add:

```go
audit.Log(ctx, app, "password_reset_complete", &resetRecord.UserID, nil, nil)
```

- [ ] **Step 6: Commit**

```bash
git add backend/internal/app/app_service/authentication/
git commit -m "feat: add audit logging to all auth flows"
```

---

### Task 5: Backend — Add Admin Middleware, Route Group, and OpenAPI Spec

**Files:**
- Create: `backend/internal/webserver/middleware/admin_middleware.go`
- Create: `backend/openapi-admin.yaml`
- Create: `backend/generated/oapi/admin/codegen.yaml`
- Create: `backend/.dockerignore`
- Modify: `backend/internal/webserver/webserver.go`
- Modify: `backend/internal/webserver/handler/baseHandler.go`
- Modify: `backend/commands.sh`

- [ ] **Step 1: Create admin middleware**

Create `backend/internal/webserver/middleware/admin_middleware.go`:

```go
package middleware

import (
	"net/http"
	"saas-template/config"
)

func NewAdminMiddleware(app *config.App) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, _ := app.Session().Get(r, "session")

			if role, ok := session.Values["role"].(string); ok && role == "admin" {
				next.ServeHTTP(w, r)
			} else {
				http.Error(w, "Forbidden", http.StatusForbidden)
			}
		})
	}
}
```

- [ ] **Step 2: Create admin OpenAPI spec**

Create `backend/openapi-admin.yaml`:

```yaml
openapi: 3.0.0
info:
  title: "SaaS Template Admin API"
  version: "1.0.0"
paths:
  /health:
    get:
      operationId: get-api-admin-v1-health
      summary: "Admin health check"
      responses:
        "200":
          description: "OK"
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/MessageResponse"

components:
  schemas:
    MessageResponse:
      type: object
      properties:
        message:
          type: string
```

- [ ] **Step 3: Create admin codegen config**

Create directory and file `backend/generated/oapi/admin/codegen.yaml`:

```yaml
package: oapi
generate:
  chi-server: true
  models: true
  client: true
  strict-server: true
output: generated/oapi/admin/generated.go
output-options:
  skip-prune: true
```

- [ ] **Step 4: Create .dockerignore**

Create `backend/.dockerignore`:

```
.env*
.air.toml
.golangci.yml
vendor/
*.md
```

- [ ] **Step 5: Update baseHandler.go to include admin interface**

Replace entire content of `backend/internal/webserver/handler/baseHandler.go`:

```go
package handler

import (
	"saas-template/config"
	"saas-template/generated/oapi"
	oapi_admin "saas-template/generated/oapi/admin"
	oapi_public "saas-template/generated/oapi/public"
)

type CombinedStrictServerInterface interface {
	oapi.StrictServerInterface
	oapi_public.StrictServerInterface
	oapi_admin.StrictServerInterface
}

type Handler struct {
	app *config.App
	CombinedStrictServerInterface
}

func NewHandler(app *config.App) *Handler {
	handler := Handler{
		app: app,
	}

	return &handler
}
```

- [ ] **Step 6: Create admin health handler**

Create `backend/internal/webserver/handler/get_api_admin_v1_health.go`:

```go
package handler

import (
	"context"
	oapi_admin "saas-template/generated/oapi/admin"
)

func (h *Handler) GetApiAdminV1Health(
	ctx context.Context,
	request oapi_admin.GetApiAdminV1HealthRequestObject,
) (oapi_admin.GetApiAdminV1HealthResponseObject, error) {
	return oapi_admin.GetApiAdminV1Health200JSONResponse{
		Message: ptr("ok"),
	}, nil
}

func ptr(s string) *string { return &s }
```

- [ ] **Step 7: Update webserver.go with admin route group**

Replace entire content of `backend/internal/webserver/webserver.go`:

```go
package webserver

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"saas-template/config"
	"saas-template/generated/oapi"
	oapi_admin "saas-template/generated/oapi/admin"
	oapi_public "saas-template/generated/oapi/public"
	"saas-template/internal/webserver/handler"
	"saas-template/internal/webserver/middleware"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type Webserver struct {
	router     *chi.Mux
	serverAddr string
}

func (ws *Webserver) Router() *chi.Mux {
	return ws.router
}

func NewWebserver(app *config.App) *Webserver {
	handler := handler.NewHandler(app)
	serverAddr := ":" + app.EnvVars().ServerPort()

	r := chi.NewRouter()
	r.Use(middleware.NewLoggerMiddleware())
	r.Use(chimiddleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{app.EnvVars().AppBaseUrl()},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	// Public API routes (rate limited)
	authRateLimiter := middleware.NewRateLimiter(1, 5) // 1 req/sec, burst of 5
	r.Group(func(r chi.Router) {
		r.Use(authRateLimiter.Middleware())
		r.Use(middleware.NewContextInjectorMiddleware())
		baseURL := "/api/public/v1"
		strictHandler := oapi_public.NewStrictHandlerWithOptions(
			handler,
			[]oapi_public.StrictMiddlewareFunc{},
			oapi_public.StrictHTTPServerOptions{
				RequestErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
					http.Error(w, err.Error(), http.StatusBadRequest)
				},
				ResponseErrorHandlerFunc: middleware.HandleErrorWithLog(app),
			},
		)
		oapi_public.HandlerFromMuxWithBaseURL(strictHandler, r, baseURL)
	})

	// Authenticated API routes (rate limited)
	apiRateLimiter := middleware.NewRateLimiter(10, 20) // 10 req/sec, burst of 20
	r.Group(func(r chi.Router) {
		r.Use(apiRateLimiter.Middleware())
		r.Use(middleware.NewContextInjectorMiddleware())
		r.Use(middleware.NewAuthMiddleware(app))
		baseURL := "/api/v1"
		strictHandler := oapi.NewStrictHandlerWithOptions(
			handler,
			[]oapi.StrictMiddlewareFunc{},
			oapi.StrictHTTPServerOptions{
				RequestErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
					http.Error(w, err.Error(), http.StatusBadRequest)
				},
				ResponseErrorHandlerFunc: middleware.HandleErrorWithLog(app),
			},
		)
		oapi.HandlerFromMuxWithBaseURL(strictHandler, r, baseURL)
	})

	// Admin API routes (rate limited, authenticated + admin role required)
	adminRateLimiter := middleware.NewRateLimiter(5, 10) // 5 req/sec, burst of 10
	r.Group(func(r chi.Router) {
		r.Use(adminRateLimiter.Middleware())
		r.Use(middleware.NewContextInjectorMiddleware())
		r.Use(middleware.NewAuthMiddleware(app))
		r.Use(middleware.NewAdminMiddleware(app))
		baseURL := "/api/admin/v1"
		strictHandler := oapi_admin.NewStrictHandlerWithOptions(
			handler,
			[]oapi_admin.StrictMiddlewareFunc{},
			oapi_admin.StrictHTTPServerOptions{
				RequestErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
					http.Error(w, err.Error(), http.StatusBadRequest)
				},
				ResponseErrorHandlerFunc: middleware.HandleErrorWithLog(app),
			},
		)
		oapi_admin.HandlerFromMuxWithBaseURL(strictHandler, r, baseURL)
	})

	return &Webserver{
		router:     r,
		serverAddr: serverAddr,
	}
}

func (ws *Webserver) Start() {
	s := &http.Server{
		Handler: ws.router,
		Addr:    ws.serverAddr,
	}

	go func() {
		log.Print("WebServer listening on " + ws.serverAddr)
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Print("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		log.Print("Server forced to shutdown:", err)
	}

	log.Print("Server exited")
}

func (ws *Webserver) PrintRoutes() {
	err := chi.Walk(
		ws.router,
		func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
			fmt.Printf("[%s]: '%s' has %d middlewares\n", method, route, len(middlewares))
			return nil
		},
	)

	if err != nil {
		log.Panicln(err)
	}
}
```

- [ ] **Step 8: Update commands.sh with admin codegen**

In `backend/commands.sh`, in the `openapi:codegen` case, after the line:

```bash
go tool oapi-codegen -config ./generated/oapi/public/codegen.yaml ./openapi-public.yaml
```

Add:

```bash
  echo -e "${GREEN}Generating OpenAPI admin code...${NC}"
  go tool oapi-codegen -config ./generated/oapi/admin/codegen.yaml ./openapi-admin.yaml
```

- [ ] **Step 9: Run OpenAPI codegen to generate admin types**

```bash
cd /Users/greg/dev/projects/saas-template/backend
mkdir -p generated/oapi/admin
go tool oapi-codegen -config ./generated/oapi/admin/codegen.yaml ./openapi-admin.yaml
```

- [ ] **Step 10: Verify the backend compiles**

```bash
cd /Users/greg/dev/projects/saas-template/backend
go build ./...
```

Expected: Successful build with no errors.

- [ ] **Step 11: Commit**

```bash
git add backend/internal/webserver/middleware/admin_middleware.go backend/openapi-admin.yaml backend/generated/oapi/admin/ backend/.dockerignore backend/internal/webserver/webserver.go backend/internal/webserver/handler/baseHandler.go backend/internal/webserver/handler/get_api_admin_v1_health.go backend/commands.sh
git commit -m "feat: add admin middleware, route group, and OpenAPI spec"
```

---

### Task 6: Frontend — Add Shared Components and Dependencies

**Files:**
- Create: `frontend/src/components/GuestRoute.tsx`
- Create: `frontend/src/components/PasswordRequirements.tsx`
- Create: `frontend/src/services/api/axios-v1-admin.ts`
- Modify: `frontend/src/providers/QueryProvider.tsx`
- Modify: `frontend/package.json`

- [ ] **Step 1: Install missing shadcn components and dependencies**

```bash
cd /Users/greg/dev/projects/saas-template/frontend
bunx shadcn@latest add avatar dropdown-menu
bun add luxon
bun add -D @types/luxon
```

- [ ] **Step 2: Create GuestRoute component**

Create `frontend/src/components/GuestRoute.tsx`:

```tsx
import { type PropsWithChildren } from "react";
import { Navigate } from "react-router";
import { useGetApiV1User } from "@/services/api/v1";

export function GuestRoute({ children }: PropsWithChildren) {
  const { data, isLoading } = useGetApiV1User({ query: { retry: false } });

  if (isLoading) return null;

  if (data) {
    return <Navigate to="/" replace />;
  }

  return <>{children}</>;
}
```

- [ ] **Step 3: Create PasswordRequirements component**

Create `frontend/src/components/PasswordRequirements.tsx`:

```tsx
import { cn } from "@/lib/utils";

interface PasswordRequirementsProps {
  password: string;
}

const requirements = [
  { label: "At least 8 characters", test: (p: string) => p.length >= 8 },
  { label: "1 uppercase letter", test: (p: string) => /[A-Z]/.test(p) },
  { label: "1 lowercase letter", test: (p: string) => /[a-z]/.test(p) },
  { label: "1 number", test: (p: string) => /\d/.test(p) },
  {
    label: "1 special character",
    test: (p: string) => /[^a-zA-Z0-9\s]/.test(p),
  },
];

export function PasswordRequirements({ password }: PasswordRequirementsProps) {
  if (!password) return null;

  return (
    <ul className="space-y-1 text-xs text-muted-foreground">
      {requirements.map((req) => {
        const met = req.test(password);
        return (
          <li
            key={req.label}
            className={cn(
              "flex items-center gap-1.5 transition-colors",
              met && "text-emerald-600"
            )}
          >
            <span className="text-[10px]">{met ? "\u2713" : "\u2022"}</span>
            {req.label}
          </li>
        );
      })}
    </ul>
  );
}
```

- [ ] **Step 4: Create admin API fetch wrapper**

Create `frontend/src/services/api/axios-v1-admin.ts`:

```tsx
export const customInstance = async <T>(
  url: string,
  options?: RequestInit,
): Promise<T> => {
  const response = await fetch(`/api/admin/v1${url}`, {
    credentials: "include",
    ...options,
  });

  if (!response.ok) {
    throw response;
  }

  return response.json() as Promise<T>;
};

export type ErrorType<Error> = Error;
export type BodyType<BodyData> = BodyData;
```

- [ ] **Step 5: Export queryClient from QueryProvider**

Replace entire content of `frontend/src/providers/QueryProvider.tsx`:

```tsx
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import type { JSX, PropsWithChildren } from "react";

export const queryClient = new QueryClient();

export const QueryProvider = ({ children }: PropsWithChildren): JSX.Element => (
  <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
);
```

- [ ] **Step 6: Update orval.config.ts with admin API**

Replace entire content of `frontend/orval.config.ts`:

```typescript
import { defineConfig } from "orval";

export default defineConfig({
  v1: {
    input: "../backend/openapi.yaml",
    output: {
      target: "./src/services/api/v1.ts",
      client: "react-query",
      httpClient: "fetch",
      mode: "single",
      override: {
        fetch: {
          includeHttpResponseReturnType: false,
        },
        mutator: {
          path: "./src/services/api/axios-v1.ts",
          name: "customInstance",
        },
      },
    },
  },
  v1Public: {
    input: "../backend/openapi-public.yaml",
    output: {
      target: "./src/services/api/v1-public.ts",
      client: "react-query",
      httpClient: "fetch",
      mode: "single",
      override: {
        fetch: {
          includeHttpResponseReturnType: false,
        },
        mutator: {
          path: "./src/services/api/axios-v1-public.ts",
          name: "customInstance",
        },
      },
    },
  },
  v1Admin: {
    input: "../backend/openapi-admin.yaml",
    output: {
      target: "./src/services/api/v1-admin.ts",
      client: "react-query",
      httpClient: "fetch",
      mode: "single",
      override: {
        fetch: {
          includeHttpResponseReturnType: false,
        },
        mutator: {
          path: "./src/services/api/axios-v1-admin.ts",
          name: "customInstance",
        },
      },
    },
  },
});
```

- [ ] **Step 7: Regenerate Orval types**

```bash
cd /Users/greg/dev/projects/saas-template/frontend
bunx orval
```

- [ ] **Step 8: Commit**

```bash
git add frontend/src/components/GuestRoute.tsx frontend/src/components/PasswordRequirements.tsx frontend/src/services/api/axios-v1-admin.ts frontend/src/providers/QueryProvider.tsx frontend/orval.config.ts frontend/src/services/api/v1-admin.ts frontend/src/services/api/v1.ts frontend/src/services/api/v1-public.ts frontend/package.json frontend/bun.lock frontend/src/components/ui/avatar.tsx frontend/src/components/ui/dropdown-menu.tsx
git commit -m "feat: add shared frontend components and admin API config"
```

---

### Task 7: Frontend — Add New Pages (ForgotPassword, ResetPassword, ErrorPage)

**Files:**
- Create: `frontend/src/pages/ForgotPassword.tsx`
- Create: `frontend/src/pages/ResetPassword.tsx`
- Create: `frontend/src/pages/ErrorPage.tsx`

- [ ] **Step 1: Create ForgotPassword page**

Create `frontend/src/pages/ForgotPassword.tsx`:

```tsx
import { useState } from "react";
import { Link } from "react-router";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { usePostApiV1ForgotPassword } from "@/services/api/v1-public";
import { Loader2, ArrowLeft, MailCheck } from "lucide-react";

export function ForgotPassword() {
  const [email, setEmail] = useState("");
  const [submitted, setSubmitted] = useState(false);
  const forgotPassword = usePostApiV1ForgotPassword();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await forgotPassword.mutateAsync({ data: { email } });
    } catch {
      // Always show success to prevent email enumeration
    }
    setSubmitted(true);
  };

  return (
    <div className="flex min-h-screen">
      {/* Left — Form */}
      <div className="flex w-full items-center justify-center px-6 py-12 lg:w-1/2 lg:shrink-0">
        <div className="w-full max-w-sm">
          {submitted ? (
            <>
              <div className="mb-8">
                <div className="mb-4 flex size-12 items-center justify-center rounded-xl bg-primary/10">
                  <MailCheck className="size-6 text-primary" />
                </div>
                <h1 className="text-2xl font-bold tracking-tight text-foreground">
                  Check your email
                </h1>
                <p className="mt-1.5 text-sm text-muted-foreground">
                  If an account with that email exists, we've sent a password
                  reset link. Check your inbox.
                </p>
              </div>

              <Button asChild variant="outline" className="h-11 w-full text-sm font-medium">
                <Link to="/signin">
                  <ArrowLeft className="size-4" />
                  Back to sign in
                </Link>
              </Button>
            </>
          ) : (
            <>
              <div className="mb-8">
                <h1 className="text-2xl font-bold tracking-tight text-foreground">
                  Forgot password?
                </h1>
                <p className="mt-1.5 text-sm text-muted-foreground">
                  Enter your email and we'll send you a reset link
                </p>
              </div>

              <form onSubmit={handleSubmit} className="space-y-5">
                <div className="space-y-2">
                  <Label htmlFor="email" className="text-sm font-medium">
                    Email
                  </Label>
                  <Input
                    id="email"
                    type="email"
                    placeholder="name@example.com"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    required
                    autoFocus
                    className="h-11"
                  />
                </div>

                <Button
                  type="submit"
                  className="h-11 w-full text-sm font-medium"
                  disabled={forgotPassword.isPending}
                >
                  {forgotPassword.isPending ? (
                    <>
                      <Loader2 className="mr-2 size-4 animate-spin" />
                      Sending...
                    </>
                  ) : (
                    "Send Reset Link"
                  )}
                </Button>
              </form>

              <p className="mt-8 text-center text-sm text-muted-foreground">
                Remember your password?{" "}
                <Link
                  to="/signin"
                  className="font-medium text-primary transition-colors hover:text-primary/80"
                >
                  Sign in
                </Link>
              </p>
            </>
          )}
        </div>
      </div>

      {/* Right — Gradient Panel */}
      <div className="hidden lg:flex lg:w-1/2 lg:items-center lg:justify-center bg-gradient-to-br from-primary/5 via-primary/10 to-primary/5">
        <div className="h-64 w-64 rounded-full bg-primary/5" />
      </div>
    </div>
  );
}
```

- [ ] **Step 2: Create ResetPassword page**

Create `frontend/src/pages/ResetPassword.tsx`:

```tsx
import { useState } from "react";
import { Link, useSearchParams, useNavigate } from "react-router";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { PasswordRequirements } from "@/components/PasswordRequirements";
import { usePostApiV1ResetPassword } from "@/services/api/v1-public";
import { Loader2, CheckCircle2, AlertCircle } from "lucide-react";

export function ResetPassword() {
  const [searchParams] = useSearchParams();
  const token = searchParams.get("token");
  const navigate = useNavigate();

  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [error, setError] = useState("");
  const [success, setSuccess] = useState(false);
  const resetPassword = usePostApiV1ResetPassword();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");

    if (password !== confirmPassword) {
      setError("Passwords don't match");
      return;
    }

    try {
      await resetPassword.mutateAsync({
        data: { token: token!, password, confirm_password: confirmPassword },
      });
      setSuccess(true);
      setTimeout(() => navigate("/signin"), 3000);
    } catch {
      setError("Invalid or expired reset token. Please request a new one.");
    }
  };

  return (
    <div className="flex min-h-screen">
      {/* Left — Form */}
      <div className="flex w-full items-center justify-center px-6 py-12 lg:w-1/2 lg:shrink-0">
        <div className="w-full max-w-sm">
          {!token ? (
            <>
              <div className="mb-8">
                <div className="mb-4 flex size-12 items-center justify-center rounded-xl bg-destructive/10">
                  <AlertCircle className="size-6 text-destructive" />
                </div>
                <h1 className="text-2xl font-bold tracking-tight text-foreground">
                  Invalid reset link
                </h1>
                <p className="mt-1.5 text-sm text-muted-foreground">
                  This password reset link is invalid or has expired.
                </p>
              </div>

              <Button asChild className="h-11 w-full text-sm font-medium">
                <Link to="/forgot-password">Request a new reset link</Link>
              </Button>
            </>
          ) : success ? (
            <>
              <div className="mb-8">
                <div className="mb-4 flex size-12 items-center justify-center rounded-xl bg-emerald-500/10">
                  <CheckCircle2 className="size-6 text-emerald-600" />
                </div>
                <h1 className="text-2xl font-bold tracking-tight text-foreground">
                  Password reset
                </h1>
                <p className="mt-1.5 text-sm text-muted-foreground">
                  Your password has been reset successfully. Redirecting to sign
                  in...
                </p>
              </div>

              <Button asChild variant="outline" className="h-11 w-full text-sm font-medium">
                <Link to="/signin">Go to sign in</Link>
              </Button>
            </>
          ) : (
            <>
              <div className="mb-8">
                <h1 className="text-2xl font-bold tracking-tight text-foreground">
                  Reset password
                </h1>
                <p className="mt-1.5 text-sm text-muted-foreground">
                  Enter your new password below
                </p>
              </div>

              <form onSubmit={handleSubmit} className="space-y-5">
                {error && (
                  <div className="rounded-lg border border-destructive/20 bg-destructive/5 px-4 py-3 text-sm text-destructive">
                    {error}
                  </div>
                )}

                <div className="space-y-2">
                  <Label htmlFor="password" className="text-sm font-medium">
                    New Password
                  </Label>
                  <Input
                    id="password"
                    type="password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    required
                    minLength={8}
                    autoFocus
                    className="h-11"
                  />
                  <PasswordRequirements password={password} />
                </div>

                <div className="space-y-2">
                  <Label htmlFor="confirm-password" className="text-sm font-medium">
                    Confirm Password
                  </Label>
                  <Input
                    id="confirm-password"
                    type="password"
                    value={confirmPassword}
                    onChange={(e) => setConfirmPassword(e.target.value)}
                    required
                    className="h-11"
                  />
                </div>

                <Button
                  type="submit"
                  className="h-11 w-full text-sm font-medium"
                  disabled={resetPassword.isPending}
                >
                  {resetPassword.isPending ? (
                    <>
                      <Loader2 className="mr-2 size-4 animate-spin" />
                      Resetting...
                    </>
                  ) : (
                    "Reset Password"
                  )}
                </Button>
              </form>
            </>
          )}
        </div>
      </div>

      {/* Right — Gradient Panel */}
      <div className="hidden lg:flex lg:w-1/2 lg:items-center lg:justify-center bg-gradient-to-br from-primary/5 via-primary/10 to-primary/5">
        <div className="h-64 w-64 rounded-full bg-primary/5" />
      </div>
    </div>
  );
}
```

- [ ] **Step 3: Create ErrorPage**

Create `frontend/src/pages/ErrorPage.tsx`:

```tsx
import { useRouteError, isRouteErrorResponse, Link } from "react-router";
import { Button } from "@/components/ui/button";

export function ErrorPage() {
  const error = useRouteError();
  const isNotFound = isRouteErrorResponse(error) && error.status === 404;

  if (isNotFound) {
    return (
      <div className="flex min-h-screen flex-col items-center justify-center px-6">
        <span className="text-[10rem] leading-none font-black tracking-tighter text-foreground/10">
          404
        </span>

        <h1 className="mt-2 text-2xl font-bold tracking-tight text-foreground">
          Page not found
        </h1>
        <p className="mt-2 text-sm text-muted-foreground">
          The page you're looking for doesn't exist or has been moved.
        </p>

        <Button asChild className="mt-8">
          <Link to="/">Go Home</Link>
        </Button>
      </div>
    );
  }

  return (
    <div className="flex min-h-screen flex-col items-center justify-center px-6">
      <h1 className="mt-6 text-2xl font-bold tracking-tight text-foreground">
        Something went wrong
      </h1>
      <p className="mt-2 text-sm text-muted-foreground">
        An unexpected error occurred. Please try again.
      </p>

      <div className="mt-8 flex gap-3">
        <Button variant="outline" onClick={() => window.location.reload()}>
          Refresh
        </Button>
        <Button asChild>
          <Link to="/">Go Home</Link>
        </Button>
      </div>
    </div>
  );
}
```

- [ ] **Step 4: Commit**

```bash
git add frontend/src/pages/ForgotPassword.tsx frontend/src/pages/ResetPassword.tsx frontend/src/pages/ErrorPage.tsx
git commit -m "feat: add forgot password, reset password, and error pages"
```

---

### Task 8: Frontend — Update Router, Auth Context, and Routes

**Files:**
- Modify: `frontend/src/router.tsx`
- Modify: `frontend/src/Routes.tsx`
- Modify: `frontend/src/contexts/authContext.tsx`

- [ ] **Step 1: Update Routes.tsx**

Replace entire content of `frontend/src/Routes.tsx`:

```tsx
export const Routes = {
  home: "/",
  forgotPassword: "/forgot-password",
  resetPassword: "/reset-password",
} as const;

export const UnauthorizedRoutes = {
  signin: "/signin",
  signup: "/signup",
  forgotPassword: "/forgot-password",
  resetPassword: "/reset-password",
} as const;

export type UnauthorizedRoutes =
  (typeof UnauthorizedRoutes)[keyof typeof UnauthorizedRoutes];

export default Routes;
```

- [ ] **Step 2: Update authContext.tsx**

Replace entire content of `frontend/src/contexts/authContext.tsx`:

```tsx
import { createContext, type JSX, type PropsWithChildren } from "react";
import { Navigate } from "react-router";
import { useGetApiV1User, type BaseUser } from "../services/api/v1";

export type AuthContextInterface = {
  user: BaseUser | null;
};

export const AuthContext = createContext<AuthContextInterface>({ user: null });

export const AuthContextProvider = (props: PropsWithChildren): JSX.Element => {
  const { data, error, isLoading } = useGetApiV1User({ query: { retry: false } });

  if (isLoading) return <></>;

  if (error) {
    return <Navigate to="/signin" replace />;
  }

  const authCtx: AuthContextInterface = {
    user: data ?? null,
  };

  return (
    <AuthContext.Provider value={authCtx}>
      {props.children}
    </AuthContext.Provider>
  );
};
```

- [ ] **Step 3: Update router.tsx**

Replace entire content of `frontend/src/router.tsx`:

```tsx
import { createBrowserRouter, Navigate } from "react-router";
import { AppLayout } from "./components/AppLayout";
import { AuthContextProvider } from "./contexts/authContext";
import { GuestRoute } from "./components/GuestRoute";
import { SignIn } from "./pages/SignIn";
import { SignUp } from "./pages/SignUp";
import { ForgotPassword } from "./pages/ForgotPassword";
import { ResetPassword } from "./pages/ResetPassword";
import { ErrorPage } from "./pages/ErrorPage";

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
            element: (
              <div className="text-muted-foreground">
                Welcome to your SaaS app. Select a page from the sidebar.
              </div>
            ),
          },
        ],
      },
      {
        path: "/signin",
        element: (
          <GuestRoute>
            <SignIn />
          </GuestRoute>
        ),
      },
      {
        path: "/signup",
        element: (
          <GuestRoute>
            <SignUp />
          </GuestRoute>
        ),
      },
      {
        path: "/forgot-password",
        element: (
          <GuestRoute>
            <ForgotPassword />
          </GuestRoute>
        ),
      },
      {
        path: "/reset-password",
        element: (
          <GuestRoute>
            <ResetPassword />
          </GuestRoute>
        ),
      },
      {
        path: "*",
        element: <ErrorPage />,
      },
    ],
  },
]);
```

- [ ] **Step 4: Commit**

```bash
git add frontend/src/router.tsx frontend/src/Routes.tsx frontend/src/contexts/authContext.tsx
git commit -m "feat: update router with guest routes, error handling, and simplified auth"
```

---

### Task 9: Frontend — Update Auth Pages (SignIn, SignUp)

**Files:**
- Modify: `frontend/src/pages/SignIn.tsx`
- Modify: `frontend/src/pages/SignUp.tsx`

- [ ] **Step 1: Update SignIn page with split layout**

Replace entire content of `frontend/src/pages/SignIn.tsx`:

```tsx
import { useState } from "react";
import { Link, useNavigate } from "react-router";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { usePostApiV1Signin } from "@/services/api/v1-public";
import { Loader2 } from "lucide-react";

export function SignIn() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const navigate = useNavigate();
  const signin = usePostApiV1Signin();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    try {
      await signin.mutateAsync({ data: { email, password } });
      navigate("/");
    } catch {
      setError("Invalid email or password");
    }
  };

  return (
    <div className="flex min-h-screen">
      {/* Left — Form */}
      <div className="flex w-full items-center justify-center px-6 py-12 lg:w-1/2 lg:shrink-0">
        <div className="w-full max-w-sm">
          {/* Header */}
          <div className="mb-8">
            <h1 className="text-2xl font-bold tracking-tight text-foreground">
              Welcome back
            </h1>
            <p className="mt-1.5 text-sm text-muted-foreground">
              Sign in to your account to continue
            </p>
          </div>

          {/* Form */}
          <form onSubmit={handleSubmit} className="space-y-5">
            {error && (
              <div className="rounded-lg border border-destructive/20 bg-destructive/5 px-4 py-3 text-sm text-destructive">
                {error}
              </div>
            )}

            <div className="space-y-2">
              <Label htmlFor="email" className="text-sm font-medium">
                Email
              </Label>
              <Input
                id="email"
                type="email"
                placeholder="name@example.com"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
                className="h-11"
              />
            </div>

            <div className="space-y-2">
              <div className="flex items-center justify-between">
                <Label htmlFor="password" className="text-sm font-medium">
                  Password
                </Label>
                <Link
                  to="/forgot-password"
                  className="text-xs text-muted-foreground transition-colors hover:text-primary"
                >
                  Forgot password?
                </Link>
              </div>
              <Input
                id="password"
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
                className="h-11"
              />
            </div>

            <Button
              type="submit"
              className="h-11 w-full text-sm font-medium"
              disabled={signin.isPending}
            >
              {signin.isPending ? (
                <>
                  <Loader2 className="mr-2 size-4 animate-spin" />
                  Signing in...
                </>
              ) : (
                "Sign In"
              )}
            </Button>
          </form>

          {/* Footer */}
          <p className="mt-8 text-center text-sm text-muted-foreground">
            Don't have an account?{" "}
            <Link
              to="/signup"
              className="font-medium text-primary transition-colors hover:text-primary/80"
            >
              Create one
            </Link>
          </p>
        </div>
      </div>

      {/* Right — Gradient Panel */}
      <div className="hidden lg:flex lg:w-1/2 lg:items-center lg:justify-center bg-gradient-to-br from-primary/5 via-primary/10 to-primary/5">
        <div className="h-64 w-64 rounded-full bg-primary/5" />
      </div>
    </div>
  );
}
```

- [ ] **Step 2: Update SignUp page with split layout and password requirements**

Replace entire content of `frontend/src/pages/SignUp.tsx`:

```tsx
import { useState } from "react";
import { Link, useNavigate } from "react-router";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { usePostApiV1Signup } from "@/services/api/v1-public";
import { PasswordRequirements } from "@/components/PasswordRequirements";
import { Loader2 } from "lucide-react";

export function SignUp() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [error, setError] = useState("");
  const navigate = useNavigate();
  const signup = usePostApiV1Signup();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    if (password !== confirmPassword) {
      setError("Passwords don't match");
      return;
    }
    try {
      await signup.mutateAsync({
        data: { email, password, confirm_password: confirmPassword },
      });
      navigate("/");
    } catch {
      setError("Signup failed. Please try again.");
    }
  };

  return (
    <div className="flex min-h-screen">
      {/* Left — Form */}
      <div className="flex w-full items-center justify-center px-6 py-12 lg:w-1/2 lg:shrink-0">
        <div className="w-full max-w-sm">
          {/* Header */}
          <div className="mb-8">
            <h1 className="text-2xl font-bold tracking-tight text-foreground">
              Create your account
            </h1>
            <p className="mt-1.5 text-sm text-muted-foreground">
              Get started — it only takes a minute
            </p>
          </div>

          {/* Form */}
          <form onSubmit={handleSubmit} className="space-y-5">
            {error && (
              <div className="rounded-lg border border-destructive/20 bg-destructive/5 px-4 py-3 text-sm text-destructive">
                {error}
              </div>
            )}

            <div className="space-y-2">
              <Label htmlFor="email" className="text-sm font-medium">
                Email
              </Label>
              <Input
                id="email"
                type="email"
                placeholder="name@example.com"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
                className="h-11"
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="password" className="text-sm font-medium">
                Password
              </Label>
              <Input
                id="password"
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
                minLength={8}
                className="h-11"
              />
              <PasswordRequirements password={password} />
            </div>

            <div className="space-y-2">
              <Label htmlFor="confirm-password" className="text-sm font-medium">
                Confirm Password
              </Label>
              <Input
                id="confirm-password"
                type="password"
                value={confirmPassword}
                onChange={(e) => setConfirmPassword(e.target.value)}
                required
                className="h-11"
              />
            </div>

            <Button
              type="submit"
              className="h-11 w-full text-sm font-medium"
              disabled={signup.isPending}
            >
              {signup.isPending ? (
                <>
                  <Loader2 className="mr-2 size-4 animate-spin" />
                  Creating account...
                </>
              ) : (
                "Sign Up"
              )}
            </Button>
          </form>

          {/* Footer */}
          <p className="mt-8 text-center text-sm text-muted-foreground">
            Already have an account?{" "}
            <Link
              to="/signin"
              className="font-medium text-primary transition-colors hover:text-primary/80"
            >
              Sign in
            </Link>
          </p>
        </div>
      </div>

      {/* Right — Gradient Panel */}
      <div className="hidden lg:flex lg:w-1/2 lg:items-center lg:justify-center bg-gradient-to-br from-primary/5 via-primary/10 to-primary/5">
        <div className="h-64 w-64 rounded-full bg-primary/5" />
      </div>
    </div>
  );
}
```

- [ ] **Step 3: Commit**

```bash
git add frontend/src/pages/SignIn.tsx frontend/src/pages/SignUp.tsx
git commit -m "feat: update auth pages with split layout and improved UX"
```

---

### Task 10: Frontend — Update Sidebar, Layout, and Theming

**Files:**
- Modify: `frontend/src/components/app-sidebar.tsx`
- Modify: `frontend/src/components/AppLayout.tsx`
- Modify: `frontend/src/index.css`
- Modify: `frontend/index.html`

- [ ] **Step 1: Update app-sidebar with user context and dropdown**

Replace entire content of `frontend/src/components/app-sidebar.tsx`:

```tsx
import { useContext } from "react";
import { Home, Settings, LogOut, ChevronsUpDown } from "lucide-react";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { useNavigate, useLocation } from "react-router";

import {
  Sidebar,
  SidebarContent,
  SidebarGroup,
  SidebarGroupContent,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarHeader,
  SidebarFooter,
  SidebarRail,
  SidebarSeparator,
} from "@/components/ui/sidebar";
import { usePostApiV1Signout } from "@/services/api/v1-public";
import { AuthContext } from "@/contexts/authContext";
import { queryClient } from "@/providers/QueryProvider";

export function AppSidebar() {
  const navigate = useNavigate();
  const location = useLocation();
  const signout = usePostApiV1Signout();
  const { user } = useContext(AuthContext);

  const userInitials = user
    ? [user.firstName, user.lastName]
        .filter(Boolean)
        .map((n) => n![0].toUpperCase())
        .join("") || user.email[0].toUpperCase()
    : "?";

  const userDisplayName = user
    ? [user.firstName, user.lastName].filter(Boolean).join(" ") || user.email
    : "";

  const handleSignout = async () => {
    try {
      await signout.mutateAsync();
    } catch {
      // continue even if request fails
    }
    queryClient.clear();
    navigate("/signin");
  };

  return (
    <Sidebar collapsible="icon">
      <SidebarHeader className="pb-0">
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton size="lg" asChild>
              <a href="/" className="group/brand">
                <div className="flex aspect-square size-9 shrink-0 items-center justify-center rounded-xl bg-sidebar-primary text-sidebar-primary-foreground transition-transform duration-200 group-hover/brand:scale-105">
                  <Home className="size-4" />
                </div>
                <div className="flex flex-col gap-0.5 leading-none">
                  <span className="text-sm font-semibold tracking-tight">
                    SaaS Template
                  </span>
                  <span className="text-[10px] text-sidebar-foreground/50">
                    v1.0.0
                  </span>
                </div>
              </a>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarHeader>

      <SidebarSeparator className="my-3" />

      <SidebarContent>
        <SidebarGroup className="pt-0">
          <SidebarGroupContent>
            <SidebarMenu>
              <SidebarMenuItem>
                <SidebarMenuButton asChild isActive={location.pathname === "/"}>
                  <a href="/">
                    <Home className="size-4" />
                    <span>Home</span>
                  </a>
                </SidebarMenuButton>
              </SidebarMenuItem>
              <SidebarMenuItem>
                <SidebarMenuButton
                  asChild
                  isActive={location.pathname === "/settings"}
                >
                  <a href="/settings">
                    <Settings className="size-4" />
                    <span>Settings</span>
                  </a>
                </SidebarMenuButton>
              </SidebarMenuItem>
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>

      <SidebarFooter>
        <SidebarSeparator className="mb-2" />
        <SidebarMenu>
          <SidebarMenuItem>
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <SidebarMenuButton
                  size="lg"
                  className="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
                >
                  <Avatar className="size-8 rounded-lg">
                    <AvatarFallback className="rounded-lg bg-sidebar-accent text-sidebar-accent-foreground text-xs font-medium">
                      {userInitials}
                    </AvatarFallback>
                  </Avatar>
                  <div className="grid flex-1 text-left text-sm leading-tight">
                    <span className="truncate font-semibold">
                      {userDisplayName}
                    </span>
                    <span className="truncate text-xs text-sidebar-foreground/50">
                      {user?.email}
                    </span>
                  </div>
                  <ChevronsUpDown className="ml-auto size-4" />
                </SidebarMenuButton>
              </DropdownMenuTrigger>
              <DropdownMenuContent
                className="w-(--radix-dropdown-menu-trigger-width) min-w-56 rounded-lg"
                side="top"
                align="end"
                sideOffset={4}
              >
                <DropdownMenuItem onClick={handleSignout}>
                  <LogOut className="mr-2 size-4" />
                  Sign out
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  );
}
```

- [ ] **Step 2: Update AppLayout with minimal header**

Replace entire content of `frontend/src/components/AppLayout.tsx`:

```tsx
import { Outlet } from "react-router";
import { AppSidebar } from "./app-sidebar";
import {
  SidebarInset,
  SidebarProvider,
  SidebarTrigger,
} from "@/components/ui/sidebar";

export function AppLayout() {
  return (
    <SidebarProvider defaultOpen={false}>
      <AppSidebar />
      <SidebarInset>
        <header className="flex h-12 shrink-0 items-center gap-2 border-b px-4">
          <SidebarTrigger className="-ml-1" />
        </header>
        <main className="flex flex-1 flex-col gap-4 p-4">
          <Outlet />
        </main>
      </SidebarInset>
    </SidebarProvider>
  );
}
```

- [ ] **Step 3: Update index.css with enhanced theming**

Replace entire content of `frontend/src/index.css`:

```css
@import url('https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&display=swap');
@import "tailwindcss";
@import "tw-animate-css";

@custom-variant dark (&:is(.dark *));

@theme inline {
  --radius-sm: calc(var(--radius) - 4px);
  --radius-md: calc(var(--radius) - 2px);
  --radius-lg: var(--radius);
  --radius-xl: calc(var(--radius) + 4px);
  --radius-2xl: calc(var(--radius) + 8px);
  --radius-3xl: calc(var(--radius) + 12px);
  --radius-4xl: calc(var(--radius) + 16px);
  --color-background: var(--background);
  --color-foreground: var(--foreground);
  --color-card: var(--card);
  --color-card-foreground: var(--card-foreground);
  --color-popover: var(--popover);
  --color-popover-foreground: var(--popover-foreground);
  --color-primary: var(--primary);
  --color-primary-foreground: var(--primary-foreground);
  --color-secondary: var(--secondary);
  --color-secondary-foreground: var(--secondary-foreground);
  --color-muted: var(--muted);
  --color-muted-foreground: var(--muted-foreground);
  --color-accent: var(--accent);
  --color-accent-foreground: var(--accent-foreground);
  --color-destructive: var(--destructive);
  --color-border: var(--border);
  --color-input: var(--input);
  --color-ring: var(--ring);
  --color-chart-1: var(--chart-1);
  --color-chart-2: var(--chart-2);
  --color-chart-3: var(--chart-3);
  --color-chart-4: var(--chart-4);
  --color-chart-5: var(--chart-5);
  --color-sidebar: var(--sidebar);
  --color-sidebar-foreground: var(--sidebar-foreground);
  --color-sidebar-primary: var(--sidebar-primary);
  --color-sidebar-primary-foreground: var(--sidebar-primary-foreground);
  --color-sidebar-accent: var(--sidebar-accent);
  --color-sidebar-accent-foreground: var(--sidebar-accent-foreground);
  --color-sidebar-border: var(--sidebar-border);
  --color-sidebar-ring: var(--sidebar-ring);
}

:root {
  --radius: 0.625rem;
  --background: oklch(1 0 0);
  --foreground: oklch(0.145 0 0);
  --card: oklch(1 0 0);
  --card-foreground: oklch(0.145 0 0);
  --popover: oklch(1 0 0);
  --popover-foreground: oklch(0.145 0 0);
  --primary: oklch(0.205 0 0);
  --primary-foreground: oklch(0.985 0 0);
  --secondary: oklch(0.97 0 0);
  --secondary-foreground: oklch(0.205 0 0);
  --muted: oklch(0.97 0 0);
  --muted-foreground: oklch(0.556 0 0);
  --accent: oklch(0.97 0 0);
  --accent-foreground: oklch(0.205 0 0);
  --destructive: oklch(0.577 0.245 27.325);
  --border: oklch(0.922 0 0);
  --input: oklch(0.922 0 0);
  --ring: oklch(0.708 0 0);
  --chart-1: oklch(0.646 0.222 41.116);
  --chart-2: oklch(0.6 0.118 184.704);
  --chart-3: oklch(0.398 0.07 227.392);
  --chart-4: oklch(0.828 0.189 84.429);
  --chart-5: oklch(0.769 0.188 70.08);
  --sidebar: oklch(0.985 0 0);
  --sidebar-foreground: oklch(0.145 0 0);
  --sidebar-primary: oklch(0.205 0 0);
  --sidebar-primary-foreground: oklch(0.985 0 0);
  --sidebar-accent: oklch(0.97 0 0);
  --sidebar-accent-foreground: oklch(0.205 0 0);
  --sidebar-border: oklch(0.922 0 0);
  --sidebar-ring: oklch(0.708 0 0);
}

.dark {
  --background: oklch(0.145 0 0);
  --foreground: oklch(0.985 0 0);
  --card: oklch(0.205 0 0);
  --card-foreground: oklch(0.985 0 0);
  --popover: oklch(0.205 0 0);
  --popover-foreground: oklch(0.985 0 0);
  --primary: oklch(0.922 0 0);
  --primary-foreground: oklch(0.205 0 0);
  --secondary: oklch(0.269 0 0);
  --secondary-foreground: oklch(0.985 0 0);
  --muted: oklch(0.269 0 0);
  --muted-foreground: oklch(0.708 0 0);
  --accent: oklch(0.269 0 0);
  --accent-foreground: oklch(0.985 0 0);
  --destructive: oklch(0.704 0.191 22.216);
  --border: oklch(1 0 0 / 10%);
  --input: oklch(1 0 0 / 15%);
  --ring: oklch(0.556 0 0);
  --chart-1: oklch(0.488 0.243 264.376);
  --chart-2: oklch(0.696 0.17 162.48);
  --chart-3: oklch(0.769 0.188 70.08);
  --chart-4: oklch(0.627 0.265 303.9);
  --chart-5: oklch(0.645 0.246 16.439);
  --sidebar: oklch(0.205 0 0);
  --sidebar-foreground: oklch(0.985 0 0);
  --sidebar-primary: oklch(0.488 0.243 264.376);
  --sidebar-primary-foreground: oklch(0.985 0 0);
  --sidebar-accent: oklch(0.269 0 0);
  --sidebar-accent-foreground: oklch(0.985 0 0);
  --sidebar-border: oklch(1 0 0 / 10%);
  --sidebar-ring: oklch(0.556 0 0);
}

@layer base {
  * {
    @apply border-border outline-ring/50;
  }
  body {
    @apply bg-background text-foreground;
    font-family: 'Inter', sans-serif;
  }
  h1, h2, h3, h4, h5, h6 {
    letter-spacing: -0.3px;
  }
}

@layer components {
  @keyframes fade-in-up {
    from {
      opacity: 0;
      transform: translateY(20px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  .animate-fade-in-up {
    animation: fade-in-up 0.6s ease-out forwards;
    opacity: 0;
  }

  .animation-delay-100 { animation-delay: 100ms; }
  .animation-delay-200 { animation-delay: 200ms; }
  .animation-delay-300 { animation-delay: 300ms; }
  .animation-delay-400 { animation-delay: 400ms; }
}
```

Note: We keep the existing neutral/grayscale color scheme (appropriate for a template) but add Inter font, heading letter-spacing, and fade-in animations. We deliberately don't port the batch purple theme — that's product-specific.

- [ ] **Step 4: Update index.html title**

In `frontend/index.html`, change `<title>frontend</title>` to `<title>SaaS Template</title>`.

- [ ] **Step 5: Verify frontend builds**

```bash
cd /Users/greg/dev/projects/saas-template/frontend
bun run build
```

Expected: Successful build.

- [ ] **Step 6: Commit**

```bash
git add frontend/src/components/app-sidebar.tsx frontend/src/components/AppLayout.tsx frontend/src/index.css frontend/index.html
git commit -m "feat: update sidebar with user context, minimal layout, and enhanced theming"
```

---

### Task 11: Infrastructure — Add Production Docker Compose and Nginx

**Files:**
- Create: `docker-compose.prod.yml`
- Create: `nginx/nginx.conf`
- Create: `frontend/Dockerfile`
- Create: `frontend/nginx.conf`

- [ ] **Step 1: Create frontend Dockerfile**

Create `frontend/Dockerfile`:

```dockerfile
FROM node:22-alpine AS builder

WORKDIR /app

RUN npm install -g bun

COPY package.json bun.lock ./
RUN bun install --frozen-lockfile

COPY . .
RUN bun run build

# Serve with nginx
FROM nginx:alpine

COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
```

- [ ] **Step 2: Create frontend nginx.conf (SPA fallback)**

Create `frontend/nginx.conf`:

```nginx
server {
    listen 80;
    root /usr/share/nginx/html;
    index index.html;

    # SPA fallback - serve index.html for all routes
    location / {
        try_files $uri $uri/ /index.html;
    }

    # Cache static assets
    location /assets/ {
        expires 1y;
        add_header Cache-Control "public, immutable";
    }
}
```

- [ ] **Step 3: Create reverse proxy nginx config**

Create `nginx/nginx.conf`:

```nginx
upstream backend {
    server backend:8008;
}

upstream frontend {
    server frontend:80;
}

server {
    listen 80;
    server_name localhost;

    client_max_body_size 20m;

    # Security headers
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-Frame-Options "DENY" always;
    add_header Referrer-Policy "strict-origin-when-cross-origin" always;
    add_header Permissions-Policy "camera=(), microphone=(), geolocation=()" always;

    location /api/ {
        proxy_pass http://backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto https;
        proxy_connect_timeout 10s;
        proxy_read_timeout 60s;
    }

    location / {
        proxy_pass http://frontend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto https;
    }
}
```

- [ ] **Step 4: Create docker-compose.prod.yml**

Create `docker-compose.prod.yml`:

```yaml
services:
  backend:
    image: ghcr.io/your-org/saas-backend:latest
    container_name: saas-backend
    restart: unless-stopped
    networks:
      - app
    environment:
      APP_ENV: production
      SERVER_PORT: "8008"
      LOG_LEVEL: info
      APP_BASE_URL: ${APP_BASE_URL}
      APP_SECRET: ${APP_SECRET}
      SESSION_SECRET: ${SESSION_SECRET}
      DATABASE_URL: ${DATABASE_URL}
      SUPABASE_URL: ${SUPABASE_URL:-}
      SUPABASE_KEY: ${SUPABASE_KEY:-}
      GOOGLE_CLIENT_ID: ${GOOGLE_CLIENT_ID:-}
      GOOGLE_CLIENT_SECRET: ${GOOGLE_CLIENT_SECRET:-}
      GITHUB_CLIENT_ID: ${GITHUB_CLIENT_ID:-}
      GITHUB_CLIENT_SECRET: ${GITHUB_CLIENT_SECRET:-}
      RESEND_API_KEY: ${RESEND_API_KEY:-}
      RESEND_FROM_EMAIL: ${RESEND_FROM_EMAIL:-noreply@example.com}
    entrypoint:
      - /bin/sh
      - -c
      - |
        set -e
        echo "Running migrations..."
        /goose -dir /migrations postgres "$DATABASE_URL" up
        echo "Starting server..."
        exec /server
    healthcheck:
      test: ["CMD", "wget", "-q", "--spider", "http://localhost:8008/api/public/v1/health"]
      interval: 30s
      timeout: 5s
      retries: 3
      start_period: 10s
    logging:
      driver: json-file
      options:
        max-size: "10m"
        max-file: "3"

  frontend:
    image: ghcr.io/your-org/saas-frontend:latest
    container_name: saas-frontend
    restart: unless-stopped
    networks:
      - app
    logging:
      driver: json-file
      options:
        max-size: "10m"
        max-file: "3"

  nginx:
    image: nginx:alpine
    container_name: saas-nginx
    restart: unless-stopped
    networks:
      - app
    ports:
      - "80:80"
    volumes:
      - ./nginx/nginx.conf:/etc/nginx/conf.d/default.conf:ro
    depends_on:
      backend:
        condition: service_healthy
      frontend:
        condition: service_started
    logging:
      driver: json-file
      options:
        max-size: "10m"
        max-file: "3"

  dozzle:
    image: amir20/dozzle:latest
    container_name: saas-dozzle
    restart: unless-stopped
    ports:
      - "127.0.0.1:8888:8080"
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro
    logging:
      driver: json-file
      options:
        max-size: "10m"
        max-file: "3"

networks:
  app:
    enable_ipv6: true
    ipam:
      config:
        - subnet: 172.28.0.0/16
        - subnet: fd00:db8:1::/64
```

- [ ] **Step 5: Commit**

```bash
git add frontend/Dockerfile frontend/nginx.conf nginx/nginx.conf docker-compose.prod.yml
git commit -m "feat: add production Docker Compose with Nginx reverse proxy"
```

---

### Task 12: Final Verification

- [ ] **Step 1: Verify backend compiles**

```bash
cd /Users/greg/dev/projects/saas-template/backend
go build ./...
```

Expected: No errors.

- [ ] **Step 2: Verify frontend builds**

```bash
cd /Users/greg/dev/projects/saas-template/frontend
bun run build
```

Expected: No errors.

- [ ] **Step 3: Verify no sensitive files staged**

```bash
cd /Users/greg/dev/projects/saas-template
git status
```

Check that no `.env` files or secrets are staged.

- [ ] **Step 4: Review all commits**

```bash
git log --oneline -10
```

Expected: Clean commit history with descriptive conventional commit messages.
