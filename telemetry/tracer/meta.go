package tracer

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/telemetry/attributes"
	sdk "go.opentelemetry.io/otel/sdk/trace"
)

// Meta extracts context metadata as OpenTelemetry span attributes.
//
// It reads the same exported metadata as the logger (via the meta package) and
// converts it to camel-cased string key/value attributes with no prefix, so a
// span carries the request/service context (request id, user id, ip, ...) used
// to correlate it with logs. It returns no attributes when ctx carries none.
func Meta(ctx context.Context) []attributes.KeyValue {
	return attributes.Strings(meta.CamelStrings(ctx, meta.NoPrefix))
}

// newMetaProcessor constructs a span processor that copies context metadata
// onto every span as it starts.
//
// It lets spans created by any instrumentation (server, database, cache, ...)
// carry the same request/service context as the logs (see [Meta]) without each
// instrumentation having to wire metadata itself. Metadata is read from the
// span's starting context, so the server/root span, whose context is not yet
// populated when it starts, is stamped explicitly by the transport metadata
// middleware instead.
func newMetaProcessor() sdk.SpanProcessor {
	return &metaProcessor{}
}

type metaProcessor struct{}

func (*metaProcessor) OnStart(parent context.Context, span sdk.ReadWriteSpan) {
	if attrs := Meta(parent); len(attrs) > 0 {
		span.SetAttributes(attrs...)
	}
}

func (*metaProcessor) OnEnd(sdk.ReadOnlySpan) {}

func (*metaProcessor) Shutdown(context.Context) error { return nil }

func (*metaProcessor) ForceFlush(context.Context) error { return nil }
