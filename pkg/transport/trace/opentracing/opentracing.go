package opentracing

import (
	"context"

	"github.com/alexfalkowski/go-service/pkg/transport/meta"
	"github.com/opentracing/opentracing-go"
)

// StartSpanOptions for opentracing.
func StartSpanOptions(ctx context.Context) []opentracing.StartSpanOption {
	return []opentracing.StartSpanOption{
		opentracing.Tag{Key: meta.RequestIDKey, Value: meta.RequestID(ctx)},
		opentracing.Tag{Key: meta.UserAgentKey, Value: meta.UserAgent(ctx)},
	}
}
