package time

import (
	"math/rand"
	"time"
)

const (
	timeout = 5

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
	max := timeout * 3              // nolint:gomnd
	num := rand.Intn(max-min) + min // #nosec G404

	return time.Duration(num) * time.Second
}

// SleepRandomWaitTime from the timeout.
func SleepRandomWaitTime() {
	time.Sleep(RandomWaitTime())
}
