package http

import (
	"github.com/alexfalkowski/go-service/net/http"
	"go.uber.org/fx"
)

// Module for fx.
var Module = fx.Options(
	fx.Provide(Mux),
	fx.Provide(http.NewRuntimeServeMux),
	fx.Provide(http.NewStandardServeMux),
	fx.Provide(http.NewServeMux),
	fx.Provide(NewServer),
	fx.Invoke(RegisterMetrics),
)
