package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Kanbenn/gophermart/internal/models"
)

func (h *Handler) RegisterNewUser(w http.ResponseWriter, r *http.Request) {
	var u models.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "unreadable json data", http.StatusBadRequest)
		return
	}
	uid, err := h.app.UserRegister(u)
	if err != nil {
		w.WriteHeader(statusFromResult(err))
		return
	}

	if err := writeAuthCookie(w, uid); err != nil {
		log.Println("h.RegisterNewUser jwt-cookie error:", uid, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var in models.User
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "unreadable json data", http.StatusBadRequest)
		return
	}
	uid, err := h.app.UserAuth(in)
	if err != nil {
		w.WriteHeader(statusFromResult(err))
		return
	}

	if err := writeAuthCookie(w, uid); err != nil {
		log.Println("h.LoginUser jwt-cookie error:", uid, err)
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetUserOrders(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value(models.CtxKeyUser).(int)
	orders, err := h.app.UserOrders(uid)
	if err != nil {
		w.WriteHeader(statusFromResult(err))
		return
	}
	writeJsnResponse(w, orders)
}

func (h *Handler) GetUserBalance(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value(models.CtxKeyUser).(int)
	orders, err := h.app.UserBalance(uid)
	if err != nil {
		w.WriteHeader(statusFromResult(err))
		return
	}
	writeJsnResponse(w, orders)
}

func (h *Handler) GetUserWithdrawHistory(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value(models.CtxKeyUser).(int)
	orders, err := h.app.UserWithdrawHistory(uid)
	if err != nil {
		w.WriteHeader(statusFromResult(err))
		return
	}
	writeJsnResponse(w, orders)
}
