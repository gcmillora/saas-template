package middleware

import (
	"log/slog"
	"net/http"
	"saas-template/config"
)

func NewAdminMiddleware(app *config.App) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, err := app.Session().Get(r, "session")
			if err != nil {
				slog.Warn("failed to read session", "error", err)
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			if role, ok := session.Values["role"].(string); ok && role == "admin" {
				next.ServeHTTP(w, r)
			} else {
				http.Error(w, "Forbidden", http.StatusForbidden)
			}
		})
	}
}
