package cache

import (
	"github.com/alexfalkowski/go-service/v2/cache/driver"
	"github.com/alexfalkowski/go-service/v2/context"
)

// Pinger verifies cache backend connectivity.
type Pinger interface {
	// Ping verifies the cache backend is reachable.
	Ping(ctx context.Context) error
}

// NewPinger extracts a ping capability from a configured cache driver.
//
// Drivers without backend connectivity checks return nil.
func NewPinger(d driver.Driver) Pinger {
	pinger, _ := d.(Pinger)

	return pinger
}
