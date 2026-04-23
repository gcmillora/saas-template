package middleware

import (
	"log/slog"
	"net/http"
	"saas-template/config"

	"github.com/go-chi/chi/v5/middleware"
)

func NewLoggerMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			defer func() {
				slog.Default().InfoContext(r.Context(), "Served",
					slog.String("method", r.Method),
					slog.String("proto", r.Proto),
					slog.String("path", r.URL.Path),
					slog.Int("status", ww.Status()),
					slog.Int("size", ww.BytesWritten()),
				)
			}()

			next.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(fn)
	}
}

func HandleErrorWithLog(app *config.App) func(w http.ResponseWriter, r *http.Request, err error) {
	fn := func(w http.ResponseWriter, r *http.Request, err error) {
		ctx := r.Context()
		slog.Default().ErrorContext(ctx, err.Error())

		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	return fn
}
