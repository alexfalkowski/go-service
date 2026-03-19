package test

import (
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/net/http/status"
)

var (
	// ErrGenerate is returned by token helpers that simulate generation failures.
	ErrGenerate = errors.New("token: generation issue")

	// ErrInvalid is returned by helpers that simulate invalid credentials or payloads.
	ErrInvalid = errors.New("token: invalid match")

	// ErrFailed is the generic failure sentinel used by test doubles in this package.
	ErrFailed = errors.New("failed")

	// ErrInternal is a test error that also exposes an HTTP status code.
	ErrInternal = &internalError{}

	_ status.Coder = ErrInternal
)

type internalError struct{}

// Error implements the error interface.
func (e *internalError) Error() string {
	return "internal"
}

// Code implements status.Coder.
func (e *internalError) Code() int {
	return 500
}
