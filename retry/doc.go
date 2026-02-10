// Package retry provides retry configuration shared across go-service.
//
// This package intentionally only defines the Config type used by transport integrations.
// Concrete retry behavior is implemented by transport-specific packages (for example
// `transport/http/retry` and `transport/grpc/retry`).
//
// Start with `Config`.
package retry
