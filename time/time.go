package time

import (
	"math/rand"
	"time"
)

const (
	timeout    = 5
	maxTimeout = timeout * 3

	// Timeout as a general guidance of the maximum time any operation should take.
	Timeout = timeout * time.Second
)

// ToMilliseconds from the duration.
func ToMilliseconds(duration time.Duration) int64 {
	return duration.Nanoseconds() / 1e6 // nolint:gomnd
}

// RandomWaitTime from the timeout.
func RandomWaitTime() time.Duration {
	min := 1
	num := rand.Intn(maxTimeout-min) + min // #nosec G404

	return time.Duration(num) * time.Second
}
