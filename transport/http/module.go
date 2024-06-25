package http

import (
	"github.com/alexfalkowski/go-service/net/http"
	"go.uber.org/fx"
)

// Module for fx.
var Module = fx.Options(
	fx.Provide(http.NewServeMux),
	fx.Invoke(http.RegisterHandler),
	fx.Provide(NewServer),
	fx.Invoke(RegisterMetrics),
)
