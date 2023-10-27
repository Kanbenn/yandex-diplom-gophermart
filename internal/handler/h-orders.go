package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/Kanbenn/gophermart/internal/models"
)

func (h *Handler) PostNewOrder(w http.ResponseWriter, r *http.Request) {
	expectedContentType := r.Header.Get("Content-Type") == "text/plain"
	number, err := io.ReadAll(r.Body)
	if err != nil || !expectedContentType {
		http.Error(w, "неверный формат запроса;", http.StatusBadRequest)
		return
	}
	uid := r.Context().Value(models.CtxKeyUser).(int)
	order := models.Order{Number: string(number), User: uid}
	e := h.app.OrderNew(order)
	w.WriteHeader(statusFromInsertNewOrderResult(e))
}

func (h *Handler) PostNewOrderWithBonus(w http.ResponseWriter, r *http.Request) {
	var order models.Order
	order.User = r.Context().Value(models.CtxKeyUser).(int)
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		log.Println("h.PostNewOrderWithBonus error at reading input json:", err)
		http.Error(w, "unreadable json data", http.StatusBadRequest)
		return
	}
	log.Println("h.PostNewOrderWithBonus got this json data for user:", order)

	err := h.app.OrderNewWithBonus(order)
	w.WriteHeader(statusFromResult(err))
}
