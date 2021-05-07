package http

import (
	"go.uber.org/fx"
)

var (
	// ServerModule for fx.
	ServerModule = fx.Options(fx.Invoke(Register), fx.Provide(NewMux))

	// ClientModule for fx.
	ClientModule = fx.Options(fx.Provide(NewClient))

	// Module for fx.
	Module = fx.Options(ServerModule, ClientModule)
)
