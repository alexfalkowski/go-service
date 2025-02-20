package test

import "errors"

var (
	// ErrGenerate for test.
	ErrGenerate = errors.New("token error")

	// ErrInvalid for test.
	ErrInvalid = errors.New("invalid match")

	// ErrFailed for test.
	ErrFailed = errors.New("failed")
)
