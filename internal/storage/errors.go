package storage

import "errors"

var (
	ErrLoginNotUnique = errors.New("this login is already taken")
	ErrUserUnknown    = errors.New("unknown user")
	ErrNotAuthorized  = errors.New("unauthorized")
	ErrBadRequest     = errors.New("bad request")
)
