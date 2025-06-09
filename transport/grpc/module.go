package grpc

import (
	"github.com/alexfalkowski/go-service/v2/transport/grpc/health"
	"go.uber.org/fx"
)

// Module for fx.
var Module = fx.Options(
	fx.Invoke(Register),
	fx.Provide(NewServer),
	fx.Provide(provide),
	health.Module,
)
