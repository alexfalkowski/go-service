// Package status provides helpers for constructing and inspecting gRPC status errors
// while using go-service types.
//
// This package is a thin wrapper around google.golang.org/grpc/status. It exists
// primarily to:
//
//   - Standardize imports across go-service codebases (use net/grpc/status and
//     net/grpc/codes instead of importing upstream gRPC packages everywhere).
//   - Ensure status construction and inspection uses the go-service codes.Code
//     alias (net/grpc/codes), which is itself an alias of gRPC's codes.Code.
//
// The functions in this package do not add new semantics beyond forwarding to
// the upstream gRPC status package. For authoritative behavior, edge cases, and
// wire-level details, consult the upstream gRPC documentation for the version
// you vend.
//
// # Background: gRPC status errors
//
// In gRPC, an RPC returns either a successful response or an error. When an error
// is returned, gRPC expects it to be (or to be convertible to) a "status" error
// consisting of:
//
//   - a status code (codes.Code), and
//   - a human-readable message,
//   - optionally, structured details (not handled by this package).
//
// Clients can inspect errors to extract the status code and decide how to react
// (retry, surface a user-facing error, map to HTTP status codes, etc.).
//
// # Constructing errors
//
// Use Error or Errorf to create an error with a specific gRPC code:
//
//	err := status.Error(codes.NotFound, "widget does not exist")
//
//	err := status.Errorf(codes.InvalidArgument, "bad widget id: %q", id)
//
// These functions delegate to google.golang.org/grpc/status.Error and
// google.golang.org/grpc/status.Errorf.
//
// # Inspecting errors
//
// Use Code to extract the gRPC code from an error:
//
//	c := status.Code(err)
//
// Code delegates to google.golang.org/grpc/status.Code.
//
// If err is nil, Code returns codes.OK (matching upstream behavior). If err is
// not a gRPC status error, upstream logic typically returns codes.Unknown. Treat
// this as a signal that the error did not originate from a standard gRPC status
// (for example it may be a local transport error, a context error, or an
// application error that wasn't converted).
//
// # Relationship to other packages
//
//   - net/grpc/codes provides the codes.Code alias and the well-known gRPC code
//     constants (OK, NotFound, Unavailable, etc.).
//   - google.golang.org/grpc/status is the upstream implementation; this package
//     forwards directly to it.
//
// # Non-goals
//
// This package intentionally does not expose the full status API surface (for
// example status.FromError, status.New, or Status.Details). Higher-level code
// that needs structured details should use the upstream status package directly.
package status
