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
// gRPC status to the runtime. LocalError marks locally produced client-control
// errors without changing their gRPC code, message, details, or wrapped causes.
// For authoritative upstream behavior, edge cases, and wire-level details,
// consult the upstream gRPC documentation for the version you vend.
//
// # Background: gRPC status errors
//
// In gRPC, an RPC returns either a successful response or an error. When an error
// is returned, gRPC expects it to be (or to be convertible to) a "status" error
// consisting of:
//
//   - a status code (codes.Code), and
//   - a human-readable message,
//   - optionally, structured details.
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
// Use New when a caller needs a Status value before converting it to an error,
// for example to attach structured details:
//
//	s, err := status.New(codes.Unavailable, "retry later").WithDetails(&status.RetryInfo{
//		RetryDelay: status.NewDuration(delay),
//	})
//	if err != nil {
//		return err
//	}
//	return s.Err()
//
// Use Errorf when the client-visible message should be formatted:
//
//	err := status.Errorf(codes.NotFound, "widget %q does not exist", name)
//
// Use SafeError when the internal cause should remain available through
// unwrapping but the client should receive a lowercase "grpc: " prefixed
// status message. Use SafeErrorf when adding formatted internal context around
// that cause:
//
//	return status.SafeErrorf(codes.Internal, err, "load tenant %s", tenantID)
//
// Return SafeError and SafeErrorf directly from gRPC handlers. Do not wrap the
// returned error with [fmt.Errorf]("%w") before returning it to gRPC: upstream
// status.FromError preserves wrapping context in the client-visible status
// message for wrapped status errors. Put internal context in the cause passed
// to SafeError or in the SafeErrorf format instead:
//
//	return status.SafeError(codes.Internal, fmt.Errorf("load tenant: %w", err))
//
// # Inspecting errors
//
// Use Code to extract the gRPC code from an error:
//
//	c := status.Code(err)
//
// Code delegates to google.golang.org/grpc/status.Code.
// Use IsLocalError when retry or client policy needs to distinguish a local
// rejection from an upstream status error carrying the same code.
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
//   - RetryInfo and NewDuration expose the protobuf detail types used by
//     go-service retry behavior without requiring transport packages to import
//     protobuf packages directly.
//
// # Client-visible messages
//
// gRPC status messages are sent to clients. Error and Errorf treat their
// messages as client-visible. SafeError and SafeErrorf send a lowercase
// "grpc: " prefixed status message while preserving an internal cause through
// Unwrap for inspection with [errors.Is] and [errors.As].
//
// # Non-goals
//
// This package intentionally exposes only the small subset of the upstream
// status API and structured detail types that go-service uses broadly: Code,
// FromError, New, Error, Errorf, SafeError, SafeErrorf, LocalError,
// IsLocalError, Status, RetryInfo, and NewDuration. Higher-level code that needs
// unrelated structured-detail helpers should add a narrow go-service alias here
// before reaching through to upstream packages from transport code.
package status
