package metrics

import (
	"github.com/alexfalkowski/go-service/pkg/metrics/prometheus/transport/http"
	"go.uber.org/fx"
)

var (
	// PrometheusModule for fx.
	PrometheusModule = fx.Invoke(http.Register)
)
