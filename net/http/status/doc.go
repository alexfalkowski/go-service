// Package status provides helpers for working with HTTP status codes and status errors in go-service.
//
// This package defines a Coder interface, error constructors (Error/Errorf and helpers like BadRequestError),
// and utilities to classify and extract status codes from errors (including mapping gRPC status errors to HTTP codes).
//
// Status error messages are diagnostic and client-visible when passed to WriteError. This is intentional for a
// framework package: callers should pass the concrete error they want clients to see, or map/wrap internal failures
// to sanitized status errors before writing the response.
package status
