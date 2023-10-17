package router

import (
	"context"
	"log"
	"net/http"

	"github.com/Kanbenn/gophermart/internal/jwtoken"
)

// had to make the static-check happy >_<
type ctxKey string

func RequireAuthMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		uid, err := parseUidFromCookie(r)
		if err != nil || uid < 1 {
			log.Println("RequireAuthMw error: unauthorized access attempt", err)
			http.Error(w, "Auth required", http.StatusUnauthorized)
			return
		}

		log.Println("RequireAuthMw: putting uid to r.Context", uid)
		var ctxKeyUid ctxKey = "uid"
		ctx := context.WithValue(r.Context(), ctxKeyUid, uid)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

func parseUidFromCookie(r *http.Request) (uid int, err error) {
	cookie, err1 := r.Cookie("auth")
	if err1 != nil {
		return 0, err1
	}
	uid, err2 := jwtoken.ParseToken(cookie.Value)
	if err2 != nil {
		return 0, err2
	}
	return uid, nil
}
