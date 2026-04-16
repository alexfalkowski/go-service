// Package tls exposes selected standard library crypto/tls types and helpers
// through the go-service import path.
//
// This package is the low-level runtime TLS surface for the repository. It
// intentionally preserves standard library semantics while allowing internal
// packages to depend on a go-service import path instead of importing
// `crypto/tls` directly.
//
// Use this package when code needs runtime TLS values such as `*tls.Config`,
// parsed certificates, protocol-version constants, or helpers such as
// `X509KeyPair`.
//
// Use `crypto/tls/config` when working with go-service TLS configuration
// values that resolve certificate and key "source strings" and materialize a
// runtime TLS config via `config.NewConfig`.
package tls
