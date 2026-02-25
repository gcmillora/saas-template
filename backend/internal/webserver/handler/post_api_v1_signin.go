package handler

import (
	oapi_public "adobo/generated/oapi/public"
	"adobo/internal/app/repository"
	"adobo/internal/webserver/middleware"
	"context"
	"errors"
	"log/slog"

	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) PostApiV1Signin(ctx context.Context, request oapi_public.PostApiV1SigninRequestObject) (oapi_public.PostApiV1SigninResponseObject, error) {
	w := middleware.GetResponseWriter(ctx)
	r := middleware.GetRequest(ctx)
	if w == nil || r == nil {
		return nil, errors.New("could not get http primitives from context")
	}

	user, err := repository.GetUserByEmail(ctx, h.app.DB(), string(request.Body.Email))
	if err != nil {
		slog.Default().Info("Signin failed: user not found", "email", request.Body.Email)
		msg := "Invalid email or password"
		return oapi_public.PostApiV1Signin401JSONResponse{Message: &msg}, nil
	}

	if user.PasswordHash == nil {
		msg := "Invalid email or password"
		return oapi_public.PostApiV1Signin401JSONResponse{Message: &msg}, nil
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(request.Body.Password)); err != nil {
		msg := "Invalid email or password"
		return oapi_public.PostApiV1Signin401JSONResponse{Message: &msg}, nil
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

	slog.Default().Info("Login successful", "user_id", user.ID)
	return oapi_public.PostApiV1Signin200TextResponse("Login successful."), nil
}
