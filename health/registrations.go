package health

import health "github.com/alexfalkowski/go-health/v2/server"

// Registrations is an alias for a slice of go-health registrations.
//
// It is typically used to pass around multiple health check registrations.
type Registrations = []*health.Registration
