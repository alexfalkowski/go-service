package test

import "github.com/alexfalkowski/go-service/v2/net/http"

// ErrResponseWriter for test.
type ErrResponseWriter struct {
	Code int
}

// Header is always empty.
func (w *ErrResponseWriter) Header() http.Header {
	return http.Header{}
}

// Write returns ErrFailed.
func (w *ErrResponseWriter) Write([]byte) (int, error) {
	return 0, ErrFailed
}

// WriteHeader stores code in the Code field.
func (w *ErrResponseWriter) WriteHeader(code int) {
	w.Code = code
}
