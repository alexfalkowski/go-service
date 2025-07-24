package http

import (
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/net/http"
)

// Pattern is an alias of http.Pattern.
func Pattern(name env.Name, pattern string) string {
	return http.Pattern(name, pattern)
}

// NewServeMux creates a new ServeMux.
func NewServeMux() *ServeMux {
	return &ServeMux{http.NewServeMux()}
}

// ServeMux is a composed of a http.ServeMux.
type ServeMux struct {
	*http.ServeMux
}
