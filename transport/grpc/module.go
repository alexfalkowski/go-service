package grpc

import (
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/transport/grpc/health"
)

// Module for fx.
var Module = di.Module(
	di.Register(Register),
	di.Constructor(NewServerLimiter),
	di.Constructor(NewServer),
	di.Constructor(registrar),
	health.Module,
)
