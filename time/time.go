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

// RFC3339 is the RFC3339 time format layout.
//
// It is an alias of time.RFC3339.
const RFC3339 = time.RFC3339

// Ticker is the go-service ticker type used across the repository.
//
// It is a type alias of time.Ticker, meaning it has identical semantics and method
// set to the standard library type.
type Ticker = time.Ticker

// Time is the go-service time type used across the repository.
//
// It is a type alias of time.Time, meaning it has identical semantics and method
// set to the standard library type.
type Time = time.Time

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

// After waits for the duration to elapse and then sends the current time
// on the returned channel.
//
// This is a thin wrapper around time.After and does not change semantics.
func After(d Duration) <-chan Time {
	return time.After(d.Duration())
}

// NewTicker returns a new [Ticker] containing a channel that will send the current time on the channel after each tick.
//
// This is a thin wrapper around time.NewTicker and does not change semantics.
func NewTicker(d Duration) *Ticker {
	return time.NewTicker(d.Duration())
}

// Now returns the current local time.
//
// This is a thin wrapper around time.Now and does not change semantics.
func Now() Time {
	return time.Now()
}

// ParseDuration parses a duration string.
//
// This is a thin wrapper around time.ParseDuration. The input uses the standard
// Go duration format such as "250ms", "5s", or "1m".
func ParseDuration(s string) (Duration, error) {
	d, err := time.ParseDuration(s)
	return Duration(d), err
}

// Since returns the time elapsed since t.
//
// This is a thin wrapper around time.Since and does not change semantics.
func Since(t Time) Duration {
	return Duration(time.Since(t))
}

// Sleep pauses the current goroutine for at least the duration d.
//
// This is a thin wrapper around time.Sleep and does not change semantics.
func Sleep(d Duration) {
	time.Sleep(d.Duration())
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
