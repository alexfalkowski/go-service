package http

import (
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/content"
	"github.com/alexfalkowski/go-service/v2/net/http/mvc"
	"github.com/alexfalkowski/go-service/v2/net/http/rest"
	"github.com/alexfalkowski/go-service/v2/net/http/rpc"
	"github.com/alexfalkowski/go-service/v2/transport/http/health"
	"github.com/alexfalkowski/go-service/v2/transport/http/telemetry/metrics"
	"github.com/alexfalkowski/go-service/v2/transport/http/token"
)

// Module wires the HTTP transport stack into Fx.
//
// It composes constructors and registrations required to run an HTTP server and to support common
// handler styles used by go-service:
//   - mux construction (`http.NewServeMux`)
//   - content negotiation and encoding (`content.NewContent`)
//   - MVC view rendering helpers (`mvc.NewFunctionMap`, `mvc.Register`)
//   - RPC and REST routing (`rpc.Register`, `rest.Register`)
//   - transport-level middleware wiring (limiter and token helpers)
//   - server construction (`NewServer`)
//   - operational endpoints (Prometheus metrics and health)
//
// This module also registers `Register`, which injects the filesystem dependency used by this package
// (required when constructing TLS configuration from certificate/key source strings).
var Module = di.Module(
	di.Register(Register),
	di.Constructor(http.NewServeMux),
	di.Constructor(content.NewContent),
	di.Constructor(mvc.NewFunctionMap),
	di.Register(mvc.Register),
	di.Register(rpc.Register),
	di.Register(rest.Register),
	di.Constructor(NewServerLimiter),
	di.Constructor(NewController),
	di.Constructor(NewToken),
	di.Constructor(token.NewGenerator),
	di.Constructor(token.NewVerifier),
	di.Constructor(NewServer),
	di.Register(metrics.Register),
	health.Module,
)
