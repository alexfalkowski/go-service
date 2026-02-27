package transport

import (
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/transport/grpc"
	"github.com/alexfalkowski/go-service/v2/transport/http"
	"github.com/alexfalkowski/go-service/v2/transport/http/events"
)

// Module wires the transport layer into Fx.
//
// It composes the transport submodules used by go-service and provides top-level lifecycle wiring for
// starting and stopping enabled server services.
//
// Included submodules:
//   - `transport/grpc.Module`: gRPC transport server wiring (including TLS filesystem registration, interceptors, and health).
//   - `transport/http.Module`: HTTP transport server wiring (routing helpers, middleware, TLS filesystem registration, metrics, and health).
//   - `transport/http/events.Module`: CloudEvents HTTP wiring (sender/receiver helpers and webhook signing/verification adapters).
//
// Server lifecycle wiring:
// This module also provides `NewServers` (to collect enabled `*server.Service` instances) and registers
// `Register`, which attaches lifecycle hooks to start and stop those services.
var Module = di.Module(
	grpc.Module,
	http.Module,
	events.Module,
	di.Constructor(NewServers),
	di.Register(Register),
)
