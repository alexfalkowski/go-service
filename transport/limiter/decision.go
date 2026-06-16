package limiter

import "github.com/alexfalkowski/go-service/v2/time"

// Decision describes the result of one limiter take attempt.
type Decision struct {
	header       string
	policyHeader string
	resetAfter   time.Duration
	allowed      bool
}

// Allowed reports whether the limited operation may proceed.
func (d Decision) Allowed() bool {
	return d.allowed
}

// Header returns the RateLimit header value for this decision.
func (d Decision) Header() string {
	return d.header
}

// PolicyHeader returns the RateLimit-Policy header value for this decision.
func (d Decision) PolicyHeader() string {
	return d.policyHeader
}

// ResetAfter returns the remaining duration before the limiter bucket resets.
func (d Decision) ResetAfter() time.Duration {
	return d.resetAfter
}

// ResetAfterSeconds returns ResetAfter rounded up to whole seconds.
func (d Decision) ResetAfterSeconds() uint64 {
	return durationSeconds(d.resetAfter)
}
