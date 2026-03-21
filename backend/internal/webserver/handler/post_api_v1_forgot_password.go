package handler

import (
	"context"
	oapi_public "saas-template/generated/oapi/public"
	"saas-template/internal/app/app_service/authentication"
)

// POST (/forgot-password)
func (h *Handler) PostApiV1ForgotPassword(
	ctx context.Context,
	request oapi_public.PostApiV1ForgotPasswordRequestObject,
) (oapi_public.PostApiV1ForgotPasswordResponseObject, error) {
	authentication.ForgotPassword(ctx, h.app, string(request.Body.Email))

	msg := "If an account with that email exists, we've sent a reset link."
	return oapi_public.PostApiV1ForgotPassword200JSONResponse{Message: &msg}, nil
}
