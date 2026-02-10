// Package status provides helpers for working with HTTP status codes and status errors in go-service.
//
// This package defines a Coder interface, error constructors (Error/Errorf and helpers like BadRequestError),
// and utilities to classify and extract status codes from errors (including mapping gRPC status errors to HTTP codes).
package status
