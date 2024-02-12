package repository

import "errors"

const (
	packagePath = "repository."
)

var (
	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("invalid credentials")

	ErrSubscriptionNotFound = errors.New("subscription not found")

	ErrNoSubscriptionTypes = errors.New("no subscription types found")
	ErrNoPaymentTypes      = errors.New("no payment types found")
)
