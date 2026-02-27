package health

import health "github.com/alexfalkowski/go-health/v2/server"

// Registrations is a convenience alias for a slice of go-health registrations.
//
// It is typically used in DI wiring to aggregate and pass around multiple health check registrations
// that should be installed on a shared `*server.Server` (from github.com/alexfalkowski/go-health/v2/server).
//
// This type does not add behavior; it exists to make dependencies that operate on sets of registrations
// more self-documenting.
type Registrations = []*health.Registration
