package test

import (
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/net/http/status"
)

var (
	// ErrGenerate for test.
	ErrGenerate = errors.New("token: generation issue")

	// ErrInvalid for test.
	ErrInvalid = errors.New("token: invalid match")

	// ErrFailed for test.
	ErrFailed = errors.New("failed")

	// ErrInternal for test.
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
