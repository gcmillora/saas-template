package handler

import (
	"adobo/generated/oapi"
	"adobo/internal/app/app_service/user"
	"context"
)

// GetApiV1User handles GET /user endpoint.
//
// This is an example of an authenticated endpoint that:
//   - Extracts session data (handled in app_service)
//   - Fetches user data from the database
//   - Returns user information
//
// Pattern: Delegate complex logic to app_service layer.
// The app_service handles session extraction, repository calls, and response mapping.
func (h *Handler) GetApiV1User(ctx context.Context, request oapi.GetApiV1UserRequestObject) (oapi.GetApiV1UserResponseObject, error) {
	// Delegate to app_service for business logic orchestration
	data, err := user.GetUser(ctx, h.app)
	if err != nil {
		return nil, err
	}

	return data, nil
}
