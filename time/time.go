package time

import (
	"time"

	"github.com/alexfalkowski/go-service/v2/runtime"
)

// Hour is a duration constant equal to 60 minutes.
//
// It is an alias of time.Hour, provided so callers can depend on go-service
// packages while using standard library time values.
const Hour = time.Hour

// Microsecond is a duration constant equal to 1e3 nanoseconds.
//
// It is an alias of time.Microsecond.
const Microsecond = time.Microsecond

// Millisecond is a duration constant equal to 1e6 nanoseconds.
//
// It is an alias of time.Millisecond.
const Millisecond = time.Millisecond

// Minute is a duration constant equal to 60 seconds.
//
// It is an alias of time.Minute.
const Minute = time.Minute

// Nanosecond is a duration constant equal to 1.
//
// It is an alias of time.Nanosecond.
const Nanosecond = time.Nanosecond

// Second is a duration constant equal to 1e9 nanoseconds.
//
// It is an alias of time.Second.
const Second = time.Second

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
// It is a type alias of time.Duration, meaning it has identical semantics and
// method set to the standard library type.
type Duration = time.Duration

// After waits for the duration to elapse and then sends the current time
// on the returned channel.
//
// This is a thin wrapper around time.After and does not change semantics.
func After(d Duration) <-chan Time {
	return time.After(d)
}

// NewTicker returns a new [Ticker] containing a channel that will send the current time on the channel after each tick.
//
// This is a thin wrapper around time.NewTicker and does not change semantics.
func NewTicker(d Duration) *Ticker {
	return time.NewTicker(d)
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
	return time.ParseDuration(s)
}

// Since returns the time elapsed since t.
//
// This is a thin wrapper around time.Since and does not change semantics.
func Since(t Time) Duration {
	return time.Since(t)
}

// Sleep pauses the current goroutine for at least the duration d.
//
// This is a thin wrapper around time.Sleep and does not change semantics.
func Sleep(d Duration) {
	time.Sleep(d)
}

// MustParseDuration parses s as a duration string and panics if parsing fails.
//
// This helper is intended for strict startup/configuration paths where an
// invalid duration is considered a fatal configuration/programming error. It
// panics by calling runtime.Must on the parse error.
//
// If you need recoverable error handling, use ParseDuration instead.
func MustParseDuration(s string) time.Duration {
	t, err := ParseDuration(s)
	runtime.Must(err)
	return t
}
