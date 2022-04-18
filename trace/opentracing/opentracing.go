package opentracing

import (
	"context"
	"fmt"
	"strings"

	"github.com/alexfalkowski/go-service/os"
	otr "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

const (
	database  = "database"
	cache     = "cache"
	transport = "transport"
)

// ServiceTracer for opentracing.
type ServiceTracer otr.Tracer

// DatabaseTracer for opentracing.
type DatabaseTracer otr.Tracer

// CacheTracer for opentracing.
type CacheTracer otr.Tracer

// TransportTracer for opentracing.
type TransportTracer otr.Tracer

// StartServiceSpanFromContext for opentracing.
// nolint:ireturn
func StartServiceSpanFromContext(ctx context.Context, tracer CacheTracer, operation, method string, opts ...otr.StartSpanOption) (context.Context, otr.Span) {
	return StartSpanFromContext(ctx, tracer, os.ExecutableName(), operation, method, opts...)
}

// StartDatabaseSpanFromContext for opentracing.
// nolint:ireturn
func StartDatabaseSpanFromContext(ctx context.Context, tracer CacheTracer, operation, method string, opts ...otr.StartSpanOption) (context.Context, otr.Span) {
	return StartSpanFromContext(ctx, tracer, database, operation, method, opts...)
}

// StartCacheSpanFromContext for opentracing.
// nolint:ireturn
func StartCacheSpanFromContext(ctx context.Context, tracer CacheTracer, operation, method string, opts ...otr.StartSpanOption) (context.Context, otr.Span) {
	return StartSpanFromContext(ctx, tracer, cache, operation, method, opts...)
}

// StartTransportSpanFromContext for opentracing.
// nolint:ireturn
func StartTransportSpanFromContext(ctx context.Context, tracer CacheTracer, operation, method string, opts ...otr.StartSpanOption) (context.Context, otr.Span) {
	return StartSpanFromContext(ctx, tracer, transport, operation, method, opts...)
}

// StartSpanFromContext for opentracing.
// nolint:ireturn
func StartSpanFromContext(ctx context.Context, tracer otr.Tracer, kind, operation, method string, opts ...otr.StartSpanOption) (context.Context, otr.Span) {
	fullMethod := fmt.Sprintf("%s %s", strings.ToLower(operation), strings.ToLower(method))
	clientSpan, ctx := otr.StartSpanFromContextWithTracer(ctx, tracer, fullMethod, opts...)

	ext.SpanKind.Set(clientSpan, ext.SpanKindEnum(kind))

	return otr.ContextWithSpan(ctx, clientSpan), clientSpan
}
