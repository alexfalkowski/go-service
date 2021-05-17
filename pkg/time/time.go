package time

import (
	"time"
)

const (
	// RFC3339 format.
	RFC3339 = time.RFC3339
)

// Now the current local time.
func Now() time.Time {
	return time.Now()
}

// Since the time elapsed since t.
func Since(t time.Time) time.Duration {
	return time.Since(t)
}

// ToMilliseconds from the duration.
func ToMilliseconds(duration time.Duration) int64 {
	return duration.Nanoseconds() / 1e6 // nolint:gomnd
}
