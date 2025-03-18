package debug

import "net/http"

// NewServeMux creates a new ServeMux.
func NewServeMux() *ServeMux {
	return &ServeMux{http.NewServeMux()}
}

// ServeMux is a wrapper around http.ServeMux.
type ServeMux struct {
	*http.ServeMux
}
