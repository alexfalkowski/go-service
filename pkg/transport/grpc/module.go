package grpc

import (
	"go.uber.org/fx"
)

var (
	// ServerModule for fx.
	ServerModule = fx.Provide(NewServer)

	// Module for fx.
	Module = fx.Options(ServerModule, fx.Provide(UnaryServerInterceptor), fx.Provide(StreamServerInterceptor))
)
