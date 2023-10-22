package handler

import (
	"net/http"

	"github.com/Kanbenn/gophermart/internal/storage"
)

func statusFromInsertNewOrderResults(e error) int {
	switch e {
	case nil:
		return http.StatusAccepted
	case storage.ErrOrderWasPostedByThisUser:
		return http.StatusOK
	case storage.ErrOrderWasPostedByAnotherUser:
		return http.StatusConflict
	case storage.ErrUserUnknown:
		return http.StatusUnauthorized
	default:
		// storage.ErrUnxpectedError
		return http.StatusInternalServerError
	}
}

func statusFromInsertNewBonusOrderResults(e error) int {
	switch e {
	case nil:
		return http.StatusOK
	case storage.ErrNotEnoughMinerals:
		return http.StatusPaymentRequired
	case storage.ErrOrderWasPostedByThisUser, storage.ErrOrderWasPostedByAnotherUser:
		return http.StatusConflict
	case storage.ErrUserUnknown:
		return http.StatusUnauthorized
	default:
		// storage.ErrUnxpectedError
		return http.StatusInternalServerError
	}
}
