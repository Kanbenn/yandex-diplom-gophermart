package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/Kanbenn/gophermart/internal/jwtoken"
	"github.com/Kanbenn/gophermart/internal/models"
	"github.com/Kanbenn/gophermart/internal/storage"
)

func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {

	var u models.UserInsert
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "unreadable json data", http.StatusBadRequest)
		return
	}

	uid, err := h.db.InsertUser(u)
	switch {
	case errors.Is(err, storage.ErrLoginNotUnique):
		log.Println("Handler.RegisterUser error:", u, err)
		http.Error(w, storage.ErrLoginNotUnique.Error(), http.StatusConflict)
		return
	case err != nil:
		log.Println("Handler.RegisterUser error:", u, err)
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	if err := writeAuthCookie(w, uid); err != nil {
		log.Println("Handler.RegisterUser error at building jwt-cookie:", uid, err)
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) AuthUser(w http.ResponseWriter, r *http.Request) {

	var in models.UserInsert
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "unreadable json data", http.StatusBadRequest)
		return
	}

	u := h.db.GetUser(in.Login)
	if u.ID == 0 || u.Password != in.Password {
		log.Println("Handler.AuthUser error: wrong login or password for", in)
		http.Error(w, "wrong login or password", http.StatusUnauthorized)
		return
	}

	if err := writeAuthCookie(w, u.ID); err != nil {
		log.Println("Handler.AuthUser error at building jwt-cookie:", u, err)
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
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
