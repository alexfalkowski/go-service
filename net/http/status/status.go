package status

import (
	"fmt"

	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"google.golang.org/grpc/status"
)

// Error returns an error that carries an HTTP status code and message.
func Error(code int, msg string) error {
	return &statusError{code: code, msg: msg}
}

// InternalServerError wraps err with StatusInternalServerError.
func InternalServerError(err error) error {
	return FromError(http.StatusInternalServerError, err)
}

// ServiceUnavailableError wraps err with StatusServiceUnavailable.
func ServiceUnavailableError(err error) error {
	return FromError(http.StatusServiceUnavailable, err)
}

// UnauthorizedError wraps err with StatusUnauthorized.
func UnauthorizedError(err error) error {
	return FromError(http.StatusUnauthorized, err)
}

// BadRequestError wraps err with StatusBadRequest.
func BadRequestError(err error) error {
	return FromError(http.StatusBadRequest, err)
}

// FromError returns err if it already carries a status code, otherwise it wraps err with code.
func FromError(code int, err error) error {
	if IsError(err) {
		return err
	}

	return Error(code, err.Error())
}

// Errorf formats a message and returns an error with the provided status code.
func Errorf(code int, format string, a ...any) error {
	return Error(code, fmt.Sprintf(format, a...))
}

// IsError reports whether err carries a status code created by this package.
func IsError(err error) bool {
	if _, ok := err.(Coder); ok {
		return true
	}

	e := &statusError{}

	return errors.As(err, &e)
}

// Code extracts the HTTP status code from err.
//
// It checks Coder, the package status type, and gRPC status mappings. If no mapping
// is found, StatusInternalServerError is returned.
func Code(err error) int {
	if coder, ok := err.(Coder); ok {
		return coder.Code()
	}

	e := &statusError{}
	if errors.As(err, &e) {
		return e.code
	}

	s, ok := status.FromError(err)
	if ok {
		return statusCodes[s.Code()]
	}

	return http.StatusInternalServerError
}

type statusError struct {
	msg  string
	code int
}

// Error returns the status message.
func (s *statusError) Error() string {
	return s.msg
}
