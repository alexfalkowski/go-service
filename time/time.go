package time

import (
	"time"

	"github.com/alexfalkowski/go-service/v2/runtime"
)

const (
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

type (
	// Time is an alias of time.Time.
	Time = time.Time

	// Duration is an alias of time.Duration.
	Duration = time.Duration
)

var (
	// Now is just an alias to time.Now.
	Now = time.Now

	// ParseDuration is just an alias to time.ParseDuration.
	ParseDuration = time.ParseDuration

	// Since is just an alias to time.Since.
	Since = time.Since

	// Sleep is just an alias to time.Sleep.
	Sleep = time.Sleep
)

// MustParseDuration for time.
func MustParseDuration(s string) time.Duration {
	t, err := ParseDuration(s)
	runtime.Must(err)

	return t
}
