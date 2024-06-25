package http

import (
	"net/http"

	nh "github.com/alexfalkowski/go-service/net/http"
	"go.uber.org/fx"
)

// Module for fx.
var Module = fx.Options(
	fx.Provide(http.NewServeMux),
	fx.Invoke(nh.Register),
	fx.Provide(NewServer),
	fx.Invoke(RegisterMetrics),
)
