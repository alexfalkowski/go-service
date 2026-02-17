package time

import (
	"time"

	"github.com/alexfalkowski/go-service/v2/runtime"
)

const (
	// Hour is an alias of time.Hour.
	Hour = time.Hour

	// Microsecond is an alias of time.Microsecond.
	Microsecond = time.Microsecond

	// Millisecond is an alias of time.Millisecond.
	Millisecond = time.Millisecond

	// Minute is an alias of time.Minute.
	Minute = time.Minute

	// Nanosecond is an alias of time.Nanosecond.
	Nanosecond = time.Nanosecond

	// Second is an alias of time.Second.
	Second = time.Second

	// RFC3339 is an alias of time.RFC3339.
	RFC3339 = time.RFC3339
)

// Time is an alias for time.Time.
type Time = time.Time

// Duration is an alias for time.Duration.
type Duration = time.Duration

// Now is just an alias to time.Now.
func Now() Time {
	return time.Now()
}

// ParseDuration is just an alias to time.ParseDuration.
func ParseDuration(s string) (Duration, error) {
	return time.ParseDuration(s)
}

// Since is just an alias to time.Since.
func Since(t Time) Duration {
	return time.Since(t)
}

// Sleep is just an alias to time.Sleep.
func Sleep(d Duration) {
	time.Sleep(d)
}

// MustParseDuration parses s as a duration and panics if parsing fails.
func MustParseDuration(s string) time.Duration {
	t, err := ParseDuration(s)
	runtime.Must(err)
	return t
}
