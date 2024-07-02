package http

import (
	"errors"
	"net/http"
)

// Error representing code and msg.
func Error(code int, msg string) error {
	return &statusError{code: code, msg: msg}
}

// Code from the error, otherwise 500.
func Code(err error) int {
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
