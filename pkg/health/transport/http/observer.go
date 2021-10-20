package http

import (
	"github.com/alexfalkowski/go-health/pkg/subscriber"
)

// HealthObserver for HTTP.
type HealthObserver struct {
	*subscriber.Observer
}
