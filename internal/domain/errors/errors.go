package errors

import "errors"

var (
	ErrInvalidInput = errors.New("invalid subscription input")
	ErrNotFound     = errors.New("subscription not found")
)
