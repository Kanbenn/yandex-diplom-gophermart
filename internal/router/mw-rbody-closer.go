package router

import "net/http"

func rBodyCloserMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		next.ServeHTTP(w, r)
		r.Body.Close()
	}
	return http.HandlerFunc(fn)
}
