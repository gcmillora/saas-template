# Handler Directory

This directory contains HTTP request handlers that implement the API endpoints defined in the OpenAPI specifications.

## Overview

Handlers follow the **strict server interface pattern** from `oapi-codegen`:
- Type-safe request/response objects generated from OpenAPI schemas
- Method names match `operationId` from `openapi.yaml` and `openapi-public.yaml`
- All handlers are methods on the `Handler` struct

## Files

### Core Files

- **`baseHandler.go`** - Handler struct definition and constructor
  - Implements `CombinedStrictServerInterface` for both main and public APIs
  - Provides access to app configuration via `h.app`

- **`handler_template.go.example`** - Complete implementation guide
  - 8 detailed examples covering all common patterns
  - Best practices and anti-patterns
  - Copy this as a reference when creating new handlers

### Example Handlers

These handlers demonstrate different patterns:

1. **`get_api_v1_me_user.go`** - Authenticated GET endpoint
   - Pattern: Delegate to app_service for orchestration
   - Shows: Session handling, repository calls, response mapping

2. **`post_api_v1_signin.go`** - Session management
   - Pattern: Direct HTTP primitive access for cookies
   - Shows: Supabase auth, session creation, repository lookup

3. **`post_api_v1_signup.go`** - App service delegation
   - Pattern: Complex business logic in app_service
   - Shows: Input mapping, multi-step operations, validation

4. **`get_api_v1_health.go`** - Simple health check
   - Pattern: Minimal logic, quick response
   - Shows: Basic endpoint structure

## Creating New Handlers

### Step 1: Define in OpenAPI

Add your endpoint to `openapi.yaml` or `openapi-public.yaml`:

```yaml
paths:
  /api/v1/items:
    get:
      operationId: get-api-v1-items  # This becomes the method name
      summary: "Get all items"
      responses:
        "200":
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: '#/components/schemas/Item'
```

### Step 2: Generate Code

```bash
./commands.sh openapi:codegen
```

This generates the interface that your handler must implement.

### Step 3: Implement Handler

Create a new file `get_api_v1_items.go`:

```go
package handler

import (
	"adobo/generated/oapi"
	"adobo/internal/app/repository"
	"context"
)

func (h *Handler) GetApiV1Items(
	ctx context.Context,
	request oapi.GetApiV1ItemsRequestObject,
) (oapi.GetApiV1ItemsResponseObject, error) {
	// Your implementation here
	items, err := repository.GetItems(ctx, h.app.DB())
	if err != nil {
		return nil, err
	}
	
	return oapi.GetApiV1Items200JSONResponse(*items), nil
}
```

## Handler Patterns

### Pattern 1: Simple Read (No Auth)

```go
func (h *Handler) GetApiV1PublicData(ctx, req) (resp, error) {
	data, err := repository.GetPublicData(ctx, h.app.DB())
	if err != nil {
		return nil, err
	}
	return oapi.GetApiV1PublicData200JSONResponse(*data), nil
}
```

### Pattern 2: Authenticated Read

```go
func (h *Handler) GetApiV1MyData(ctx, req) (resp, error) {
	// Extract session
	sessionData, err := middleware.GetSessionData(ctx, h.app.Session())
	if err != nil {
		return nil, err
	}
	
	tenantID, _ := uuid.Parse(sessionData.TenantID)
	
	// Tenant-scoped query
	data, err := repository.GetDataByTenant(ctx, h.app.DB(), tenantID)
	if err != nil {
		return nil, err
	}
	
	return oapi.GetApiV1MyData200JSONResponse(*data), nil
}
```

### Pattern 3: Create Operation

```go
func (h *Handler) PostApiV1Items(ctx, req) (resp, error) {
	sessionData, err := middleware.GetSessionData(ctx, h.app.Session())
	if err != nil {
		return nil, err
	}
	
	tenantID, _ := uuid.Parse(sessionData.TenantID)
	
	// Prepare model
	newItem := model.ItemTbl{
		ID:       uuid.New(),
		TenantID: tenantID,
		Name:     req.Body.Name,
	}
	
	// Create via mutation
	created, err := mutation.CreateItem(ctx, h.app.DB(), newItem)
	if err != nil {
		return nil, err
	}
	
	return oapi.PostApiV1Items201JSONResponse(*created), nil
}
```

### Pattern 4: Complex Logic (Use App Service)

```go
func (h *Handler) PostApiV1ComplexOperation(ctx, req) (resp, error) {
	// Delegate complex logic to app_service
	result, err := app_service.PerformComplexOperation(ctx, h.app, req.Body)
	if err != nil {
		return nil, err
	}
	
	return result, nil
}
```

### Pattern 5: Session/Cookie Operations

```go
func (h *Handler) PostApiV1Logout(ctx, req) (resp, error) {
	// Get HTTP primitives
	w := middleware.GetResponseWriter(ctx)
	r := middleware.GetRequest(ctx)
	if w == nil || r == nil {
		return nil, errors.New("missing HTTP primitives")
	}
	
	// Manipulate session
	session, err := h.app.Session().Get(r, "session")
	if err != nil {
		return nil, err
	}
	
	session.Options.MaxAge = -1
	if err := session.Save(r, w); err != nil {
		return nil, err
	}
	
	return oapi.PostApiV1Logout200Response{}, nil
}
```

## Best Practices

### ✅ DO

- Keep handlers thin - delegate to app_services
- Use `middleware.GetSessionData()` for authentication
- Always verify tenant ownership in queries
- Return errors directly (middleware handles formatting)
- Use repositories for reads, mutations for writes
- Log important operations with `h.app.Logger()`
- Follow REST conventions for status codes (200, 201, 204, etc.)

### ❌ DON'T

- Put business logic in handlers
- Bypass tenant isolation
- Return database errors directly (wrap them)
- Mix repository and mutation logic
- Forget input validation
- Expose internal IDs or sensitive data
- Use raw SQL (use repository/mutation layers)

## Testing

When writing tests for handlers:

```go
func TestGetApiV1Items(t *testing.T) {
	// 1. Setup test app
	app := testutil.NewTestApp(t)
	defer app.Cleanup()
	
	// 2. Create handler
	h := NewHandler(app)
	
	// 3. Create request
	req := oapi.GetApiV1ItemsRequestObject{}
	
	// 4. Call handler
	resp, err := h.GetApiV1Items(context.Background(), req)
	
	// 5. Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}
```

## Additional Resources

- See `handler_template.go.example` for detailed examples
- Review existing handlers for real-world patterns
- Check `../middleware/` for authentication helpers
- Read `../../app/` for repository and mutation patterns
