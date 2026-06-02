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

// SafeMessageError is an error test double that exposes a safe message.
type SafeMessageError struct {
	Message string
}

// Error implements the error interface.
func (e SafeMessageError) Error() string {
	return "internal"
}

// SafeMessage returns the configured safe message.
func (e SafeMessageError) SafeMessage() string {
	return e.Message
}

// EmptySafeMessageError is an error test double that unwraps another error but has no safe message.
type EmptySafeMessageError struct {
	Err error
}

// Error implements the error interface.
func (e EmptySafeMessageError) Error() string {
	return "internal"
}

// SafeMessage returns an empty safe message.
func (e EmptySafeMessageError) SafeMessage() string {
	return ""
}

// Unwrap returns the wrapped error.
func (e EmptySafeMessageError) Unwrap() error {
	return e.Err
}

// NilError is an error test double used for typed nil checks.
type NilError struct{}

// Error implements the error interface.
func (e *NilError) Error() string {
	return "nil"
}

// MessageError is an error test double that returns Message.
type MessageError struct {
	Message string
}

// Error implements the error interface.
func (e *MessageError) Error() string {
	return e.Message
}

// Stringer is an [fmt.Stringer] test double that returns Value.
type Stringer struct {
	Value string
}

// String implements [fmt.Stringer].
func (s *Stringer) String() string {
	return s.Value
}

// CoderError is an error test double that exposes an HTTP status code.
type CoderError struct {
	StatusCode int
}

// Error implements the error interface.
func (e *CoderError) Error() string {
	return "custom"
}

// Code returns the configured status code.
func (e *CoderError) Code() int {
	return e.StatusCode
}

type internalError struct{}

// Error implements the error interface.
func (e *internalError) Error() string {
	return "internal"
}

// Code implements [status.Coder].
func (e *internalError) Code() int {
	return 500
}
