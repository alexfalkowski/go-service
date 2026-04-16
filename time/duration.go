package time

import (
	"time"

	"github.com/alexfalkowski/go-service/v2/encoding/json"
	"github.com/alexfalkowski/go-service/v2/runtime"
)

// Hour is the go-service representation of the standard library `time.Hour`
// constant.
//
// Its value is one hour.
const Hour Duration = Duration(time.Hour)

// Microsecond is the go-service representation of the standard library
// `time.Microsecond` constant.
//
// Its value is one microsecond.
const Microsecond Duration = Duration(time.Microsecond)

// Millisecond is the go-service representation of the standard library
// `time.Millisecond` constant.
//
// Its value is one millisecond.
const Millisecond Duration = Duration(time.Millisecond)

// Minute is the go-service representation of the standard library `time.Minute`
// constant.
//
// Its value is one minute.
const Minute Duration = Duration(time.Minute)

// Nanosecond is the go-service representation of the standard library
// `time.Nanosecond` constant.
//
// Its value is one nanosecond.
const Nanosecond Duration = Duration(time.Nanosecond)

// Second is the go-service representation of the standard library `time.Second`
// constant.
//
// Its value is one second.
const Second Duration = Duration(time.Second)

// Duration is the go-service duration type used across the repository.
//
// It is a named type over the standard library `time.Duration` so it can expose
// config-friendly text and JSON marshaling helpers while remaining easy to
// convert at API boundaries.
//
// Duration values serialize as Go duration strings such as `250ms`, `5s`, or
// `3m15s`.
type Duration time.Duration

// Duration converts d to the standard library `time.Duration` type.
//
// Use it when calling APIs that accept the standard library duration type.
func (d Duration) Duration() time.Duration {
	return time.Duration(d)
}

// String returns d in the standard Go duration format.
func (d Duration) String() string {
	return d.Duration().String()
}

// MarshalText encodes d using the same duration string returned by
// [Duration.String].
func (d Duration) MarshalText() ([]byte, error) {
	return []byte(d.String()), nil
}

// UnmarshalText parses a Go duration string into d.
//
// Accepted inputs use the same format as [ParseDuration], such as `250ms`,
// `5s`, or `1m`.
func (d *Duration) UnmarshalText(text []byte) error {
	duration, err := ParseDuration(string(text))
	if err != nil {
		return err
	}

	*d = duration
	return nil
}

// MarshalJSON encodes d as a quoted Go duration string, such as `"5s"`.
func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

// UnmarshalJSON decodes a quoted Go duration string into d.
//
// Non-string JSON values are rejected by the underlying JSON decoder before the
// duration parser runs.
func (d *Duration) UnmarshalJSON(data []byte) error {
	var text string
	if err := json.Unmarshal(data, &text); err != nil {
		return err
	}

	return d.UnmarshalText([]byte(text))
}

// ParseDuration parses a duration string.
//
// This is a thin wrapper around `time.ParseDuration` that returns the
// repository's [Duration] type. The input uses the standard Go duration format
// such as `250ms`, `5s`, or `1m`.
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
