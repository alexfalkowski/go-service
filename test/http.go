package test

import (
	"net/http"
)

// BadResponseWriter for test.
type BadResponseWriter struct {
	Code int
}

// Header is always empty.
func (w *BadResponseWriter) Header() http.Header {
	return http.Header{}
}

// Write returns ErrFailed.
func (w *BadResponseWriter) Write([]byte) (int, error) {
	return 0, ErrFailed
}

// WriteHeader stored the code in the field Code.
func (w *BadResponseWriter) WriteHeader(code int) {
	w.Code = code
}
