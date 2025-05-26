package test

import (
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/net/http/status"
)

var (
	// ErrGenerate for test.
	ErrGenerate = errors.New("token error")

	// ErrInvalid for test.
	ErrInvalid = errors.New("invalid match")

	// ErrFailed for test.
	ErrFailed = errors.New("failed")

	// ErrInternal for test.
	ErrInternal = &internalError{}

	_ status.Coder = ErrInternal
)

type internalError struct{}

func (e *internalError) Error() string {
	return "internal"
}

func (e *internalError) Code() int {
	return 500
}
