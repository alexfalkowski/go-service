package grpc

import (
	"go.uber.org/fx"
)

var (
	// ServerModule for fx.
	ServerModule = fx.Options(fx.Provide(NewServer), fx.Provide(NewConfig), fx.Provide(UnaryServerInterceptor), fx.Provide(StreamServerInterceptor))

	// Module for fx.
	Module = fx.Options(ServerModule)
)
