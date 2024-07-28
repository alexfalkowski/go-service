package http

import (
	"net/http"

	"github.com/alexfalkowski/go-service/net/http/mvc"
	"github.com/alexfalkowski/go-service/net/http/rpc"
	"go.uber.org/fx"
)

// Module for fx.
var Module = fx.Options(
	fx.Provide(http.NewServeMux),
	fx.Provide(mvc.NewViews),
	fx.Provide(mvc.NewRouter),
	fx.Provide(rpc.NewRouter),
	fx.Provide(NewServer),
	fx.Invoke(RegisterMetrics),
)
