package http

import (
	hm "github.com/alexfalkowski/go-service/net/http/mux"
	"go.uber.org/fx"
)

// Module for fx.
var Module = fx.Options(
	fx.Provide(Mux),
	fx.Provide(hm.NewRuntimeServeMux),
	fx.Provide(hm.NewStandardServeMux),
	fx.Provide(NewServer),
	fx.Invoke(RegisterMetrics),
)
