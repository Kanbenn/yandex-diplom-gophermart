package router

import (
	"log"
	"net/http"
	"strings"
)

func RequireJsnMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		if !isJsnContentType(r) {
			log.Println("RequireJsnMw error: non-json contnent-type", r.Header.Get("Content-Type"))
			http.Error(w, "use the application/json Contnent-Type", http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func isJsnContentType(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Content-Type"), "application/json")
}
