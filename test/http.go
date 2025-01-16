package test

import (
	"net/http"
)

// BadResponseWriter for test.
type BadResponseWriter struct {
	Code int
}

func (w *BadResponseWriter) Header() http.Header {
	return http.Header{}
}

func (w *BadResponseWriter) Write([]byte) (int, error) {
	return 0, ErrFailed
}

func (w *BadResponseWriter) WriteHeader(code int) {
	w.Code = code
}
