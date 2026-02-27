package options

import "github.com/alexfalkowski/go-service/v2/time"

// Map contains string key-value pairs used to represent transport- or feature-specific options.
//
// It is commonly embedded into larger configuration structs to allow passing implementation-specific
// knobs without expanding the strongly-typed schema.
type Map map[string]string

// Duration returns a duration option for key if present; otherwise it returns fallback.
//
// The stored value must be a Go duration string accepted by time.ParseDuration (for example "250ms",
// "30s", or "5m").
//
// Failure mode: if key is present but the value cannot be parsed as a duration, Duration will panic
// because it uses time.MustParseDuration. Use this helper only when the option values are validated
// or come from trusted configuration.
func (m Map) Duration(key string, fallback time.Duration) time.Duration {
	if val, ok := m[key]; ok {
		return time.MustParseDuration(val)
	}
	return fallback
}
