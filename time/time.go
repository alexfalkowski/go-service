package time

import (
	"math/rand"
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

// RandomWaitTime from the timeout.
func RandomWaitTime() time.Duration {
	min := 1
	num := rand.Intn(timeout-min) + min // #nosec G404

	return time.Duration(num) * time.Second
}
