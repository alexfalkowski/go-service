// Package health provides health server wiring for go-service.
//
// This package integrates the go-health server with the application lifecycle by constructing a
// `*server.Server` (from github.com/alexfalkowski/go-health/v2/server) and starting/stopping it using
// Fx/Dig lifecycle hooks.
//
// # Registrations and checks
//
// The go-health server exposes health endpoints based on registrations and checkers managed by the
// go-health package. This package provides the `Registrations` alias to make it easier to pass around
// lists of health check registrations in go-service wiring.
//
// Health check implementations and constructors live under `health/checker`.
//
// Start with `Module` and `NewServer`.
package health
