package http

import (
	"github.com/alexfalkowski/go-service/v3/net/http"
	"github.com/alexfalkowski/go-service/v3/net/http/content/stream"
	"github.com/alexfalkowski/go-service/v3/net/http/content/unary"
	"github.com/alexfalkowski/go-service/v3/net/http/mvc"
	"github.com/alexfalkowski/go-service/v3/net/http/rest"
	"github.com/alexfalkowski/go-service/v3/net/http/rpc"
	"github.com/alexfalkowski/go-service/v3/net/server"
	"github.com/alexfalkowski/go-service/v3/transport/http/health"
	"github.com/alexfalkowski/go-service/v3/transport/http/telemetry/metrics"
	"github.com/alexfalkowski/go-service/v3/transport/http/token"
	"github.com/alexfalkowski/go-signal/v2"
)

// Module wires the HTTP transport stack into [github.com/alexfalkowski/go-signal/v2].
//
// It composes constructors and registrations required to run an HTTP server and to support common
// handler styles used by go-service:
//   - mux, route policy, and router construction ([http.NewServeMux], [http.NewRoutePolicy], [http.NewRouter])
//   - content negotiation and encoding ([unary.NewContent], [contentstream.NewContent])
//   - MVC view rendering helpers ([mvc.NewFunctionMap], [mvc.NewServer])
//   - RPC and REST routing ([rpc.NewServer], [rest.NewServer], both resolving [stream.Options] as an
//     go-signal dependency computed from *[Config] by [newStreamOptions]), including the streaming route
//     helpers (rest.StreamRoute/StreamGet/StreamRouteRequest/StreamPost/StreamPut/StreamPatch and
//     rpc.StreamRoute): on a bidirectional stream, a successful Send or Recv extends both the response
//     write deadline and the request read deadline (see
//     [github.com/alexfalkowski/go-service/v3/net/http/content/stream.NewRequestHandler]); a send-only
//     stream extends only the write deadline. Bidirectional streaming routes also bound each decoded
//     request value instead of the request body's cumulative size (see
//     [github.com/alexfalkowski/go-service/v3/net/http/content/stream.NewRequestHandler]). At shutdown,
//     the shared server drain signal cancels stream handler contexts; handlers must return after
//     observing that cancellation when waiting outside Send or Recv.
//   - transport-level middleware wiring (limiter and token helpers)
//   - server construction ([NewServer])
//   - operational endpoints (Prometheus metrics and health)
//
// This module also registers [Register], which injects the filesystem dependency used by this package
// (required when constructing TLS configuration from source strings).
var Module = signal.Module(
	signal.Register(Register),
	signal.Constructor(http.NewServeMux),
	signal.Constructor(http.NewRoutePolicy),
	signal.Constructor(http.NewRouter),
	signal.Constructor(unary.NewContent),
	signal.Constructor(stream.NewContent),
	signal.Constructor(mvc.NewFunctionMap),
	signal.Constructor(mvc.NewServer),
	signal.Constructor(newStreamOptions),
	signal.Constructor(rest.NewServer),
	signal.Constructor(rpc.NewServer),
	signal.Constructor(NewServerLimiter),
	signal.Constructor(NewToken),
	signal.Constructor(token.NewGenerator),
	signal.Constructor(token.NewVerifier),
	signal.Constructor(NewServer),
	signal.Register(metrics.Register),
	health.Module,
)

// newStreamOptions resolves the [stream.Options] dependency for [rest.NewServer] and [rpc.NewServer].
//
// net/http/rest and net/http/rpc must not import this package (see AGENTS.md), so cfg's values are
// computed here and resolved via go-signal, mirroring how NewServer already derives its own
// ReadTimeout/WriteTimeout/max receive size from the same *Config: the read/write timeouts resolve
// through the same options-aware precedence (options key, falling back to the lower-level default)
// NewServer uses for its own ReadTimeout/WriteTimeout, so a service that sets explicit streaming
// budgets gets a matching server deadline instead of a silently different one.
func newStreamOptions(cfg *Config, drain *server.Drain) stream.Options {
	var so stream.Options

	if cfg.IsEnabled() {
		so = stream.Options{
			ReadTimeout:    cfg.GetReadTimeout(),
			WriteTimeout:   cfg.GetWriteTimeout(),
			MaxReceiveSize: cfg.GetMaxReceiveSize(),
			Drain:          drain.Done(),
		}
	}

	return so
}
