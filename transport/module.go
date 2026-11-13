package transport

import (
	"github.com/alexfalkowski/go-service/v3/net/server"
	"github.com/alexfalkowski/go-service/v3/transport/grpc"
	"github.com/alexfalkowski/go-service/v3/transport/http"
	"github.com/alexfalkowski/go-service/v3/transport/http/events"
	"github.com/alexfalkowski/go-signal/v2"
)

// Module wires the transport layer into [github.com/alexfalkowski/go-signal/v2].
//
// It composes the transport submodules used by go-service and provides top-level lifecycle wiring for
// starting and stopping enabled server services.
//
// Included submodules:
//   - [github.com/alexfalkowski/go-service/v3/transport/grpc.Module]: gRPC transport server wiring (including TLS filesystem registration, interceptors, and health).
//   - [github.com/alexfalkowski/go-service/v3/transport/http.Module]: HTTP transport server wiring (routing helpers, middleware, TLS filesystem registration, metrics, and health).
//   - [github.com/alexfalkowski/go-service/v3/transport/http/events.Module]: CloudEvents HTTP receiver and webhook verification wiring.
//
// Server lifecycle wiring:
// This module also provides [NewServers] (to collect enabled *[server.Service] instances) and registers
// [github.com/alexfalkowski/go-service/v3/net/server.Register], which attaches lifecycle hooks to start, drain, and stop those services.
//
// In typical service applications this module is consumed through [github.com/alexfalkowski/go-service/v3/module.Server], so callers normally do
// not need to invoke lower-level transport registration or lifecycle helpers themselves.
var Module = signal.Module(
	signal.Constructor(server.NewDrain),
	grpc.Module,
	http.Module,
	events.Module,
	signal.Constructor(NewAccessController),
	signal.Constructor(NewServers),
	signal.Register(server.Register),
)
