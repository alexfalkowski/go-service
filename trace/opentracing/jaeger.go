package opentracing

import (
	"github.com/alexfalkowski/go-service/os"
	"github.com/alexfalkowski/go-service/trace/opentracing/jaeger"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// NewJaegerServiceTracer for opentracing.
// nolint:ireturn
func NewJaegerServiceTracer(lc fx.Lifecycle, logger *zap.Logger, cfg *jaeger.Config) (ServiceTracer, error) {
	return jaeger.NewTracer(lc, os.ExecutableName(), logger, cfg)
}

// NewJaegerDatabaseTracer for opentracing.
// nolint:ireturn
func NewJaegerDatabaseTracer(lc fx.Lifecycle, logger *zap.Logger, cfg *jaeger.Config) (DatabaseTracer, error) {
	return jaeger.NewTracer(lc, database, logger, cfg)
}

// NewJaegerCacheTracer for opentracing.
// nolint:ireturn
func NewJaegerCacheTracer(lc fx.Lifecycle, logger *zap.Logger, cfg *jaeger.Config) (CacheTracer, error) {
	return jaeger.NewTracer(lc, cache, logger, cfg)
}

// NewJaegerTransportTracer for opentracing.
// nolint:ireturn
func NewJaegerTransportTracer(lc fx.Lifecycle, logger *zap.Logger, cfg *jaeger.Config) (TransportTracer, error) {
	return jaeger.NewTracer(lc, transport, logger, cfg)
}
