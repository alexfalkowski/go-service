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

// Module wires HTTP transport server, routing, and middleware into Fx.
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
