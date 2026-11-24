// Package time provides time-related helpers, aliases, and optional network time providers
// used by go-service.
//
// This package serves two purposes:
//
//  1. Standard library compatibility via wrappers and aliases.
//     It exposes standard Go time operations through the go-service import path.
//
//  2. Optional network time sourcing.
//     In environments where local wall-clock time may drift or needs stronger
//     guarantees, this package can construct a Network provider that fetches time
//     from external services (for example NTP or NTS).
//
// # Standard library compatibility
//
// [Time], [Location], [Timer], and [Ticker] alias the corresponding standard
// library types. [Duration] wraps [time.Duration] with text and JSON marshaling
// that uses Go duration strings.
//
// Use [Parse] and [ParseInLocation] with standard layouts such as [RFC3339]
// and [DateOnly]. Construct timestamps with [Date], [Unix] (whole seconds),
// [UnixMilli], or [UnixMicro]. Timer helpers such as [AfterFunc], [NewTimer],
// [NewTicker], and [Tick] accept the package's [Duration].
//
// # Strict helpers
//
// [MustParseDuration] parses a Go duration string (see [time.ParseDuration]) and
// panics if parsing fails (via [runtime.Must]).
//
// This is intended for strict startup/configuration code paths where an invalid
// duration is a programmer/configuration error and should fail fast. If you need
// recoverable error handling, use [ParseDuration] directly.
//
// # Network time providers
//
// The [Network] interface provides a single method:
//
//   - Now() (Time, error): returns the current time as reported by the provider.
//
// [NewNetwork] constructs a Network implementation based on [Config].
// Enablement is modeled by presence: a nil *[Config] is treated as disabled and
// NewNetwork returns (nil, nil). A non-nil *[Config] is enabled even when
// [Config.Kind] is empty; empty or unknown kinds return [ErrNotFound].
//
// [Config.Kind] selects the provider implementation. This package currently
// supports:
//
//   - "ntp": Network Time Protocol (NTP).
//   - "nts": Network Time Security (NTS), which provides authenticated time as
//     defined by RFC 8915.
//
// If [Config.Kind] is not recognized, [NewNetwork] returns [ErrNotFound].
//
// Provider implementations may wrap and prefix errors to provide clearer context
// (for example "ntp: ..." or "nts: ...").
//
// # Dependency injection (go-signal)
//
// [Module] wires [NewNetwork] into go-signal as a constructor so applications can optionally
// depend on a Network provider when configured.
//
// # Notes
//
// Network time providers require external connectivity and correct configuration
// (for example a valid server address). Services should treat network time as an
// optional dependency unless their operational requirements demand it.
package time
