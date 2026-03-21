package handler

import (
	"context"
	oapi_public "saas-template/generated/oapi/public"
	"saas-template/internal/app/app_service/authentication"
)

// POST (/reset-password)
func (h *Handler) PostApiV1ResetPassword(
	ctx context.Context,
	request oapi_public.PostApiV1ResetPasswordRequestObject,
) (oapi_public.PostApiV1ResetPasswordResponseObject, error) {
	result, err := authentication.ResetPassword(
		ctx,
		h.app,
		request.Body.Token,
		request.Body.Password,
		request.Body.ConfirmPassword,
	)

	if err != nil {
		msg := err.Error()
		return oapi_public.PostApiV1ResetPassword400JSONResponse{
			Message: &msg,
		}, nil
	}

	if result != nil && len(result.ValidationErrors) > 0 {
		msg := "Password does not meet requirements"
		errs := result.ValidationErrors
		return oapi_public.PostApiV1ResetPassword400JSONResponse{
			Message: &msg,
			Data: &struct {
				Errors *[]string `json:"errors,omitempty"`
			}{Errors: &errs},
		}, nil
	}

	msg := "Password has been reset successfully."
	return oapi_public.PostApiV1ResetPassword200JSONResponse{Message: &msg}, nil
}
