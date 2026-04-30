package authmodel

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid login name or password")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserDisabled       = errors.New("user disabled")
	ErrTokenBlocked       = errors.New("token blocked")
	ErrInvalidToken       = errors.New("invalid token")
)
