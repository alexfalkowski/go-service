package time

import (
	"encoding/json"
	"time"

	"github.com/alexfalkowski/go-service/v2/runtime"
)

// Hour is a duration constant equal to 60 minutes.
//
// It is an alias of time.Hour, provided so callers can depend on go-service
// packages while using standard library time values.
const Hour Duration = Duration(time.Hour)

// Microsecond is a duration constant equal to 1e3 nanoseconds.
//
// It is an alias of time.Microsecond.
const Microsecond Duration = Duration(time.Microsecond)

// Millisecond is a duration constant equal to 1e6 nanoseconds.
//
// It is an alias of time.Millisecond.
const Millisecond Duration = Duration(time.Millisecond)

// Minute is a duration constant equal to 60 seconds.
//
// It is an alias of time.Minute.
const Minute Duration = Duration(time.Minute)

// Nanosecond is a duration constant equal to 1.
//
// It is an alias of time.Nanosecond.
const Nanosecond Duration = Duration(time.Nanosecond)

// Second is a duration constant equal to 1e9 nanoseconds.
//
// It is an alias of time.Second.
const Second Duration = Duration(time.Second)

// Duration is the go-service duration type used across the repository.
//
// It is a named type over the standard library time.Duration so it can expose
// config-friendly marshaling helpers while remaining easy to convert at API
// boundaries.
type Duration time.Duration

// Duration converts d to the standard library time.Duration type.
func (d Duration) Duration() time.Duration {
	return time.Duration(d)
}

// String returns the Go duration string for d.
func (d Duration) String() string {
	return d.Duration().String()
}

// MarshalText encodes d using the standard Go duration string format.
func (d Duration) MarshalText() ([]byte, error) {
	return []byte(d.String()), nil
}

// UnmarshalText decodes a Go duration string into d.
func (d *Duration) UnmarshalText(text []byte) error {
	duration, err := ParseDuration(string(text))
	if err != nil {
		return err
	}

	*d = duration
	return nil
}

// MarshalJSON encodes d as a quoted Go duration string.
func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

// UnmarshalJSON decodes a quoted Go duration string into d.
func (d *Duration) UnmarshalJSON(data []byte) error {
	var text string
	if err := json.Unmarshal(data, &text); err != nil {
		return err
	}

	return d.UnmarshalText([]byte(text))
}

// ParseDuration parses a duration string.
//
// This is a thin wrapper around time.ParseDuration. The input uses the standard
// Go duration format such as "250ms", "5s", or "1m".
func ParseDuration(s string) (Duration, error) {
	d, err := time.ParseDuration(s)
	return Duration(d), err
}

// MustParseDuration parses s as a duration string and panics if parsing fails.
//
// This helper is intended for strict startup/configuration paths where an
// invalid duration is considered a fatal configuration/programming error. It
// panics by calling runtime.Must on the parse error.
//
// If you need recoverable error handling, use ParseDuration instead.
func MustParseDuration(s string) Duration {
	t, err := ParseDuration(s)
	runtime.Must(err)
	return t
}
