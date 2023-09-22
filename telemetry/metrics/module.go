package metrics

import (
	"github.com/alexfalkowski/go-service/telemetry/metrics/prometheus/transport/http"
	"go.uber.org/fx"
)

// Module for fx.
var Module = fx.Options(
	fx.Invoke(http.Register),
)
