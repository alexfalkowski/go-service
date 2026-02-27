// Package codes provides go-service aliases for gRPC status codes.
//
// This package re-exports google.golang.org/grpc/codes behind go-service types
// and constants so callers can:
//   - depend on go-service packages consistently (instead of importing the
//     upstream gRPC package directly everywhere), and
//   - use a stable, narrowly-scoped API surface when working with status codes
//     across go-service subsystems (for example net/grpc/status).
//
// # What is a gRPC code?
//
// A gRPC status code is a canonical classification of the outcome of an RPC.
// Codes are used to communicate success/failure semantics across language and
// transport boundaries.
//
// In Go, codes are represented by the upstream type:
//
//	google.golang.org/grpc/codes.Code
//
// This package defines:
//
//	// Code is an alias for codes.Code.
//	type Code = codes.Code
//
// and re-exports selected code constants (for example OK, NotFound, Unavailable,
// etc.).
//
// # Usage
//
// This package is commonly used with the sibling package net/grpc/status:
//
//	err := status.Error(codes.NotFound, "widget does not exist")
//
//	code := status.Code(err) // returns a codes.Code
//
// # Notes
//
// This package intentionally contains aliases/constants only. It does not define
// new semantics, mapping logic, or wrappers beyond naming and import
// consolidation. For detailed behavioral descriptions of each code, consult the
// upstream gRPC documentation.
package codes
