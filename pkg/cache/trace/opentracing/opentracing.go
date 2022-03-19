package opentracing

import (
	"context"

	potr "github.com/alexfalkowski/go-service/pkg/trace/opentracing"
	otr "github.com/opentracing/opentracing-go"
)

// StartSpanFromContext for cache.
// nolint:ireturn
func StartSpanFromContext(ctx context.Context, operation, method string, opts ...otr.StartSpanOption) (context.Context, otr.Span) {
	return potr.StartSpanFromContext(ctx, "cache", operation, method, opts...)
}
