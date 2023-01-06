package time

import (
	"time"
)

const (
	timeout = 30

	// Timeout as a general guidance of the maximum time any operation should take.
	Timeout = timeout * time.Second
)

// ToMilliseconds from the duration.
func ToMilliseconds(duration time.Duration) int64 {
	return duration.Nanoseconds() / 1e6
}
