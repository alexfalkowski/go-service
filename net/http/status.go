package http

import (
	"errors"
	"net/http"
)

// Error representing code and msg.
func Error(code int, msg string) error {
	return &statusError{code: code, msg: msg}
}

// IsError verifies if the package created this error.
func IsError(err error) bool {
	e := &statusError{}

	return errors.As(err, &e)
}

// Code from the error, otherwise 500.
func Code(err error) int {
	if err == nil {
		return http.StatusOK
	}

	e := &statusError{}
	if errors.As(err, &e) {
		return e.code
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
