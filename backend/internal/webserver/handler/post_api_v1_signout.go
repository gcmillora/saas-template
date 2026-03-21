package handler

import (
	"context"
	"errors"
	oapi_public "saas-template/generated/oapi/public"
	"saas-template/internal/app/app_service/authentication"
	"saas-template/internal/webserver/middleware"

	"github.com/google/uuid"
)

// POST (/signout)
func (h *Handler) PostApiV1Signout(
	ctx context.Context,
	request oapi_public.PostApiV1SignoutRequestObject,
) (oapi_public.PostApiV1SignoutResponseObject, error) {
	w := middleware.GetResponseWriter(ctx)
	r := middleware.GetRequest(ctx)
	if w == nil || r == nil {
		return nil, errors.New("could not get http primitives from context")
	}

	session, err := h.app.Session().Get(r, "session")
	if err != nil {
		return nil, err
	}

	if userIDStr, ok := session.Values["user_id"].(string); ok {
		if userID, err := uuid.Parse(userIDStr); err == nil {
			if tenantIDStr, ok := session.Values["tenant_id"].(string); ok {
				if tenantID, err := uuid.Parse(tenantIDStr); err == nil {
					authentication.SignOut(ctx, h.app, userID, tenantID)
				}
			}
		}
	}

	session.Options.MaxAge = -1
	if err := session.Save(r, w); err != nil {
		return nil, err
	}

	msg := "Signout successful."
	return oapi_public.PostApiV1Signout200JSONResponse{Message: &msg}, nil
}
