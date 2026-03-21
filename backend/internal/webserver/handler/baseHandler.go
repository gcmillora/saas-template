// Package handler implements HTTP request handlers for the API.
//
// Handlers follow the strict server interface pattern from oapi-codegen,
// which generates type-safe interfaces from OpenAPI specifications.
//
// Handler Structure:
//   - All handlers are methods on the Handler struct
//   - Method names match the operationId from openapi.yaml
//   - Request/Response types are generated from OpenAPI schemas
//
// The Handler struct provides access to:
//   - h.app: Configuration and dependencies (DB, Session, Logger, etc.)
//
// See handler_template.go.example for implementation patterns and best practices.
package handler

import (
	"saas-template/config"
	"saas-template/generated/oapi"
	oapi_public "saas-template/generated/oapi/public"
)

// CombinedStrictServerInterface combines both the main API and public API interfaces.
// This allows a single Handler struct to implement endpoints from multiple OpenAPI specs.
type CombinedStrictServerInterface interface {
	oapi.StrictServerInterface        // Main API (openapi.yaml)
	oapi_public.StrictServerInterface // Public API (openapi-public.yaml)
}

// Handler implements all API endpoints defined in the OpenAPI specifications.
// Each handler method:
//   - Receives a context and a typed request object
//   - Returns a typed response object or an error
//   - Has access to app configuration via h.app
//
// Example handler implementations can be found in:
//   - get_api_v1_me_user.go (authenticated GET endpoint)
//   - post_api_v1_signin.go (session management)
//   - post_api_v1_signup.go (app_service orchestration)
type Handler struct {
	app *config.App
	CombinedStrictServerInterface
}

// NewHandler creates a new Handler instance with access to the application configuration.
// The Handler implements all endpoints defined in the OpenAPI specifications.
func NewHandler(app *config.App) *Handler {
	handler := Handler{
		app: app,
	}

	return &handler
}
