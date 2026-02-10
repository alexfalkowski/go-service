// Package status provides helpers for working with gRPC status errors in go-service.
//
// This package re-exports small helpers from google.golang.org/grpc/status behind go-service types
// (for example using codes.Code) so callers can depend on go-service packages consistently.
package status
