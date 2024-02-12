package repository

import "errors"

const (
	packagePath = "repository."
)

var (
	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("invalid credentials")

	ErrSubscriptionNotFound = errors.New("subscription not found")
)
