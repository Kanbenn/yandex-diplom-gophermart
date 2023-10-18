package handler

import (
	"io"
	"net/http"

	"github.com/Kanbenn/gophermart/internal/models"
	"github.com/Kanbenn/gophermart/internal/storage"
	"github.com/Kanbenn/gophermart/pkg/luhn"
)

func (h *Handler) PostOrder(w http.ResponseWriter, r *http.Request) {
	expectedContentType := r.Header.Get("Content-Type") == "text/plain"
	number, err := io.ReadAll(r.Body)
	if err != nil || !expectedContentType {
		http.Error(w, "неверный формат запроса;", http.StatusBadRequest)
		return
	}
	if valid := luhn.IsValidLuhnNumber(number); !valid {
		http.Error(w, "неверный формат номера заказа;", http.StatusUnprocessableEntity)
		return
	}
	uid := r.Context().Value(models.CtxKeyUser).(int)
	order := models.OrderInsert{Number: string(number), UID: uid}

	e := h.db.InsertOrder(order)
	w.WriteHeader(statusFromInsertResult(e))
}

func (h *Handler) GetOrders(w http.ResponseWriter, r *http.Request) {

}

func statusFromInsertResult(e error) int {
	switch e {
	case nil:
		return http.StatusAccepted
	case storage.ErrOrderWasPostedByThisUser:
		return http.StatusOK
	case storage.ErrOrderWasPostedByAnotherUser:
		return http.StatusConflict
	default:
		// storage.ErrUnxpectedError
		return http.StatusInternalServerError
	}
}
