package http

import (
	"fmt"
	"net/http"
)

// ResponseWriter with status for http.
type ResponseWriter struct {
	StatusCode int

	http.ResponseWriter
}

// WriteHeader sends an HTTP response header with the provided status code.
func (r *ResponseWriter) WriteHeader(statusCode int) {
	r.StatusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

// Status of the response, e.g. "200 OK".
func (r *ResponseWriter) Status() string {
	return fmt.Sprintf("%d %s", r.StatusCode, http.StatusText(r.StatusCode))
}
