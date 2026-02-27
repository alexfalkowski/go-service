package http

import (
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/net/http"
)

// Pattern builds a debug route pattern for name and pattern.
//
// This is a thin wrapper around `net/http.Pattern`, which prefixes the provided pattern with
// the service name to form routes like:
//
//	/<name><pattern>
//
// Example:
//
//	Pattern(name, "/debug/pprof/") // -> "/my-service/debug/pprof/"
func Pattern(name env.Name, pattern string) string {
	return http.Pattern(name, pattern)
}

// NewServeMux constructs a new debug ServeMux.
//
// The returned mux is a small wrapper over go-service's `net/http.ServeMux` and is intended as
// the shared router where debug subpackages register their handlers (pprof, fgprof, statsviz,
// psutil, etc.).
func NewServeMux() *ServeMux {
	return &ServeMux{http.NewServeMux()}
}

// ServeMux wraps `net/http.ServeMux` for debug routing.
//
// It exists primarily to provide a stable type that can be injected via DI as "the debug mux".
type ServeMux struct {
	*http.ServeMux
}
