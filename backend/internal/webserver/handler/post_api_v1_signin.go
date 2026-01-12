package handler

import (
	oapi_public "adobo/generated/oapi/public"
	"adobo/internal/app/repository"
	"adobo/internal/webserver/middleware"
	"context"
	"errors"
	"log/slog"
)

// PostApiV1Signin handles POST /signin endpoint.
//
// This is an example of a session management endpoint that:
//   - Authenticates user with Supabase
//   - Retrieves user data from local database
//   - Creates a session cookie with user_id and tenant_id
//   - Returns success message
//
// Pattern: Direct HTTP primitive access for session/cookie operations.
// Use middleware.GetResponseWriter() and middleware.GetRequest() to access
// http.ResponseWriter and http.Request when needed for cookies/headers.
func (h *Handler) PostApiV1Signin(ctx context.Context, request oapi_public.PostApiV1SigninRequestObject) (oapi_public.PostApiV1SigninResponseObject, error) {
	// 1. Get HTTP primitives for session management
	w := middleware.GetResponseWriter(ctx)
	r := middleware.GetRequest(ctx)
	if w == nil || r == nil {
		return nil, errors.New("could not get http.ResponseWriter or http.Request from context")
	}

	// 2. Authenticate with Supabase
	res, err := h.app.Supabase().Client().Auth.SignInWithEmailPassword(
		*request.Body.Email,
		*request.Body.Password,
	)
	if err != nil {
		return nil, err
	}

	// 3. Get or create session
	session, err := h.app.Session().Get(r, "session")
	if err != nil {
		return nil, err
	}

	// 4. Fetch user from local database using Supabase auth ID
	user, err := repository.GetUserByAuthID(ctx, h.app.DB(), res.User.ID)
	if err != nil {
		return nil, err
	}

	// 5. Store user_id and tenant_id in session
	session.Values["user_id"] = res.User.ID.String()
	session.Values["tenant_id"] = user.TenantID.String()
	if err := session.Save(r, w); err != nil {
		return nil, err
	}

	// 6. Log successful login
	slog.Default().Info("Login successful", "auth_id", res.User.ID)

	return oapi_public.PostApiV1Signin200TextResponse("Login successful. Session cookie is set."), nil
}
