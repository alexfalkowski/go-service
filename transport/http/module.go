package http

import (
	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/di"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/content"
	"github.com/alexfalkowski/go-service/v2/net/http/mvc"
	"github.com/alexfalkowski/go-service/v2/net/http/rest"
	"github.com/alexfalkowski/go-service/v2/net/http/rpc"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/alexfalkowski/go-service/v2/transport/http/health"
	"github.com/alexfalkowski/go-service/v2/transport/http/telemetry/metrics"
	"github.com/alexfalkowski/go-service/v2/transport/http/token"
	sync "github.com/alexfalkowski/go-sync"
)

// Module wires the HTTP transport stack into [go.uber.org/fx].
//
// It composes constructors and registrations required to run an HTTP server and to support common
// handler styles used by go-service:
//   - mux, route policy, and router construction ([http.NewServeMux], [http.NewRoutePolicy], [http.NewRouter])
//   - content negotiation and encoding ([content.NewContent])
//   - MVC view rendering helpers ([mvc.NewFunctionMap], [mvc.Register])
//   - RPC and REST routing ([rpc.Register], [rest.Register] — called from [registerRoutes] rather than
//     as their own separate Fx invoke targets, so their timeout/maxReceiveSize arguments are plain
//     values this package computes from its own *[Config], not separate Fx-resolved dependencies),
//     including the streaming route helpers (rest.StreamRoute/StreamGet/StreamRouteRequest/StreamPost/
//     StreamPut/StreamPatch and rpc.StreamRoute): a successful Send extends the response write deadline,
//     while a successful Recv on a bidirectional stream extends the request read deadline (see
//     [github.com/alexfalkowski/go-service/v2/net/http/content.NewRequestStreamHandler]); bidirectional
//     streaming routes also bound each decoded request value instead of the request body's cumulative size
//     (see [github.com/alexfalkowski/go-service/v2/net/http/content.NewRequestStreamHandler])
//   - transport-level middleware wiring (limiter and token helpers)
//   - server construction ([NewServer])
//   - operational endpoints (Prometheus metrics and health)
//
// This module also registers [Register], which injects the filesystem dependency used by this package
// (required when constructing TLS configuration from source strings).
var Module = di.Module(
	di.Register(Register),
	di.Constructor(http.NewServeMux),
	di.Constructor(http.NewRoutePolicy),
	di.Constructor(http.NewRouter),
	di.Constructor(content.NewContent),
	di.Constructor(mvc.NewFunctionMap),
	di.Register(mvc.Register),
	di.Register(registerRoutes),
	di.Constructor(NewServerLimiter),
	di.Constructor(NewToken),
	di.Constructor(token.NewGenerator),
	di.Constructor(token.NewVerifier),
	di.Constructor(NewServer),
	di.Register(metrics.Register),
	health.Module,
)

// registerRoutes wires [rest.Register] and [rpc.Register] with the router, content, and buffer pool
// dependencies Fx resolves normally, plus a timeout/maxReceiveSize pair computed here from cfg and
// passed as plain arguments.
//
// Calling rest.Register/rpc.Register directly (rather than as their own separate Fx invoke targets)
// means timeout ([time.Duration]) and maxReceiveSize ([bytes.Size]) never need to exist as their own DI
// graph nodes: those are common, easily-collided types (nothing else in the DI graph keys on them, but
// nothing stops something else from needing to someday), so resolving them by bare type would be
// fragile in a way resolving *[Config] is not. net/http/rest and net/http/rpc must not import this
// package (see AGENTS.md), so cfg's values are computed here and passed in, mirroring how NewServer
// already derives its own ReadTimeout/WriteTimeout/max receive size from the same *Config.
func registerRoutes(cfg *Config, router *http.Router, cont *content.Content, pool *sync.BufferPool) {
	var timeout time.Duration
	var maxReceiveSize bytes.Size

	if cfg.IsEnabled() {
		timeout = cfg.GetTimeout()
		maxReceiveSize = cfg.GetMaxReceiveSize()
	}

	rest.Register(router, cont, pool, timeout, maxReceiveSize)
	rpc.Register(router, cont, pool, timeout, maxReceiveSize)
}
