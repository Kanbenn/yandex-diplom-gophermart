package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/Kanbenn/gophermart/internal/models"
	"github.com/Kanbenn/gophermart/pkg/luhn"
)

func (h *Handler) PostNewOrder(w http.ResponseWriter, r *http.Request) {
	expectedContentType := r.Header.Get("Content-Type") == "text/plain"
	number, err := io.ReadAll(r.Body)
	if err != nil || !expectedContentType {
		http.Error(w, "неверный формат запроса;", http.StatusBadRequest)
		return
	}
	if !luhn.IsValidLuhnNumber(number) {
		http.Error(w, "неверный формат номера заказа;", http.StatusUnprocessableEntity)
		return
	}
	uid := r.Context().Value(models.CtxKeyUser).(int)
	order := models.OrderNew{Number: string(number), User: uid}

	e := h.db.InsertOrder(order)
	w.WriteHeader(statusFromInsertOrderResults(e))
}

func (h *Handler) PostNewOrderWithBonus(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value(models.CtxKeyUser).(int)
	log.Println("h.PostNewOrderWithBonus:", uid)
	var order models.OrderNew
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		log.Println("h.PostNewOrderWithBonus error at reading input json:", err)
		http.Error(w, "unreadable json data", http.StatusBadRequest)
		return
	}
	log.Println("h.PostNewOrderWithBonus got this json data:", order)
	if !luhn.IsValidLuhnNumber([]byte(order.Number)) {
		log.Println("h.PostNewOrderWithBonus luhn formula error:", order)
		http.Error(w, "неверный формат номера заказа;", http.StatusUnprocessableEntity)
		return
	}
	order.User = uid
	order.Status = "PROCESSED"

	err := h.db.InsertOrderWithBonus(order)
	log.Println("h.PostNewOrderWithBonus result from pg.Insert:", order, err)
	w.WriteHeader(statusFromInsertOrderWithBonusResults(err))
}
