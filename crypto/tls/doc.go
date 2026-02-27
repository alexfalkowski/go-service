// Package tls provides helpers for constructing `crypto/tls` configurations from go-service config.
//
// This package focuses on turning go-service TLS configuration (`tls.Config`) into a standard library
// `*crypto/tls.Config` via `NewConfig`.
//
// # Defaults
//
// `NewConfig` applies conservative defaults suitable for service-to-service deployments:
//
//   - Minimum TLS version: 1.2
//   - Client authentication: `tls.RequireAndVerifyClientCert` (mTLS)
//
// If you need different semantics (for example optional client certs, server-only TLS, or TLS 1.3-only),
// construct a `*crypto/tls.Config` directly or wrap the returned config and adjust fields as needed.
package tls
