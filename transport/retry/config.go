package retry

import (
	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/alexfalkowski/go-service/v2/time"
)

// DefaultBackoff is the shared retry backoff used when [Config.Backoff] is unset.
const DefaultBackoff time.Duration = 100 * time.Millisecond

// DefaultJitterPercent is the shared retry jitter applied around the configured base backoff.
const DefaultJitterPercent uint64 = 20

// MaxAttempts is the maximum retry attempt count, including the initial attempt.
//
// Keep Config.Attempts validation in sync with this value.
const MaxAttempts uint64 = 10

// Config configures retry behavior for an operation.
//
// This package defines configuration only; concrete retry behavior is implemented
// by transport-specific packages (for example transport/http/retry and
// transport/grpc/retry). As a result, the exact retry policy (what is retryable,
// whether jitter is applied, etc.) is transport-defined, but these fields are
// the common knobs most implementations interpret similarly.
//
// Timeout and Backoff are typed durations encoded as Go duration strings in config
// (see [time.ParseDuration]), such as "250ms", "5s", or "1m".
type Config struct {
	// Strategy selects the backoff growth applied between attempts.
	//
	// Supported values are "constant" (the default), "exponential", and "fibonacci".
	// Backoff is the base duration: constant reuses it for every wait, while exponential
	// and fibonacci grow it on each attempt. An empty value applies the constant strategy.
	Strategy string `yaml:"strategy,omitempty" json:"strategy,omitempty" toml:"strategy,omitempty" validate:"omitempty,oneof=constant exponential fibonacci"`

	// Timeout is a transport-specific timeout duration.
	//
	// Transports that support per-attempt deadlines use it to bound attempts independently.
	// Other transports may rely on the caller's request context or client timeout instead.
	//
	// Value encoding: Go duration string (for example "250ms", "5s").
	// A zero value applies time.DefaultTimeout. Negative values are invalid.
	Timeout time.Duration `yaml:"timeout,omitempty" json:"timeout,omitempty" toml:"timeout,omitempty" validate:"gte=0"`

	// Backoff is the delay between attempts after a failed attempt.
	//
	// Transports commonly treat this as a base backoff duration and may apply
	// additional behavior on top (for example jitter), depending on the
	// implementation.
	//
	// Value encoding: Go duration string (for example "100ms", "1s").
	// A zero value applies DefaultBackoff. Negative values are invalid.
	Backoff time.Duration `yaml:"backoff,omitempty" json:"backoff,omitempty" toml:"backoff,omitempty" validate:"gte=0"`

	// MaxBackoff caps the per-attempt backoff duration for the exponential and fibonacci
	// strategies. Without a cap, those strategies grow the wait between attempts
	// unboundedly (aside from any overall context or client timeout).
	//
	// Value encoding: Go duration string (for example "5s", "30s").
	// A zero value leaves backoff growth uncapped, preserving today's behavior. Negative
	// values are invalid.
	MaxBackoff time.Duration `yaml:"max_backoff,omitempty" json:"max_backoff,omitempty" toml:"max_backoff,omitempty" validate:"gte=0"`

	// Attempts is the maximum number of attempts, including the initial attempt.
	//
	// The go-service transport convention is:
	//   - Attempts == 0 leaves retries disabled using the transport's zero-value behavior.
	//   - Attempts == 1 disables retries (single attempt only).
	//   - Attempts > 1 allows up to Attempts-1 retries after the initial attempt.
	//
	// Decoded configuration rejects values above MaxAttempts. Direct public API construction is
	// clamped by MaxAttempts and MaxRetries so transport callers cannot bypass the repository cap.
	Attempts uint64 `yaml:"attempts,omitempty" json:"attempts,omitempty" toml:"attempts,omitempty" validate:"lte=10"`
}

// MaxAttempts returns the configured total attempt count, including the initial attempt.
//
// It preserves the zero value so transports that treat zero specially can retain their
// upstream behavior. Values above MaxAttempts are clamped to MaxAttempts.
func (c *Config) MaxAttempts() uint64 {
	if c == nil {
		return 0
	}
	if c.Attempts > MaxAttempts {
		return MaxAttempts
	}

	return c.Attempts
}

// GetTimeout returns the configured transport-specific retry timeout.
//
// A nil receiver or a non-positive value falls back to [time.DefaultTimeout].
func (c *Config) GetTimeout() time.Duration {
	if c == nil || c.Timeout <= 0 {
		return time.DefaultTimeout
	}

	return c.Timeout
}

// GetBackoff returns the configured retry backoff.
//
// A nil receiver or a non-positive value falls back to DefaultBackoff.
func (c *Config) GetBackoff() time.Duration {
	if c == nil || c.Backoff <= 0 {
		return DefaultBackoff
	}

	return c.Backoff
}

// GetMaxBackoff returns the configured maximum backoff duration.
//
// A nil receiver or a non-positive value returns zero, meaning backoff growth is uncapped.
func (c *Config) GetMaxBackoff() time.Duration {
	if c == nil || c.MaxBackoff <= 0 {
		return 0
	}

	return c.MaxBackoff
}

// GetStrategy returns the configured retry backoff strategy.
//
// A nil receiver or an empty value falls back to the "constant" strategy.
func (c *Config) GetStrategy() string {
	if c == nil || strings.IsEmpty(c.Strategy) {
		return "constant"
	}

	return c.Strategy
}

// MaxRetries returns the maximum number of retries after the initial attempt.
//
// For example:
//   - Attempts == 0 or 1 returns 0 retries.
//   - Attempts == 2 returns 1 retry.
//   - Attempts == 3 returns 2 retries.
func (c *Config) MaxRetries() uint64 {
	attempts := c.MaxAttempts()
	if attempts <= 1 {
		return 0
	}

	return attempts - 1
}
