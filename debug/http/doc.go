// Package http provides debug HTTP routing helpers for go-service.
//
// This package provides a small wrapper around go-service's `net/http` primitives to build a
// dedicated debug endpoint router. Other debug subpackages (for example `debug/pprof`,
// `debug/fgprof`, `debug/statsviz`, and `debug/psutil`) register their handlers on the mux
// returned by `NewServeMux`.
//
// # Routing conventions
//
// Handlers are typically registered using `Pattern`, which prefixes the provided route pattern
// with the service name to form a stable debug route namespace.
//
// Start with `ServeMux`, `NewServeMux`, and `Pattern`.
package http
