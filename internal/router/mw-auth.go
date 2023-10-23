package router

import (
	"context"
	"log"
	"net/http"

	"github.com/Kanbenn/gophermart/internal/jwtoken"
	"github.com/Kanbenn/gophermart/internal/models"
)

func requireAuthMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		uid, err := parseTokenFromCookie(r)
		if err != nil {
			log.Println("RequireAuthMw error: unauthorized access attempt", err)
			http.Error(w, "Auth required", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), models.CtxKeyUser, uid)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

func parseTokenFromCookie(r *http.Request) (uid int, err error) {
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
