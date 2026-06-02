// Package header provides shared helpers for working with network protocol headers in go-service.
//
// This package currently focuses on Bearer Authorization header parsing and shared forwarding header names so
// HTTP and gRPC integrations can use the same metadata semantics.
//
// Forwarded IP headers are accepted as trusted inputs by transport metadata extraction. Services that rely on
// extracted IPs for access logs, policy, or rate limiting should only receive traffic through trusted edge
// infrastructure that strips or overwrites client-supplied forwarding headers.
//
// Start with [ParseBearer] and [ForwardedIPs].
package header
