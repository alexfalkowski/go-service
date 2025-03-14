package debug

import "net/http"

// ServeMux is an alias for http.ServeMux.
type ServeMux struct {
	*http.ServeMux
}

// NewServeMux creates a new ServeMux.
func NewServeMux() *ServeMux {
	return &ServeMux{http.NewServeMux()}
}
