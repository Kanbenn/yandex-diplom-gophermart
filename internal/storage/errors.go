package storage

import (
	"errors"

	"github.com/lib/pq"
)

var (
	ErrLoginNotUnique    = errors.New("this login is already taken")
	ErrUserUnknown       = errors.New("unknown user")
	ErrNotAuthorized     = errors.New("unauthorized")
	ErrNotEnoughMinerals = errors.New("user's balance doesn't have enough bonus points")

	ErrOrderWasPostedByThisUser    = errors.New("this order number was already posted")
	ErrOrderWasPostedByAnotherUser = errors.New("this order number was already posted by other user")

	ErrUnxpectedError = errors.New("unexpected server error")
)

func isNotUniqueInsert(e error) bool {
	if err, ok := e.(*pq.Error); ok && err.Code == "23505" {
		return true
	}
	return false
}

func isNotEnoughMinerals(e error) bool {
	if err, ok := e.(*pq.Error); ok && err.Code == "23514" {
		return true
	}
	return false
}

func isForeignKeyViolation(e error) bool {
	if err, ok := e.(*pq.Error); ok && err.Code == "23503" {
		return true
	}
	return false
}
