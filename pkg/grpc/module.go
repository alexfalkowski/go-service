package grpc

import (
	"go.uber.org/fx"
)

var (
	// ServerModule for fx.
	ServerModule = fx.Provide(NewServer)

	// ServerOptionsModule for fx.
	ServerOptionsModule = fx.Provide(NewServerOptions)

	// Module for fx.
	Module = fx.Options(ServerModule, ServerOptionsModule)
)
