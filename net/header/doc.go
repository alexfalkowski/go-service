// Package header provides shared helpers for working with network protocol headers in go-service.
//
// This package currently focuses on Bearer Authorization header parsing and shared forwarding header names so
// HTTP and gRPC integrations can use the same metadata semantics.
//
// Start with `ParseBearer` and `ForwardedIPs`.
package header
