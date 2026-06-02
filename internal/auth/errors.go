package auth

import "errors"

var (
	ErrUserEmailAlreadyExists = errors.New("email already exist")
	ErrUserNotFound           = errors.New("user not found")
	ErrInvalidPassword        = errors.New("invalid password")
)
