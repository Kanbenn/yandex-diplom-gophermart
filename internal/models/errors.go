package models

import "errors"

var (
	ErrLoginNotUnique    = errors.New("this login is already taken")
	ErrUserUnknown       = errors.New("unknown user")
	ErrNotAuthorized     = errors.New("unauthorized")
	ErrNotEnoughMinerals = errors.New("user's balance doesn't have enough bonus points")

	ErrOrderWasPostedByThisUser    = errors.New("this order number was already posted")
	ErrOrderWasPostedByAnotherUser = errors.New("this order number was already posted by other user")

	ErrLuhnFormulaViolation = errors.New("неверный формат номера заказа;")
	ErrNoContent            = errors.New("no content")
	ErrUnxpectedError       = errors.New("unexpected server error")
)
