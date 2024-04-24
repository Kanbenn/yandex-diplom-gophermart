package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Kanbenn/gophermart/internal/app"
	"github.com/Kanbenn/gophermart/internal/config"
	"github.com/Kanbenn/gophermart/internal/jwtoken"
)

type Handler struct {
	cfg config.Config
	app *app.App
}

func New(cfg config.Config, app *app.App) *Handler {
	return &Handler{cfg, app}
}

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

func writeJsnResponse(w http.ResponseWriter, v any) {
	out, err := json.Marshal(v)
	if err != nil {
		log.Println("handler error at writing json:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}
