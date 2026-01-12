package middleware

import (
	"context"
	"net/http"
)

type contextKey string

const (
	ResponseWriterKey contextKey = "responseWriter"
	RequestKey        contextKey = "request"
)

func NewContextInjectorMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), ResponseWriterKey, w)
			ctx = context.WithValue(ctx, RequestKey, r)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetResponseWriter(ctx context.Context) http.ResponseWriter {
	if rw, ok := ctx.Value(ResponseWriterKey).(http.ResponseWriter); ok {
		return rw
	}
	return nil
}

func GetRequest(ctx context.Context) *http.Request {
	if r, ok := ctx.Value(RequestKey).(*http.Request); ok {
		return r
	}
	return nil
}
