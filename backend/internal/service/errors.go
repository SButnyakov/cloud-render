package service

import "errors"

var (
	ErrUserAlreadyExists  = errors.New("user already exsits")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")

	ErrFailedToGetSubscription = errors.New("failed to get subscription")

	ErrPaymentTypeNotFound      = errors.New("failed to get payment type")
	ErrSubscriptionTypeNotFound = errors.New("failed to get subscription type")

	ErrExternalError = errors.New("external error")
)
