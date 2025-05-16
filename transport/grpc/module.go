package grpc

import "go.uber.org/fx"

// Module for fx.
var Module = fx.Options(
	fx.Invoke(Register),
	fx.Provide(NewServer),
	fx.Provide(provide),
)
