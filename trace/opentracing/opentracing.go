package opentracing

import (
	"context"
	"fmt"
	"strings"

	otr "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

// StartSpanFromContext for opentracing.
func StartSpanFromContext(ctx context.Context, tracer otr.Tracer, kind, operation, method string, opts ...otr.StartSpanOption) (context.Context, otr.Span) {
	fullMethod := fmt.Sprintf("%s %s", strings.ToLower(operation), strings.ToLower(method))
	clientSpan, ctx := otr.StartSpanFromContextWithTracer(ctx, tracer, fullMethod, opts...)

	ext.SpanKind.Set(clientSpan, ext.SpanKindEnum(kind))

	return otr.ContextWithSpan(ctx, clientSpan), clientSpan
}
