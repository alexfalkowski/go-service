package grpc

import (
	"github.com/alexfalkowski/go-service/v3/net/grpc"
	"github.com/alexfalkowski/go-service/v3/transport/grpc/health"
	"github.com/alexfalkowski/go-service/v3/transport/grpc/token"
	"github.com/alexfalkowski/go-signal/v2"
)

// Module wires the gRPC transport stack into [github.com/alexfalkowski/go-signal/v2].
//
// It composes constructors and registrations required to run a gRPC server and to support common
// cross-cutting concerns used by go-service:
//
//   - transport registration for TLS filesystem access ([Register])
//   - server-side rate limiting wiring ([NewServerLimiter])
//   - token service wiring ([NewToken])
//   - token generator/verifier adapters for interceptor wiring
//     ([github.com/alexfalkowski/go-service/v3/transport/grpc/token.NewGenerator],
//     [github.com/alexfalkowski/go-service/v3/transport/grpc/token.NewVerifier])
//   - server construction ([NewServer]) and service registration (`registrar`)
//   - health service wiring ([github.com/alexfalkowski/go-service/v3/transport/grpc/health.Module])
//
// This module also registers [Register], which injects the filesystem dependency used by this package
// (required when constructing TLS configuration from source strings).
var Module = signal.Module(
	signal.Register(Register),
	signal.Constructor(NewServerLimiter),
	signal.Constructor(NewToken),
	signal.Constructor(grpc.NewMethodPolicy),
	signal.Constructor(token.NewGenerator),
	signal.Constructor(token.NewVerifier),
	signal.Constructor(NewServer),
	signal.Constructor(registrar),
	signal.Constructor(serviceRegistrar),
	health.Module,
)
