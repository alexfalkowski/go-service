package limiter

import "github.com/alexfalkowski/go-service/v2/time"

// Decision describes the result of one limiter take attempt.
type Decision struct {
	header     string
	resetAfter time.Duration
	allowed    bool
}

// Allowed reports whether the limited operation may proceed.
func (d Decision) Allowed() bool {
	return d.allowed
}

// Header returns the RateLimit header value for this decision.
func (d Decision) Header() string {
	return d.header
}

// ResetAfter returns the remaining duration before the limiter bucket resets.
func (d Decision) ResetAfter() time.Duration {
	return d.resetAfter
}

// ResetAfterSeconds returns ResetAfter rounded up to whole seconds.
func (d Decision) ResetAfterSeconds() uint64 {
	if d.resetAfter <= 0 {
		return 0
	}

	seconds := uint64(d.resetAfter / time.Second)
	if d.resetAfter%time.Second != 0 {
		seconds++
	}

	return seconds
}
