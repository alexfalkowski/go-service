// Package config defines the go-service TLS configuration model and constructor
// helpers.
//
// This package is the configuration-facing counterpart to `crypto/tls`.
// It contains:
//   - the lightweight `Config` struct used to describe certificate and key
//     source strings in decoded application config, and
//   - `NewConfig`, which resolves those source strings and materializes a
//     runtime `*crypto/tls.Config`.
//
// Use this package in higher-level config and transport code. Use `crypto/tls`
// directly when code already has a runtime TLS config and only needs low-level
// TLS types or helpers.
package config
