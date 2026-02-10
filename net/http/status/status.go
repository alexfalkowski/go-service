package status

import (
	"fmt"

	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"google.golang.org/grpc/status"
)

// Error representing code and msg.
func Error(code int, msg string) error {
	return &statusError{code: code, msg: msg}
}

// InternalServerError for status.
func InternalServerError(err error) error {
	return FromError(http.StatusInternalServerError, err)
}

// ServiceUnavailableError for status.
func ServiceUnavailableError(err error) error {
	return FromError(http.StatusServiceUnavailable, err)
}

// UnauthorizedError for status.
func UnauthorizedError(err error) error {
	return FromError(http.StatusUnauthorized, err)
}

// BadRequestError for status.
func BadRequestError(err error) error {
	return FromError(http.StatusBadRequest, err)
}

// FromError creates an error from an error.
func FromError(code int, err error) error {
	if IsError(err) {
		return err
	}

	return Error(code, err.Error())
}

// Errorf representing code and a formatted message.
func Errorf(code int, format string, a ...any) error {
	return Error(code, fmt.Sprintf(format, a...))
}

// IsError verifies if the package created this error.
func IsError(err error) bool {
	if _, ok := err.(Coder); ok {
		return true
	}

	e := &statusError{}

	return errors.As(err, &e)
}

// Code from the error. If nil 200, otherwise 500.
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
