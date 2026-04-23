# Repository Guidelines

## Project Structure & Module Organization

Backend code is in `backend/`. API specs live in `backend/openapi*.yaml`; generated Go code is in `backend/generated` and must not be edited. Runtime code is in `backend/internal`: handlers in `webserver/handler`, middleware in `webserver/middleware`, business logic in `app/app_service`, reads in `app/repository`, and writes in `app/mutation`. Frontend code is in `frontend/src`: pages in `pages`, shared components in `components`, shadcn primitives in `components/ui`, hooks in `hooks`, contexts in `contexts`, and Orval hooks in `services/api`.

## Backend Implementation Rules

For new endpoints, update the correct OpenAPI file first, then run `cd backend && bash commands.sh openapi:codegen`. Keep handlers thin: parse request/session data, call one app service function, and map expected errors to typed responses. Do not put business logic or database calls in handlers.

App services own orchestration, validation, authorization-sensitive logic, and audit logging. Use `middleware.GetSessionData(ctx, app.Session())` for session data. Repositories are read-only and named `Get*`. Mutations own inserts, updates, and deletes and are named `Create*`, `Update*`, or `Delete*`; maintain timestamps where applicable.

## Frontend Implementation Rules

Pages belong in `frontend/src/pages` and should compose reusable components. Shared app components belong in `frontend/src/components`; generic primitives belong in `frontend/src/components/ui`. Prefer existing shadcn/ui components before custom markup. Add shadcn components from `frontend/` with `bunx --bun shadcn@latest add <component>`.

Use generated API hooks from `frontend/src/services/api`; do not edit generated API files. Toasts use shadcn's Sonner wrapper at `components/ui/sonner.tsx`; mount `Toaster` once in `App.tsx` and call `toast()` from `sonner`.

## Build, Test, and Development Commands

```bash
./setup.sh
cd backend && bash commands.sh webserver
cd backend && bash commands.sh test
cd backend && bash commands.sh lint
cd backend && bash commands.sh migration:up
cd frontend && bun run dev
cd frontend && bun run lint
cd frontend && bun run build
```

## Coding Style & Naming Conventions

Use `gofmt` for Go. OpenAPI `operationId` controls names: `post-api-v1-signin` becomes `PostApiV1Signin` and `post_api_v1_signin.go`. Keep one app service file per operation. Frontend uses TypeScript, React 19, Tailwind v4, and `@/` imports. Keep code generic; do not introduce product-specific branding, routes, roles, or currency defaults.

## Testing & Review Expectations

Run `cd backend && go test ./...`, `cd frontend && bun run lint`, and `cd frontend && bun run build` before handoff. Add focused backend tests for app service, repository, mutation, and security-flow changes. PRs should explain behavior changes, list verification commands, link issues, and include screenshots for UI changes.
