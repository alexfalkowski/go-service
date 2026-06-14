// Package checker provides health check implementations used by go-service.
//
// This package contains checker.Checker implementations (from [github.com/alexfalkowski/go-health/v2/checker])
// that can be registered with a go-health server to expose liveness/readiness-style endpoints.
//
// The checkers in this package are intended to be composed by services depending on which repository
// subsystems they use, such as configured database pools and cache drivers.
//
// Start with [CacheChecker], [NewCacheChecker], [DBChecker], and [NewDBChecker].
package checker
