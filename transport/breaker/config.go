package breaker

import "github.com/alexfalkowski/go-service/v2/time"

// Config configures common circuit breaker behavior.
//
// This shared config only models transport-agnostic breaker mechanics. HTTP
// status-code and gRPC status-code failure classification belong to their
// transport-specific breaker config packages.
//
// A nil *Config is treated as "breaker disabled" by config-driven client
// wiring. When a Config is present, zero values keep [DefaultSettings].
type Config struct {
	// Interval is the closed-state interval used to clear breaker counts.
	//
	// In config files it is encoded as a Go duration string (for example "30s").
	// A zero value keeps the default. Negative values are invalid.
	Interval time.Duration `yaml:"interval,omitempty" json:"interval,omitempty" toml:"interval,omitempty" validate:"gte=0"`

	// Timeout is how long the breaker stays open before allowing half-open probes.
	//
	// In config files it is encoded as a Go duration string (for example "10s").
	// A zero value keeps the default. Negative values are invalid.
	Timeout time.Duration `yaml:"timeout,omitempty" json:"timeout,omitempty" toml:"timeout,omitempty" validate:"gte=0"`

	// MaxRequests is the maximum number of requests allowed while the breaker is half-open.
	//
	// A zero value keeps the default.
	MaxRequests uint32 `yaml:"max_requests,omitempty" json:"max_requests,omitempty" toml:"max_requests,omitempty"`

	// ConsecutiveFailures is the number of consecutive failures that opens the breaker.
	//
	// A zero value keeps the default.
	ConsecutiveFailures uint32 `yaml:"consecutive_failures,omitempty" json:"consecutive_failures,omitempty" toml:"consecutive_failures,omitempty"`

	// FailureRatio is the failure ratio (0 < r <= 1) that opens the breaker once MinRequests is reached.
	//
	// A zero value keeps ConsecutiveFailures-based tripping. When set, it takes precedence over ConsecutiveFailures.
	FailureRatio float64 `yaml:"failure_ratio,omitempty" json:"failure_ratio,omitempty" toml:"failure_ratio,omitempty" validate:"omitempty,gt=0,lte=1"`

	// MinRequests is the minimum request volume within the interval before FailureRatio is evaluated.
	//
	// A zero value means any request volume is eligible once FailureRatio is set.
	MinRequests uint32 `yaml:"min_requests,omitempty" json:"min_requests,omitempty" toml:"min_requests,omitempty"`
}

// IsEnabled reports whether breaker configuration is present.
func (c *Config) IsEnabled() bool {
	return c != nil
}

// Settings returns breaker settings using [DefaultSettings] plus configured overrides.
func (c *Config) Settings() Settings {
	settings := DefaultSettings
	if c == nil {
		return settings
	}

	if c.MaxRequests > 0 {
		settings.MaxRequests = c.MaxRequests
	}
	if c.Interval > 0 {
		settings.Interval = c.Interval.Duration()
	}
	if c.Timeout > 0 {
		settings.Timeout = c.Timeout.Duration()
	}
	if c.FailureRatio > 0 {
		ratio := c.FailureRatio
		minRequests := c.MinRequests
		settings.ReadyToTrip = func(counts Counts) bool {
			return counts.Requests > 0 && counts.Requests >= minRequests && float64(counts.TotalFailures)/float64(counts.Requests) >= ratio
		}
	} else if c.ConsecutiveFailures > 0 {
		failures := c.ConsecutiveFailures
		settings.ReadyToTrip = func(counts Counts) bool {
			return counts.ConsecutiveFailures >= failures
		}
	}

	return settings
}
