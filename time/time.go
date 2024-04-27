package time

import (
	"time"

	"github.com/alexfalkowski/go-service/runtime"
)

const (
	// Timeout as a general guidance of the maximum time any operation should take.
	Timeout = 30 * time.Second

	// Backoff as a general guidance of the scalar time any retry should take.
	Backoff = 100 * time.Millisecond
)

// MustParseDuration for time.
func MustParseDuration(s string) time.Duration {
	t, err := time.ParseDuration(s)
	runtime.Must(err)

	return t
}
