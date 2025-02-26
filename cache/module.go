package cache

import (
	"github.com/alexfalkowski/go-service/cache/cachego"
	"go.uber.org/fx"
)

// Module for fx.
var Module = fx.Options(
	fx.Provide(cachego.New),
	fx.Provide(New),
	fx.Invoke(Register),
)
