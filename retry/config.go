package retry

import "github.com/alexfalkowski/go-service/v2/time"

// Config configures retry behavior for an operation.
//
// This package defines configuration only; concrete retry behavior is implemented
// by transport-specific packages (for example transport/http/retry and
// transport/grpc/retry). As a result, the exact retry policy (what is retryable,
// whether jitter is applied, etc.) is transport-defined, but these fields are
// the common knobs most implementations interpret similarly.
//
// Timeout and Backoff are typed durations encoded as Go duration strings in config
// (see time.ParseDuration), such as "250ms", "5s", or "1m".
type Config struct {
	// Timeout is the per-attempt timeout duration.
	//
	// When interpreted by a transport, each attempt is typically bounded
	// independently (for example by deriving a per-attempt context deadline).
	//
	// Value encoding: Go duration string (for example "250ms", "5s").
	Timeout time.Duration `yaml:"timeout,omitempty" json:"timeout,omitempty" toml:"timeout,omitempty"`

	// Backoff is the delay between attempts after a failed attempt.
	//
	// Transports commonly treat this as a base backoff duration and may apply
	// additional behavior on top (for example jitter), depending on the
	// implementation.
	//
	// Value encoding: Go duration string (for example "100ms", "1s").
	Backoff time.Duration `yaml:"backoff,omitempty" json:"backoff,omitempty" toml:"backoff,omitempty"`

	// Attempts is the maximum number of attempts, including the initial attempt.
	//
	// The go-service transport convention is:
	//   - Attempts == 0 leaves retries disabled using the transport's zero-value behavior.
	//   - Attempts == 1 disables retries (single attempt only).
	//   - Attempts > 1 allows up to Attempts-1 retries after the initial attempt.
	Attempts uint64 `yaml:"attempts,omitempty" json:"attempts,omitempty" toml:"attempts,omitempty"`
}

// MaxAttempts returns the configured total attempt count, including the initial attempt.
//
// It preserves the zero value so transports that treat zero specially can retain their
// upstream behavior.
func (c *Config) MaxAttempts() uint64 {
	if c == nil {
		return 0
	}

	return c.Attempts
}

// MaxRetries returns the maximum number of retries after the initial attempt.
//
// For example:
//   - Attempts == 0 or 1 returns 0 retries.
//   - Attempts == 2 returns 1 retry.
//   - Attempts == 3 returns 2 retries.
func (c *Config) MaxRetries() uint64 {
	if c == nil || c.Attempts <= 1 {
		return 0
	}

	return c.Attempts - 1
}
