package handler

import (
	oapi_public "adobo/generated/oapi/public"
	"context"
)

// GetApiV1Health handles GET /health endpoint.
//
// This is an example of a simple health check endpoint.
// Typically used by load balancers and monitoring systems.
//
// Pattern: Minimal logic, quick response.
// Optionally check critical dependencies (DB, external APIs).
func (h *Handler) GetApiV1Health(ctx context.Context, request oapi_public.GetApiV1HealthRequestObject) (oapi_public.GetApiV1HealthResponseObject, error) {
	// Simple health check - return 200 OK
	// For production, consider checking:
	//   - Database connectivity: h.app.DB().Ping()
	//   - Cache availability: h.app.Cache()

	return oapi_public.GetApiV1Health200TextResponse("OK"), nil
}
