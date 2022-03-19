package http

import (
	"github.com/alexfalkowski/go-health/subscriber"
)

// HealthObserver for HTTP.
type HealthObserver struct {
	*subscriber.Observer
}

// LivenessObserver for HTTP.
type LivenessObserver struct {
	*subscriber.Observer
}

// ReadinessObserver for HTTP.
type ReadinessObserver struct {
	*subscriber.Observer
}
