package handler

import (
	"net/http"

	"github.com/Kanbenn/gophermart/internal/models"
)

func statusFromInsertNewOrderResult(e error) int {
	switch e {
	case nil:
		return http.StatusAccepted
	case models.ErrLuhnFormulaViolation:
		return http.StatusUnprocessableEntity
	case models.ErrOrderWasPostedByThisUser:
		return http.StatusOK
	case models.ErrOrderWasPostedByAnotherUser:
		return http.StatusConflict
	case models.ErrUserUnknown:
		return http.StatusUnauthorized
	default:
		// models.ErrUnxpectedError
		return http.StatusInternalServerError
	}
}

func statusFromResult(err error) int {
	switch err {
	case nil:
		return http.StatusOK
	case models.ErrLuhnFormulaViolation:
		return http.StatusUnprocessableEntity
	case models.ErrNotEnoughMinerals:
		return http.StatusPaymentRequired
	case models.ErrNoContent:
		return http.StatusNoContent
	case models.ErrLoginNotUnique,
		models.ErrOrderWasPostedByThisUser,
		models.ErrOrderWasPostedByAnotherUser:
		return http.StatusConflict
	case models.ErrUserUnknown:
		return http.StatusUnauthorized
	default:
		// models.ErrUnxpectedError
		return http.StatusInternalServerError
	}
}
