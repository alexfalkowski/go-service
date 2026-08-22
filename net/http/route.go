package http

import (
	"net/http"

	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/net/url"
	"github.com/alexfalkowski/go-service/v2/strings"
)

// ParseServiceMethod derives a logical "service" and "method" name from an HTTP request.
//
// This helper is intended for consistent telemetry naming. It attempts to derive names from the request
// path when it follows the conventional go-service route shape:
//
//	/<service>/<method>
//
// If the request path matches that shape, ParseServiceMethod returns the extracted service/method pair.
//
// Otherwise it falls back to:
//   - method: lower-cased HTTP method (e.g. "get", "post")
//   - service: a best-effort name derived from the path:
//   - "root" when the path is empty or "/"
//   - otherwise the path without the leading "/" (e.g. "/health" -> "health")
func ParseServiceMethod(req *http.Request) (string, string) {
	path := req.URL.Path
	if service, method, ok := url.SplitPath(path); ok {
		return service, method
	}

	method := strings.ToLower(req.Method)

	if strings.IsEmpty(path) {
		return "root", method
	}

	path = path[1:]
	if strings.IsEmpty(path) {
		return "root", method
	}

	return path, method
}

// Pattern constructs a route pattern of the form "/<name><pattern>".
//
// This helper is used to namespace routes by service name so different services can share a router/mux
// without colliding, and so route names are consistent across telemetry, server registration, and tests.
//
// Example:
//
//	Pattern(name, "/debug/pprof/") // -> "/my-service/debug/pprof/"
func Pattern(name env.Name, pattern string) string {
	return strings.Concat("/", name.String(), pattern)
}
