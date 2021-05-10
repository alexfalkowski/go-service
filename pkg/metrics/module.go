package metrics

import (
	"github.com/alexfalkowski/go-service/pkg/metrics/prometheus/transport/http"
	"go.uber.org/fx"
)

// Module for fx.
var Module = fx.Invoke(http.Register)
