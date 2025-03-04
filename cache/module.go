package cache

import (
	"github.com/alexfalkowski/go-service/cache/driver"
	"go.uber.org/fx"
)

// Module for fx.
var Module = fx.Options(
	fx.Provide(driver.NewDriver),
	fx.Provide(NewCache),
	fx.Invoke(Register),
)
