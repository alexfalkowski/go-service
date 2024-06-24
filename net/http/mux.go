package http

import (
	"net/http"
)

// NewServeMux for http.
func NewServeMux() *http.ServeMux {
	return http.NewServeMux()
}
