package time

import (
	"time"

	"github.com/alexfalkowski/go-service/runtime"
)

const (
	// Minute represents 60 seconds.
	Minute = time.Minute

	// Second represents 1 second.
	Second = time.Second

	// RFC3339 formats time to 2006-01-02T15:04:05Z07:00.
	RFC3339 = time.RFC3339
)

type (
	// Time represents an instant in time with nanosecond precision.
	Time = time.Time

	// Duration represents the elapsed time between two instants.
	Duration = time.Duration
)

var (

	// Now is just an alias to time.Now.
	Now = time.Now

	// Since is just an alias to time.Since.
	Since = time.Since

	// ParseDuration is just an alias to time.ParseDuration.
	ParseDuration = time.ParseDuration
)

// MustParseDuration for time.
func MustParseDuration(s string) time.Duration {
	t, err := ParseDuration(s)
	runtime.Must(err)

	return t
}
