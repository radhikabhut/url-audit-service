package errors

import "errors"

var (
	ErrInvalidURL       = errors.New("invalid url format")
	ErrURLTimeout       = errors.New("request to target url timed out")
	ErrConnectionFailed = errors.New("failed to connect to target url")
	ErrRateLimit        = errors.New("rate limit exceeded")
	ErrInternal         = errors.New("internal server error")
	ErrNotFound         = errors.New("resource not found")
)
