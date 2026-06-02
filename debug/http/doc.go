// Package http provides debug HTTP routing helpers for go-service.
//
// This package provides a small wrapper around go-service's [github.com/alexfalkowski/go-service/v2/net/http] primitives to build a
// dedicated debug endpoint router. The root debug package registers its built-in
// endpoint handlers on the mux returned by [NewServeMux].
//
// # Routing conventions
//
// Handlers are typically registered using [Pattern], which prefixes the provided route pattern
// with the service name to form a stable debug route namespace.
//
// Start with [ServeMux], [NewServeMux], and [Pattern].
package http
