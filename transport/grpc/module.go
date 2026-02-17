package grpc

import (
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/transport/grpc/health"
	"github.com/alexfalkowski/go-service/v2/transport/grpc/token"
)

// Module wires gRPC transport server, auth, and middleware into Fx.
var Module = di.Module(
	di.Register(Register),
	di.Constructor(NewServerLimiter),
	di.Constructor(NewController),
	di.Constructor(NewToken),
	di.Constructor(token.NewGenerator),
	di.Constructor(token.NewVerifier),
	di.Constructor(NewServer),
	di.Constructor(registrar),
	health.Module,
)
