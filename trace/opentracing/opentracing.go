package opentracing

import (
	"context"
	"fmt"
	"strings"

	"github.com/alexfalkowski/go-service/trace/opentracing/datadog"
	"github.com/alexfalkowski/go-service/trace/opentracing/jaeger"
	otr "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"go.uber.org/fx"
)

// StartSpanFromContext for opentracing.
func StartSpanFromContext(ctx context.Context, tracer otr.Tracer, kind, operation, method string, opts ...otr.StartSpanOption) (context.Context, otr.Span) {
	fullMethod := fmt.Sprintf("%s %s", strings.ToLower(operation), strings.ToLower(method))
	clientSpan, ctx := otr.StartSpanFromContextWithTracer(ctx, tracer, fullMethod, opts...)

	ext.SpanKind.Set(clientSpan, ext.SpanKindEnum(kind))

	return otr.ContextWithSpan(ctx, clientSpan), clientSpan
}

// TracerParams for opentracing.
type TracerParams struct {
	Lifecycle fx.Lifecycle
	Name      string
	Config    *Config
}

// NewTracer for opentracing.
func NewTracer(params TracerParams) (otr.Tracer, error) {
	if params.Config.IsJaeger() {
		return jaeger.NewTracer(jaeger.TracerParams{Lifecycle: params.Lifecycle, Name: params.Name, Host: params.Config.Host})
	}

	if params.Config.IsDataDog() {
		return datadog.NewTracer(datadog.TracerParams{Lifecycle: params.Lifecycle, Name: params.Name, Host: params.Config.Host}), nil
	}

	return nil, nil
}
