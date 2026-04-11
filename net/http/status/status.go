package status

import (
	"fmt"

	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"google.golang.org/grpc/status"
)

// Error constructs an error that carries an HTTP status code and a message.
//
// The returned error implements the Coder interface so handlers and helpers can extract the HTTP
// status code via Code(err).
func Error(code int, msg string) error {
	return &statusError{code: code, msg: msg}
}

// InternalServerError wraps err with StatusInternalServerError (500) unless err already carries a status code.
//
// This is a convenience wrapper over FromError.
func InternalServerError(err error) error {
	return FromError(http.StatusInternalServerError, err)
}

// ServiceUnavailableError wraps err with StatusServiceUnavailable (503) unless err already carries a status code.
//
// This is a convenience wrapper over FromError.
func ServiceUnavailableError(err error) error {
	return FromError(http.StatusServiceUnavailable, err)
}

// UnauthorizedError wraps err with StatusUnauthorized (401) unless err already carries a status code.
//
// This is a convenience wrapper over FromError.
func UnauthorizedError(err error) error {
	return FromError(http.StatusUnauthorized, err)
}

// BadRequestError wraps err with StatusBadRequest (400) unless err already carries a status code.
//
// This is a convenience wrapper over FromError.
func BadRequestError(err error) error {
	return FromError(http.StatusBadRequest, err)
}

// FromError returns err unchanged if it already carries a status code; otherwise it wraps err with code.
//
// This helper is intentionally idempotent for errors already produced by this package (or any error
// implementing Coder): calling FromError on such errors does not overwrite the original status code.
// Raw *http.MaxBytesError values are normalized to StatusRequestEntityTooLarge (413) regardless of the
// provided code so oversized request bodies surface consistently.
//
// Note: err must be non-nil. Passing a nil error will panic because err.Error() will be called.
func FromError(code int, err error) error {
	if IsError(err) {
		return err
	}

	if isMaxBytesError(err) {
		code = http.StatusRequestEntityTooLarge
	}

	return Error(code, err.Error())
}

// Errorf formats a message and returns an error with the provided status code.
//
// This is a convenience wrapper over Error(code, fmt.Sprintf(...)).
func Errorf(code int, format string, a ...any) error {
	return Error(code, fmt.Sprintf(format, a...))
}

// IsError reports whether err carries a status code.
//
// It returns true for:
//   - errors produced by this package, and
//   - any error implementing the Coder interface.
func IsError(err error) bool {
	_, ok := coderFromError(err)
	return ok
}

// Code extracts the HTTP status code from err.
//
// Resolution order:
//  1. If err implements Coder, return coder.Code().
//  2. If err is a raw *http.MaxBytesError, return StatusRequestEntityTooLarge (413).
//  3. If err is a gRPC status error, map its gRPC code to an HTTP status code using the statusCodes table.
//  4. Otherwise return StatusInternalServerError (500).
func Code(err error) int {
	if coder, ok := coderFromError(err); ok {
		return coder.Code()
	}

	if isMaxBytesError(err) {
		return http.StatusRequestEntityTooLarge
	}

	s, ok := status.FromError(err)
	if ok {
		return statusCodes[s.Code()]
	}

	return http.StatusInternalServerError
}

func coderFromError(err error) (Coder, bool) {
	var coder Coder
	if errors.As(err, &coder) {
		return coder, true
	}

	return nil, false
}

func isMaxBytesError(err error) bool {
	_, ok := errors.AsType[*http.MaxBytesError](err)
	return ok
}

type statusError struct {
	msg  string
	code int
}

// Code returns the HTTP status code carried by the error.
func (s *statusError) Code() int {
	return s.code
}

// Error returns the status message.
func (s *statusError) Error() string {
	return s.msg
}
