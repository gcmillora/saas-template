package middleware

import (
	"net/http"
	"saas-template/config"
)

func NewAdminMiddleware(app *config.App) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, _ := app.Session().Get(r, "session")

			if role, ok := session.Values["role"].(string); ok && role == "admin" {
				next.ServeHTTP(w, r)
			} else {
				http.Error(w, "Forbidden", http.StatusForbidden)
			}
		})
	}
}
