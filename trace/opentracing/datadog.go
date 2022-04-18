package opentracing

import (
	"github.com/alexfalkowski/go-service/os"
	"github.com/alexfalkowski/go-service/trace/opentracing/datadog"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// NewDataDogDatabaseTracer for opentracing.
// nolint:ireturn
func NewDataDogDatabaseTracer(lc fx.Lifecycle, logger *zap.Logger, cfg *datadog.Config) DatabaseTracer {
	return datadog.NewTracer(lc, os.ExecutableName(), logger, cfg)
}

// NewDataDogCacheTracer for opentracing.
// nolint:ireturn
func NewDataDogCacheTracer(lc fx.Lifecycle, logger *zap.Logger, cfg *datadog.Config) CacheTracer {
	return datadog.NewTracer(lc, os.ExecutableName(), logger, cfg)
}

// NewDataDogTransportTracer for opentracing.
// nolint:ireturn
func NewDataDogTransportTracer(lc fx.Lifecycle, logger *zap.Logger, cfg *datadog.Config) TransportTracer {
	return datadog.NewTracer(lc, os.ExecutableName(), logger, cfg)
}
