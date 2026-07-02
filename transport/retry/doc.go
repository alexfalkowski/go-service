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
//   - Timeout: a transport-specific timeout knob. Transports that support per-attempt deadlines bound each attempt independently.
//   - Backoff: a delay inserted between failed attempts.
//
// Timeout and Backoff are typed durations. In config files they are encoded as
// Go duration strings (see [time.ParseDuration]), such as "250ms", "5s", or "1m".
//
// # Transport interpretation
//
// Because this package does not implement retries directly, the exact semantics depend on
// the transport implementation. In general, transports typically apply the configuration
// like this:
//
//   - Attempts limits how many times the operation will be tried.
//   - Timeout may be applied per attempt, depending on the transport.
//   - Backoff is applied between attempts, often as a fixed sleep; some transports may
//     apply additional logic (for example jitter) on top of this base value.
//
// Transports may also define what constitutes a retryable failure (for example, specific
// HTTP status codes, gRPC status codes, or transport-level errors) and may respect
// cancellation/deadlines from the caller's context regardless of Config.
//
// # Defaults and disabling
//
// This package defines shared defaults for common knobs. Callers should consult the
// transport packages to understand any transport-specific behavior.
//
// A common convention is:
//   - Attempts == 0: retries disabled using the transport's zero-value behavior.
//   - Attempts == 1: retries disabled.
//   - Timeout == 0: treated as "unspecified" and replaced with [time.DefaultTimeout] by transports that use it.
//   - Backoff == 0: treated as "unspecified" and replaced with DefaultBackoff.
//
// # Usage
//
// Embed Config into larger service configuration when callers need shared retry mechanics.
// Transport-specific retry configuration, such as HTTP or gRPC client retry config, embeds
// Config and adds protocol-specific failure classification before it is passed to transport
// wiring.
//
// Start with Config in config.go.
package retry
