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

	e := h.db.InsertNewOrder(order)
	w.WriteHeader(statusFromInsertNewOrderResults(e))
}

func (h *Handler) GetUserOrders(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value(models.CtxKeyUser).(int)
	orders, err := h.db.SelectUserOrders(uid)
	if err != nil {
		log.Println("h.GetOrders error from Pg:", orders, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(orders) < 1 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	out, err := json.Marshal(orders)
	if err != nil {
		log.Println("h.GetOrders error at writing json:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

func (h *Handler) PostNewBonusOrder(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value(models.CtxKeyUser).(int)
	log.Println("h.PostNewBonusOrder:", uid)
	var order models.OrderNew
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		log.Println("h.PostNewBonusOrder error at reading input json:", err)
		http.Error(w, "unreadable json data", http.StatusBadRequest)
		return
	}
	log.Println("h.PostNewBonusOrder got this json data:", order)
	if !luhn.IsValidLuhnNumber([]byte(order.Number)) {
		log.Println("h.PostNewBonusOrder luhn formula error:", order)
		http.Error(w, "неверный формат номера заказа;", http.StatusUnprocessableEntity)
		return
	}
	order.User = uid

	err := h.db.InsertNewBonusOrder(order)
	log.Println("h.PostNewBonusOrder result from pg.Insert:", order, err)
	w.WriteHeader(statusFromInsertNewBonusOrderResults(err))
}
