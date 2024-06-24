package http

import (
	"fmt"
	"net/http"
)

// NewServeMux for HTTP.
func NewServeMux(s *http.ServeMux) ServeMux {
	return &StandardServeMux{s}
}

// ServeMux for HTTP.
type ServeMux interface {
	// Handle a verb, pattern with the func.
	Handle(verb, pattern string, fn http.HandlerFunc) error

	// Handler from the mux.
	Handler() http.Handler
}

// NewServeMux for http.
func NewStandardServeMux() *http.ServeMux {
	return http.NewServeMux()
}

// StandardServeMux for HTTP.
type StandardServeMux struct {
	*http.ServeMux
}

func (s *StandardServeMux) Handle(verb, pattern string, fn http.HandlerFunc) error {
	s.HandleFunc(fmt.Sprintf("%s %s", verb, pattern), fn)

	return nil
}

func (s *StandardServeMux) Handler() http.Handler {
	return s.ServeMux
}
