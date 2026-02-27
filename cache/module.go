package cache

import (
	"github.com/alexfalkowski/go-service/v2/cache/driver"
	"github.com/alexfalkowski/go-service/v2/di"
)

// Module wires the cache subsystem into Fx.
//
// It provides, in order:
//   - a cache Driver (see cache/driver.NewDriver)
//   - a *Cache (see NewCache)
//   - package-level registration (see Register) so generic helpers (Get/Persist) can be used
//
// # Disabled behavior
//
// When caching is disabled via configuration, driver.NewDriver returns a nil Driver and NewCache returns
// a nil *Cache. Register is still invoked with nil, which makes the package-level helpers behave as if
// caching is disabled (no-ops / zero values) rather than failing.
var Module = di.Module(
	di.Constructor(driver.NewDriver),
	di.Constructor(NewCache),
	di.Register(Register),
)
