package handler

import (
	"net/http"

	"github.com/Kanbenn/gophermart/internal/jwtoken"
)

func writeAuthCookie(w http.ResponseWriter, uid int) error {
	tokenString, err := jwtoken.MakeToken(uid)
	if err != nil {
		return err
	}
	newCookie := &http.Cookie{
		Name:  "auth",
		Value: tokenString,
	}
	http.SetCookie(w, newCookie)
	return nil
}
