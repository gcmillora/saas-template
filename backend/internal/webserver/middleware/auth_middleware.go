package middleware

import (
	"saas-template/config"
	"context"
	"net/http"
)

func NewAuthMiddleware(app *config.App) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, _ := app.Session().Get(r, "session")

			if userID, ok := session.Values["user_id"].(string); ok && userID != "" {
				ctx := context.WithValue(r.Context(), "user_id", userID)
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
			}
		})
	}
}
