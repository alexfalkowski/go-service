// Package config defines the go-service TLS configuration model and constructor
// helpers.
//
// This package is the configuration-facing counterpart to `crypto/tls`.
// It contains:
//   - the lightweight `Config` struct used to describe certificate, key, CA,
//     and server-name settings in decoded application config, and
//   - helper functions that resolve those source strings into runtime TLS key
//     material.
//
// Use `config/server.NewConfig` and `config/client.NewConfig` when code needs a
// runtime `*crypto/tls.Config` for a specific TLS role. Use `crypto/tls`
// directly when code already has a runtime TLS config and only needs low-level
// TLS types or helpers.
package config
