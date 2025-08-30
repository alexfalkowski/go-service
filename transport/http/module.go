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
)

// Module for fx.
var Module = di.Module(
	di.Register(Register),
	di.Constructor(http.NewServeMux),
	di.Constructor(content.NewContent),
	di.Constructor(mvc.NewFunctionMap),
	di.Register(mvc.Register),
	di.Register(rpc.Register),
	di.Register(rest.Register),
	di.Constructor(NewServerLimiter),
	di.Constructor(NewServer),
	di.Register(metrics.Register),
	health.Module,
)
