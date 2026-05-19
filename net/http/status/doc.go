// Package status provides helpers for working with HTTP status codes and status errors in go-service.
//
// This package defines a Coder interface, error constructors (Error/Errorf and helpers like BadRequestError),
// and utilities to classify and extract status codes from errors (including mapping gRPC status errors to HTTP codes).
//
// Status error messages created with Error/Errorf are client-visible when passed to WriteError. Wrapped
// internal failures created with FromError or SafeError keep their diagnostic Error text, but WriteError
// sends a lowercase "http: " prefixed status message instead.
package status
