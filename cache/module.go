package cache

import (
	"github.com/alexfalkowski/go-service/v2/cache/driver"
	"github.com/alexfalkowski/go-service/v2/di"
)

// Module wires cache driver, cache, and package-level registration into Fx.
var Module = di.Module(
	di.Constructor(driver.NewDriver),
	di.Constructor(NewCache),
	di.Register(Register),
)
