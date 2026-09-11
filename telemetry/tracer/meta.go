package tracer

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/alexfalkowski/go-service/v2/telemetry/attributes"
	sdk "go.opentelemetry.io/otel/sdk/trace"
)

// Meta extracts context metadata as OpenTelemetry span attributes.
//
// It reads the same exported metadata as the logger (via the meta package) as
// string key/value attributes with no prefix, so a span carries the
// request/service context (request_id, user_id, ip_addr, ...) used to correlate
// it with logs. Standard metadata keys use snake_case; arbitrary keys are
// preserved. It returns no attributes when ctx carries none.
func Meta(ctx context.Context, limit meta.Limit) []attributes.KeyValue {
	return attributes.Strings(meta.Attributes(ctx, limit))
}

type metaProcessor struct {
	limit meta.Limit
}

func (p *metaProcessor) OnStart(parent context.Context, span sdk.ReadWriteSpan) {
	if attrs := Meta(parent, p.limit); len(attrs) > 0 {
		span.SetAttributes(attrs...)
	}
}

func (*metaProcessor) OnEnd(sdk.ReadOnlySpan) {}

func (*metaProcessor) Shutdown(context.Context) error { return nil }

func (*metaProcessor) ForceFlush(context.Context) error { return nil }
