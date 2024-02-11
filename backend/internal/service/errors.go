package service

import "errors"

var (
	ErrUserAlreadyExists  = errors.New("user already exsits")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
)
