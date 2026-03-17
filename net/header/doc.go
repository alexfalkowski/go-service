// Package header provides shared helpers for working with network protocol headers in go-service.
//
// This package currently focuses on Authorization header parsing so both HTTP and gRPC integrations can
// share the same scheme validation and error semantics.
//
// Start with `ParseAuthorization`.
package header
