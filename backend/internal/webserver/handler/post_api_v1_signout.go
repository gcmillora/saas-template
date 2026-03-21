package handler

import (
	oapi_public "saas-template/generated/oapi/public"
	"saas-template/internal/webserver/middleware"
	"context"
	"errors"
)

func (h *Handler) PostApiV1Signout(ctx context.Context, request oapi_public.PostApiV1SignoutRequestObject) (oapi_public.PostApiV1SignoutResponseObject, error) {
	w := middleware.GetResponseWriter(ctx)
	r := middleware.GetRequest(ctx)
	if w == nil || r == nil {
		return nil, errors.New("could not get http primitives from context")
	}

	session, err := h.app.Session().Get(r, "session")
	if err != nil {
		return nil, err
	}

	session.Options.MaxAge = -1
	if err := session.Save(r, w); err != nil {
		return nil, err
	}

	return oapi_public.PostApiV1Signout200TextResponse("Signout successful."), nil
}
