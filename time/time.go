package time

import (
	"time"

	"github.com/alexfalkowski/go-service/runtime"
)

// MustParseDuration for time.
func MustParseDuration(s string) time.Duration {
	t, err := time.ParseDuration(s)
	runtime.Must(err)

	return t
}
