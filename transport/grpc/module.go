package grpc

import (
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/transport/grpc/health"
	"github.com/alexfalkowski/go-service/v2/transport/grpc/token"
)

// Module wires the gRPC transport stack into [go.uber.org/fx].
//
// It composes constructors and registrations required to run a gRPC server and to support common
// cross-cutting concerns used by go-service:
//
//   - transport registration for TLS filesystem access ([Register])
//   - server-side rate limiting wiring ([NewServerLimiter])
//   - token access-controller construction and token service wiring ([NewController], [NewToken])
//   - token generator/verifier adapters for interceptor wiring
//     ([github.com/alexfalkowski/go-service/v2/transport/grpc/token.NewGenerator],
//     [github.com/alexfalkowski/go-service/v2/transport/grpc/token.NewVerifier])
//   - server construction ([NewServer]) and service registration (`registrar`)
//   - health service wiring ([github.com/alexfalkowski/go-service/v2/transport/grpc/health.Module])
//
// This module also registers [Register], which injects the filesystem dependency used by this package
// (required when constructing TLS configuration from source strings).
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
