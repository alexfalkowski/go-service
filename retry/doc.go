// Package retry provides shared retry configuration used across go-service transports.
//
// This package intentionally contains configuration only. It does not implement a retry
// algorithm itself. Concrete retry behavior is implemented by transport-specific packages
// (for example transport/http/retry and transport/grpc/retry), which interpret Config in
// a way that is appropriate for the underlying protocol.
//
// # Configuration model
//
// Config is designed to capture the three most common knobs for retries:
//
//   - Attempts: the maximum number of attempts for a single logical operation,
//     including the initial attempt (so Attempts=1 means "no retries").
//   - Timeout: a per-attempt timeout. Each attempt is bounded independently.
//   - Backoff: a delay inserted between failed attempts.
//
// Timeout and Backoff are encoded as Go duration strings (see time.ParseDuration),
// such as "250ms", "5s", or "1m".
//
// # Transport interpretation
//
// Because this package does not implement retries directly, the exact semantics depend on
// the transport implementation. In general, transports typically apply the configuration
// like this:
//
//   - Attempts limits how many times the operation will be tried.
//   - Timeout is applied per attempt (usually by deriving a context with a deadline).
//   - Backoff is applied between attempts, often as a fixed sleep; some transports may
//     apply additional logic (for example jitter) on top of this base value.
//
// Transports may also define what constitutes a retryable failure (for example, specific
// HTTP status codes, gRPC status codes, or transport-level errors) and may respect
// cancellation/deadlines from the caller's context regardless of Config.
//
// # Defaults and disabling
//
// This package does not define defaults. Callers should consult the transport packages
// to understand default values and how zero values are handled.
//
// A common convention is:
//   - Attempts == 0: treated as "unspecified" and replaced with a transport default.
//   - Attempts == 1: retries disabled.
//   - Timeout == "" or Backoff == "": treated as "unspecified" and replaced with
//     transport defaults.
//
// # Usage
//
// Embed retry.Config into a larger service configuration and pass it to the relevant
// transport wiring. For example, a service might expose retry configuration for outbound
// HTTP calls and gRPC calls using the same struct, while each transport package applies
// the details appropriate for that protocol.
//
// Start with Config in config.go.
package retry
