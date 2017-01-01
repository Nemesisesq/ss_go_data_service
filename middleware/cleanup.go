package middleware

import (
	"context"
	"github.com/codegangsta/negroni"
	"net/http"
)

func CleanupMiddleware() negroni.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		cleanup_chan := make(chan string)
		ctx := context.WithValue(r.Context(), "cleanup", cleanup_chan)
		next.ServeHTTP(w, r.WithContext(ctx))
		cleanup_chan <- "cleanup"
	}
}
