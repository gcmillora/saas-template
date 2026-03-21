package handler

import (
	"context"
	"errors"
	oapi_public "saas-template/generated/oapi/public"
	"saas-template/internal/app/app_service/authentication"
	"saas-template/internal/webserver/middleware"
)

// POST (/signin)
func (h *Handler) PostApiV1Signin(
	ctx context.Context,
	request oapi_public.PostApiV1SigninRequestObject,
) (oapi_public.PostApiV1SigninResponseObject, error) {
	w := middleware.GetResponseWriter(ctx)
	r := middleware.GetRequest(ctx)
	if w == nil || r == nil {
		return nil, errors.New("could not get http primitives from context")
	}

	user, err := authentication.SignIn(
		ctx,
		h.app,
		string(request.Body.Email),
		request.Body.Password,
	)
	if err != nil {
		msg := err.Error()
		return oapi_public.PostApiV1Signin401JSONResponse{Message: &msg}, nil
	}

	session, err := h.app.Session().Get(r, "session")
	if err != nil {
		return nil, err
	}

	session.Values["user_id"] = user.ID.String()
	session.Values["tenant_id"] = user.TenantID.String()
	session.Values["role"] = user.Role
	if err := session.Save(r, w); err != nil {
		return nil, err
	}

	msg := "Login successful."
	return oapi_public.PostApiV1Signin200JSONResponse{Message: &msg}, nil
}
