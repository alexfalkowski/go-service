package grpc

import (
	"github.com/alexfalkowski/go-service/transport/grpc/security/token"
	"go.uber.org/fx"
)

// Module for fx.
var Module = fx.Options(
	token.Module,
	fx.Provide(NewServer),
)
