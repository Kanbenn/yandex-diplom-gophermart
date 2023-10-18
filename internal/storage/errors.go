package storage

import "errors"

var (
	// user errors
	ErrLoginNotUnique = errors.New("this login is already taken")
	ErrUserUnknown    = errors.New("unknown user")
	ErrNotAuthorized  = errors.New("unauthorized")

	// insertOrder errors
	ErrOrderWasPostedByThisUser    = errors.New("this order number was already posted")
	ErrOrderWasPostedByAnotherUser = errors.New("this order number was already posted by other user")

	ErrUnxpectedError = errors.New("unexpected server error")
	ErrBadRequest     = errors.New("bad request")
)
