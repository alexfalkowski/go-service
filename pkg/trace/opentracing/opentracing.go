package opentracing

import (
	"context"
	"fmt"

	otr "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

// StartSpanFromContext for opentracing.
// nolint:ireturn
func StartSpanFromContext(ctx context.Context, kind, operation, method string, opts ...otr.StartSpanOption) (context.Context, otr.Span) {
	fullMethod := fmt.Sprintf("%s.%s", operation, method)
	clientSpan, ctx := otr.StartSpanFromContextWithTracer(ctx, otr.GlobalTracer(), fullMethod, opts...)

	ext.SpanKind.Set(clientSpan, ext.SpanKindEnum(kind))

	return otr.ContextWithSpan(ctx, clientSpan), clientSpan
}
