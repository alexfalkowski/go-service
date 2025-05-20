package http

import (
	"net/http"

	"github.com/alexfalkowski/go-service/v2/net/http/content"
	"github.com/alexfalkowski/go-service/v2/net/http/mvc"
	"github.com/alexfalkowski/go-service/v2/net/http/rest"
	"github.com/alexfalkowski/go-service/v2/net/http/rpc"
	"github.com/alexfalkowski/go-service/v2/transport/http/telemetry/metrics"
	"go.uber.org/fx"
)

// Module for fx.
var Module = fx.Options(
	fx.Invoke(Register),
	fx.Provide(http.NewServeMux),
	fx.Provide(content.NewContent),
	fx.Provide(mvc.NewFunctionMap),
	fx.Invoke(mvc.Register),
	fx.Invoke(rpc.Register),
	fx.Invoke(rest.Register),
	fx.Provide(NewServer),
	fx.Invoke(metrics.Register),
)
