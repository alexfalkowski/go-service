package status

import (
	"errors"
	"fmt"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Coder allows errors to implement so we can return the code needed.
type Coder interface {
	// Code reflects the status code to return, e.g: http.StatusNotFound.
	Code() int
}

// Taken from https://github.com/grpc-ecosystem/grpc-gateway/blob/main/runtime/errors.go
var statusCodes = map[codes.Code]int{
	codes.OK:                 http.StatusOK,
	codes.Canceled:           499,
	codes.Unknown:            http.StatusInternalServerError,
	codes.InvalidArgument:    http.StatusBadRequest,
	codes.DeadlineExceeded:   http.StatusGatewayTimeout,
	codes.NotFound:           http.StatusNotFound,
	codes.AlreadyExists:      http.StatusConflict,
	codes.PermissionDenied:   http.StatusForbidden,
	codes.Unauthenticated:    http.StatusUnauthorized,
	codes.ResourceExhausted:  http.StatusTooManyRequests,
	codes.FailedPrecondition: http.StatusBadRequest,
	codes.Aborted:            http.StatusConflict,
	codes.OutOfRange:         http.StatusBadRequest,
	codes.Unimplemented:      http.StatusNotImplemented,
	codes.Internal:           http.StatusInternalServerError,
	codes.Unavailable:        http.StatusInternalServerError,
	codes.DataLoss:           http.StatusInternalServerError,
}

// WriteError will write the error to the response writer.
func WriteError(res http.ResponseWriter, err error) {
	status := Code(err)

	http.Error(res, err.Error(), status)
}

// Error representing code and msg.
func Error(code int, msg string) error {
	return &statusError{code: code, msg: msg}
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

func (s *statusError) Error() string {
	return s.msg
}
