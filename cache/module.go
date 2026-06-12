package cache

import (
	"github.com/alexfalkowski/go-service/v2/cache/driver"
	"github.com/alexfalkowski/go-service/v2/di"
)

// Module wires the cache subsystem into [go.uber.org/fx].
//
// It provides, in order:
//   - a cache [driver.Driver] (see [driver.NewDriver])
//   - an optional cache [Pinger] for backends that support connectivity checks
//   - a *[Cache] (see [NewCache])
//   - package-level registration (see [Register]) so generic helpers ([Get]/[Persist]) can be used
//
// # Disabled behavior
//
// When caching is disabled via configuration, [driver.NewDriver] returns a nil [driver.Driver] and [NewCache]
// returns a nil *[Cache]. [Register] is still invoked with nil, which makes the package-level helpers behave as if
// caching is disabled (no-ops / zero values) rather than failing.
var Module = di.Module(
	di.Constructor(driver.NewDriver),
	di.Constructor(NewPinger),
	di.Constructor(NewCache),
	di.Register(Register),
)
