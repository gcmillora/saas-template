package middleware

import (
	"context"
	"errors"

	"github.com/gorilla/sessions"
)

// SessionData represents the decoded session values
type SessionData struct {
	UserID   string
	TenantID string
	Role     string
}

// GetSessionData retrieves all session values and maps them to SessionData
func GetSessionData(ctx context.Context, store sessions.Store) (*SessionData, error) {
	r := GetRequest(ctx)
	if r == nil {
		return nil, errors.New("could not get http.Request from context")
	}

	session, err := store.Get(r, "session")
	if err != nil {
		return nil, err
	}

	data := &SessionData{}

	if userID, ok := session.Values["user_id"].(string); ok {
		data.UserID = userID
	}

	if tenantID, ok := session.Values["tenant_id"].(string); ok {
		data.TenantID = tenantID
	}

	if role, ok := session.Values["role"].(string); ok {
		data.Role = role
	}

	return data, nil
}
