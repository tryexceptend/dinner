package services

import "errors"

var (
	ErrAttemptLimitExceeded = errors.New("user attempt limit exceeded")
	ErrEmptyFood            = errors.New("food is empty")
)
