package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/Kanbenn/gophermart/internal/models"
	"github.com/Kanbenn/gophermart/internal/storage"
)

func (h *Handler) RegisterNewUser(w http.ResponseWriter, r *http.Request) {

	var u models.UserInsert
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "unreadable json data", http.StatusBadRequest)
		return
	}

	uid, err := h.db.InsertNewUser(u)
	switch {
	case errors.Is(err, storage.ErrLoginNotUnique):
		http.Error(w, storage.ErrLoginNotUnique.Error(), http.StatusConflict)
		return
	case err != nil:
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	log.Println("Handler.RegisterUser:", u, err)

	if err := writeAuthCookie(w, uid); err != nil {
		log.Println("Handler.RegisterUser error at building jwt-cookie:", uid, err)
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {

	var in models.UserInsert
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "unreadable json data", http.StatusBadRequest)
		return
	}

	u := h.db.SelectUserAuth(in.Login)
	if u.ID < 1 || u.Password != in.Password {
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

func (h *Handler) GetUserAllOrders(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value(models.CtxKeyUser).(int)
	orders, err := h.db.SelectUserAllOrders(uid)
	if err != nil {
		log.Println("h.GetOrders error from Pg:", orders, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(orders) < 1 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	writeJsnResponse(w, orders)
}

func (h *Handler) GetUserBalance(w http.ResponseWriter, r *http.Request) {

	uid := r.Context().Value(models.CtxKeyUser).(int)

	orders, err := h.db.SelectUserBalance(uid)
	if err != nil {
		log.Println("h.GetUserBalance error:", uid, orders, err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	writeJsnResponse(w, orders)
}

func (h *Handler) GetUserWithdrawHistory(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value(models.CtxKeyUser).(int)

	orders, err := h.db.SelectUserWithdrawHistory(uid)
	if err != nil {
		log.Println("h.GetUserHistory err:", orders, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(orders) < 1 {
		log.Println("h.GetUserHistory err: нет ни одного списания", uid, orders, err)
		w.WriteHeader(http.StatusNoContent)
		return
	}
	writeJsnResponse(w, orders)
}
