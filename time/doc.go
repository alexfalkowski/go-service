// Package time provides time-related helpers, aliases, and optional network time providers
// used by go-service.
//
// This package serves two purposes:
//
//  1. Standard library compatibility via aliases.
//     It re-exports a small subset of the Go standard library time API so code across
//     go-service can consistently import go-service packages while still using the
//     underlying time semantics.
//
//     Exported aliases include Time (time.Time), Duration (time.Duration), common
//     duration constants (Second, Minute, Hour, etc.), and RFC3339.
//
//  2. Optional network time sourcing.
//     In environments where local wall-clock time may drift or needs stronger
//     guarantees, this package can construct a Network provider that fetches time
//     from external services (for example NTP or NTS).
//
// # Standard library aliases
//
// The following identifiers are thin wrappers/aliases of the standard library and do
// not change semantics:
//
//   - Time and Duration alias time.Time and time.Duration.
//   - Now, Since, Sleep, and ParseDuration forward to time.Now, time.Since,
//     time.Sleep, and time.ParseDuration respectively.
//   - Constants such as Second, Minute, Hour, and RFC3339 alias the standard library
//     values.
//
// Use these when you want to keep dependencies within the go-service module while
// remaining fully compatible with the standard library time types.
//
// # Strict helpers
//
// MustParseDuration parses a Go duration string (see time.ParseDuration) and panics
// if parsing fails (via runtime.Must).
//
// This is intended for strict startup/configuration code paths where an invalid
// duration is a programmer/configuration error and should fail fast. If you need
// recoverable error handling, use ParseDuration directly.
//
// # Network time providers
//
// The Network interface provides a single method:
//
//   - Now() (Time, error): returns the current time as reported by the provider.
//
// NewNetwork constructs a Network implementation based on *Config. Enablement is
// modeled by presence: a nil *Config is treated as disabled and NewNetwork returns
// (nil, nil).
//
// Config.Kind selects the provider implementation. This package currently supports:
//
//   - "ntp": Network Time Protocol (NTP).
//   - "nts": Network Time Security (NTS), which provides authenticated time as
//     defined by RFC 8915.
//
// If Config.Kind is not recognized, NewNetwork returns ErrNotFound.
//
// Provider implementations may wrap and prefix errors to provide clearer context
// (for example "ntp: ..." or "nts: ...").
//
// # Dependency injection (Fx)
//
// Module wires NewNetwork into Fx as a constructor so applications can optionally
// depend on a Network provider when configured.
//
// # Notes
//
// Network time providers require external connectivity and correct configuration
// (for example a valid server address). Services should treat network time as an
// optional dependency unless their operational requirements demand it.
package time
