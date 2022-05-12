package metrics

import (
	"github.com/alexfalkowski/go-service/metrics/prometheus/transport/http"
	"go.uber.org/fx"
)

// PrometheusModule for fx.
var PrometheusModule = fx.Invoke(http.Register)
