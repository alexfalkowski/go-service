package options

import (
	"fmt"
	"math"
	"strconv"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/runtime"
	"github.com/alexfalkowski/go-service/v2/time"
)

// Map contains string key-value pairs used to represent transport- or feature-specific options.
//
// It is commonly embedded into larger configuration structs to allow passing implementation-specific
// knobs without expanding the strongly-typed schema.
//
// Option values are treated as trusted startup configuration. Helpers panic for invalid present values
// so low-level knob mistakes fail fast during startup rather than surfacing on request paths.
type Map map[string]string

// Duration returns a duration option for key if present; otherwise it returns fallback.
//
// The stored value must be a Go duration string accepted by [time.ParseDuration] (for example "250ms",
// "30s", or "5m").
func (m Map) Duration(key string, fallback time.Duration) time.Duration {
	if val, ok := m[key]; ok {
		return time.MustParseDuration(val)
	}
	return fallback
}

// NonNegativeDuration returns a duration option for key and panics if the resolved value is negative.
//
// It behaves like Duration for parsing and fallback resolution, then enforces
// that startup configuration cannot disable timeout-like protections with a
// negative duration.
func (m Map) NonNegativeDuration(key string, fallback time.Duration) time.Duration {
	duration := m.Duration(key, fallback)
	if duration < 0 {
		runtime.Must(fmt.Errorf("options: %s must be non-negative: %s", key, duration))
	}

	return duration
}

// Uint32 returns an unsigned integer option for key if present; otherwise it returns fallback.
//
// The stored value must be a base-10 unsigned integer accepted by [strconv.ParseUint] for a 32-bit size.
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
// The stored value must be a human-readable decimal size string accepted by [bytes.ParseSize], such
// as "64B", "2MB", or "4GB".
//
// Size only parses the option value. It does not enforce [bytes.MaxConfigSize];
// typed config fields that require that cap use the config package's
// `config_size` validation rule after decoding.
func (m Map) Size(key string, fallback bytes.Size) bytes.Size {
	if val, ok := m[key]; ok {
		return bytes.MustParseSize(val)
	}

	return fallback
}

// IntSize returns a byte-size option as an int.
//
// It resolves the value with Size, so it panics if the option value cannot be
// parsed as a size or if the resolved size overflows int. It does not apply
// [bytes.MaxConfigSize].
func (m Map) IntSize(key string, fallback bytes.Size) int {
	size := m.Size(key, fallback)
	if size.Bytes() > math.MaxInt {
		runtime.Must(fmt.Errorf("options: %s exceeds max int: %s", key, size))
	}

	return int(size.Bytes())
}

// Int32Size returns a byte-size option as an int32.
//
// It resolves the value with Size, so it panics if the option value cannot be
// parsed as a size or if the resolved size overflows int32. It does not apply
// [bytes.MaxConfigSize].
func (m Map) Int32Size(key string, fallback bytes.Size) int32 {
	size := m.Size(key, fallback)
	if size.Bytes() > math.MaxInt32 {
		runtime.Must(fmt.Errorf("options: %s exceeds max int32: %s", key, size))
	}

	//nolint:gosec // Size is range-checked against MaxInt32 above.
	return int32(size.Bytes())
}

// Uint32Size returns a byte-size option as a uint32.
//
// It resolves the value with Size, so it panics if the option value cannot be
// parsed as a size or if the resolved size overflows uint32. It does not apply
// [bytes.MaxConfigSize].
func (m Map) Uint32Size(key string, fallback bytes.Size) uint32 {
	size := m.Size(key, fallback)
	if size.Bytes() > math.MaxUint32 {
		runtime.Must(fmt.Errorf("options: %s exceeds max uint32: %s", key, size))
	}

	//nolint:gosec // Size is range-checked against MaxUint32 above.
	return uint32(size.Bytes())
}
