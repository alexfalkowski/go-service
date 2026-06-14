package limiter

import "github.com/alexfalkowski/go-service/v2/time"

// DefaultMaxKeys is the default maximum number of caller-derived keys with independent limiter buckets.
const DefaultMaxKeys uint64 = 4096

// Config configures the in-memory rate limiter.
//
// The limiter is typically constructed by [NewLimiter], which interprets these fields as:
//
//   - Kind: selects how request keys are derived from context metadata (see NewKeyMap for default kinds).
//     Prefer "user-id" when authenticated identity is available, or "transport-service-method" for
//     per-operation quotas. Other examples include "service-method", "ip", and "user-agent".
//
//   - Interval: a typed duration encoded as a Go duration string (e.g. "1s", "1m")
//     defining the refill/measurement window used by the underlying token bucket store.
//
//   - Tokens: the maximum number of tokens available per Interval for a given key.
//
//   - MaxKeys: the maximum number of caller-derived keys that get independent buckets.
//     Additional distinct keys share one overflow bucket.
//
// # Optional pointers and "enabled" semantics
//
// This config is intentionally optional. By convention across go-service configuration types, a nil
// *[Config] is treated as "limiter disabled" (see IsEnabled).
type Config struct {
	// Kind selects which limiter key kind is used for limiting (see NewKeyMap for default kinds).
	//
	// If Kind is not present in the KeyMap passed to NewLimiter, NewLimiter returns ErrMissingKey.
	Kind string `yaml:"kind,omitempty" json:"kind,omitempty" toml:"kind,omitempty"`

	// Interval is a typed duration that defines the refill window.
	//
	// In config files it is encoded as a Go duration string, for example "1s" or "1m".
	Interval time.Duration `yaml:"interval,omitempty" json:"interval,omitempty" toml:"interval,omitempty" validate:"gte=0"`

	// Tokens is the maximum number of tokens available per interval, per derived key.
	//
	// When Tokens is 0, the underlying in-memory store applies its default of 1 token.
	// Prefer configuring a positive value so the quota is explicit in service config.
	Tokens uint64 `yaml:"tokens,omitempty" json:"tokens,omitempty" toml:"tokens,omitempty"`

	// MaxKeys is the maximum number of caller-derived keys that get independent buckets.
	//
	// Additional distinct keys share one overflow bucket so high-cardinality key floods cannot create
	// unbounded store entries. A zero value applies [DefaultMaxKeys].
	MaxKeys uint64 `yaml:"max_keys,omitempty" json:"max_keys,omitempty" toml:"max_keys,omitempty"`
}

// IsEnabled reports whether limiter configuration is present.
//
// By convention, a nil *[Config] is treated as "disabled".
func (c *Config) IsEnabled() bool {
	return c != nil
}

// GetMaxKeys returns the configured active key cap.
//
// A nil receiver or a zero value falls back to [DefaultMaxKeys].
func (c *Config) GetMaxKeys() uint64 {
	if c == nil || c.MaxKeys == 0 {
		return DefaultMaxKeys
	}

	return c.MaxKeys
}
