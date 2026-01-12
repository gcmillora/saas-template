package handler

import (
	oapi_public "adobo/generated/oapi/public"
	"adobo/internal/app/app_service/authentication"
	"context"
)

// PostApiV1Signup handles POST /signup endpoint.
//
// This is an example of a handler that delegates to an app_service for:
//   - Complex business logic (user creation, validation)
//   - Multiple operations (Supabase auth + local DB)
//   - Input validation beyond schema validation
//
// Pattern: Create a typed body struct and pass to app_service.
// The app_service encapsulates all the complexity, keeping the handler thin.
func (h *Handler) PostApiV1Signup(
	ctx context.Context,
	request oapi_public.PostApiV1SignupRequestObject,
) (oapi_public.PostApiV1SignupResponseObject, error) {
	// 1. Map request body to app_service input type
	body := authentication.PostSignupBody{
		Email:           *request.Body.Email,
		Password:        *request.Body.Password,
		ConfirmPassword: *request.Body.ConfirmPassword,
		TenantId:        *request.Body.TenantId,
	}

	// 2. Delegate to app_service for signup orchestration
	// The app_service handles:
	//   - Password validation
	//   - Supabase user creation
	//   - Local database user record creation
	err := authentication.SignUp(ctx, h.app, body)
	if err != nil {
		return nil, err
	}

	// 3. Return success response
	return oapi_public.PostApiV1Signup200TextResponse("Signup successful."), nil
}
