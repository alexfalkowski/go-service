package http

import (
	"net/http"

	"github.com/alexfalkowski/go-service/net/http/content"
	"github.com/alexfalkowski/go-service/net/http/mvc"
	"github.com/alexfalkowski/go-service/net/http/rest"
	"github.com/alexfalkowski/go-service/net/http/rpc"
	"go.uber.org/fx"
)

// Module for fx.
var Module = fx.Options(
	fx.Provide(http.NewServeMux),
	fx.Invoke(content.Register),
	fx.Provide(content.NewContent),
	fx.Provide(mvc.NewViews),
	fx.Provide(mvc.NewRouter),
	fx.Invoke(rpc.Register),
	fx.Invoke(rest.Register),
	fx.Provide(NewServer),
	fx.Invoke(RegisterMetrics),
)
