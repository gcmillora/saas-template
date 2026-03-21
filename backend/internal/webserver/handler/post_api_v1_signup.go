package handler

import (
	"context"
	"errors"
	"log/slog"
	oapi_public "saas-template/generated/oapi/public"
	"saas-template/internal/app/app_service/authentication"
	"saas-template/internal/webserver/middleware"
)

// POST (/signup)
func (h *Handler) PostApiV1Signup(
	ctx context.Context,
	request oapi_public.PostApiV1SignupRequestObject,
) (oapi_public.PostApiV1SignupResponseObject, error) {
	w := middleware.GetResponseWriter(ctx)
	r := middleware.GetRequest(ctx)
	if w == nil || r == nil {
		return nil, errors.New("could not get http primitives from context")
	}

	body := authentication.PostSignupBody{
		Email:           string(request.Body.Email),
		Password:        request.Body.Password,
		ConfirmPassword: request.Body.ConfirmPassword,
		FirstName:       request.Body.FirstName,
		LastName:        request.Body.LastName,
	}

	user, err := authentication.SignUp(ctx, h.app, body)
	if err != nil {
		msg := err.Error()
		return oapi_public.PostApiV1Signup400JSONResponse{Message: &msg}, nil
	}

	session, err := h.app.Session().Get(r, "session")
	if err != nil {
		return nil, err
	}

	session.Values["user_id"] = user.ID.String()
	session.Values["tenant_id"] = user.TenantID.String()
	if err := session.Save(r, w); err != nil {
		return nil, err
	}

	slog.Default().Info("Signup successful", "user_id", user.ID)
	msg := "Signup successful."
	return oapi_public.PostApiV1Signup200JSONResponse{Message: &msg}, nil
}
