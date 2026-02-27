// Package checker provides health check implementations used by go-service.
//
// This package contains `checker.Checker` implementations (from github.com/alexfalkowski/go-health/v2/checker)
// that can be registered with a go-health server to expose liveness/readiness-style endpoints.
//
// The checkers in this package are intended to be composed by services depending on which subsystems they use
// (database, caches, upstream dependencies, etc.).
//
// Start with the package-level constructors (for example `NewDBChecker`).
package checker
