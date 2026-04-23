package middleware

import (
	"context"
	"net/http"
	"saas-template/config"
)

const userIDKey contextKey = "user_id"

func NewAuthMiddleware(app *config.App) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, _ := app.Session().Get(r, "session")

			if userID, ok := session.Values["user_id"].(string); ok && userID != "" {
				ctx := context.WithValue(r.Context(), userIDKey, userID)
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
			}
		})
	}
}
