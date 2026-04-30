package options

import (
	"strconv"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/alexfalkowski/go-service/v2/time"
)

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
//
// Design note: option values are treated as trusted startup configuration. Panics are intentional
// fail-fast behavior for invalid low-level knobs, not request-path error handling.
func (m Map) Duration(key string, fallback time.Duration) time.Duration {
	if val, ok := m[key]; ok {
		return time.MustParseDuration(val)
	}
	return fallback
}

// Uint32 returns an unsigned integer option for key if present; otherwise it returns fallback.
//
// The stored value must be a base-10 unsigned integer accepted by strconv.ParseUint for a 32-bit size.
//
// Failure mode: if key is present but the value cannot be parsed as a uint32, Uint32 will panic because
// it treats invalid startup/configuration values as fatal.
//
// Design note: option values are treated as trusted startup configuration. Panics are intentional
// fail-fast behavior for invalid low-level knobs, not request-path error handling.
func (m Map) Uint32(key string, fallback uint32) uint32 {
	if val, ok := m[key]; ok {
		num, err := strconv.ParseUint(val, 10, 32)
		runtime.Must(err)

		return uint32(num)
	}

	return fallback
}

// Size returns a byte-size option for key if present; otherwise it returns fallback.
//
// The stored value must be a human-readable SI size string accepted by bytes.ParseSize, such as "64B",
// "2MB", or "4GB".
//
// Failure mode: if key is present but the value cannot be parsed as a size, Size will panic because
// it uses bytes.MustParseSize. Use this helper only when the option values are validated or come from
// trusted configuration.
//
// Design note: option values are treated as trusted startup configuration. Panics are intentional
// fail-fast behavior for invalid low-level knobs, not request-path error handling.
func (m Map) Size(key string, fallback bytes.Size) bytes.Size {
	if val, ok := m[key]; ok {
		return bytes.MustParseSize(val)
	}

	return fallback
}
