package domain

import "errors"

var (
	ErrNotFound                  = errors.New("NOT_FOUND")
	ErrInternalError             = errors.New("INTERNAL_ERROR")
	ErrBadRequest                = errors.New("BAD_REQUEST")
	ErrSubscriptionAlreadyExists = errors.New("SUBSCRIPTION_ALREADY_EXISTS")
)
