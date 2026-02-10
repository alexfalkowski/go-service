// Package hooks provides HTTP webhook hook middleware and wiring for go-service.
//
// This package integrates the Standard Webhooks signing and verification flow into HTTP transports.
// It provides:
//   - server-side verification middleware (rejecting invalid webhook requests), and
//   - client-side signing middleware (adding webhook signature headers).
//
// Start with `NewHandler` for server-side verification and `NewRoundTripper` for client-side signing.
package hooks
