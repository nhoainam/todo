package http_middleware

import (
	"context"
	"net/http"
)

type contextKey string

const ResponseWriterKey contextKey = "http_response_writer"

func InjectHTTPMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), ResponseWriterKey, w)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
