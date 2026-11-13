package cache

import (
	"github.com/alexfalkowski/go-service/v3/cache/driver"
	"github.com/alexfalkowski/go-signal/v2"
)

// Module wires the cache subsystem into [github.com/alexfalkowski/go-signal/v2].
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
var Module = signal.Module(
	signal.Constructor(driver.NewDriver),
	signal.Constructor(NewPinger),
	signal.Constructor(NewCache),
	signal.Register(Register),
)
