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
// Most functions in this package forward to the upstream gRPC status package.
// SafeError adds go-service safe-message behavior while still exposing a normal
// gRPC status to the runtime. For authoritative upstream behavior, edge cases,
// and wire-level details, consult the upstream gRPC documentation for the
// version you vend.
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
// Use Error to create an error with a specific gRPC code:
//
//	err := status.Error(codes.NotFound, "widget does not exist")
//
// Use Errorf when the client-visible message should be formatted:
//
//	err := status.Errorf(codes.NotFound, "widget %q does not exist", name)
//
// Use SafeError when the internal cause should remain available through
// unwrapping but the client should receive the standard status text.
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
// # Client-visible messages
//
// gRPC status messages are sent to clients. Error and Errorf treat their
// messages as client-visible. SafeError sends the standard status text while
// preserving an internal cause through Unwrap for inspection with errors.Is and
// errors.As.
//
// # Non-goals
//
// This package intentionally exposes only the small subset of the upstream
// status API that go-service uses broadly: Code, FromError, Error, Errorf,
// SafeError, and the Status type alias. Higher-level code that needs additional
// constructors or richer structured-detail helpers should use the upstream
// status package directly.
package status
