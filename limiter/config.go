package limiter

// Config configures the in-memory rate limiter.
//
// The limiter is typically constructed by `NewLimiter`, which interprets these fields as:
//
//   - Kind: selects how request keys are derived from context metadata (see NewKeyMap for default kinds).
//     Examples: "user-agent", "ip", "token".
//
//   - Interval: a Go duration string (e.g. "1s", "1m") defining the refill/measurement window used by the
//     underlying token bucket store.
//
//   - Tokens: the maximum number of tokens available per Interval for a given key.
//
// # Optional pointers and "enabled" semantics
//
// This config is intentionally optional. By convention across go-service configuration types, a nil
// *Config is treated as "limiter disabled" (see IsEnabled).
type Config struct {
	// Kind selects which limiter key kind is used for limiting (see NewKeyMap for default kinds).
	//
	// If Kind is not present in the KeyMap passed to NewLimiter, NewLimiter returns ErrMissingKey.
	Kind string `yaml:"kind,omitempty" json:"kind,omitempty" toml:"kind,omitempty"`

	// Interval is a Go duration string (for example "1s" or "1m") that defines the refill window.
	//
	// NewLimiter parses this value using time.MustParseDuration; invalid values will panic during
	// limiter construction.
	Interval string `yaml:"interval,omitempty" json:"interval,omitempty" toml:"interval,omitempty"`

	// Tokens is the maximum number of tokens available per interval, per derived key.
	//
	// When Tokens is 0, the underlying store behavior is implementation-defined (it may reject all
	// requests or behave as an always-empty bucket). Prefer configuring a positive value.
	Tokens uint64 `yaml:"tokens,omitempty" json:"tokens,omitempty" toml:"tokens,omitempty"`
}

// IsEnabled reports whether limiter configuration is present.
//
// By convention, a nil *Config is treated as "disabled".
func (c *Config) IsEnabled() bool {
	return c != nil
}
