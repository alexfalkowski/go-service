package metrics

import (
	"github.com/alexfalkowski/go-service/pkg/metrics/prometheus/transport/http"
	"go.uber.org/fx"
)

// PrometheusModule for fx.
// nolint:gochecknoglobals
var PrometheusModule = fx.Invoke(http.Register)
