package opentracing

import (
	"context"
	"fmt"

	pkgMeta "github.com/alexfalkowski/go-service/pkg/meta"
	"github.com/alexfalkowski/go-service/pkg/time"
	"github.com/alexfalkowski/go-service/pkg/transport/nsq/handler"
	"github.com/alexfalkowski/go-service/pkg/transport/nsq/meta"
	"github.com/nsqio/go-nsq"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
)

const (
	nsqID        = "nsq.id"
	nsqBody      = "nsq.body"
	nsqTimestamp = "nsq.timestamp"
	nsqAttempts  = "nsq.attempts"
	nsqAddress   = "nsq.address"
	nsqDuration  = "nsq.duration_ms"
	nsqStartTime = "nsq.start_time"
	component    = "component"
	nsqComponent = "nsq"
)

// NewHandler for zap.
func NewHandler(h handler.Handler) handler.Handler {
	return &traceHandler{Handler: h}
}

type traceHandler struct {
	handler.Handler
}

func (h *traceHandler) Handle(ctx context.Context, message *nsq.Message) (context.Context, error) {
	start := time.Now().UTC()
	tracer := opentracing.GlobalTracer()
	operationName := fmt.Sprintf("Consume msg %s(%s)", pkgMeta.Attribute(ctx, meta.Topic), pkgMeta.Attribute(ctx, meta.Channel))
	opts := []opentracing.StartSpanOption{
		opentracing.Tag{Key: nsqStartTime, Value: start.Format(time.RFC3339)},
		opentracing.Tag{Key: nsqID, Value: message.ID[:]},
		opentracing.Tag{Key: nsqBody, Value: message.Body},
		opentracing.Tag{Key: nsqTimestamp, Value: message.Timestamp},
		opentracing.Tag{Key: nsqAttempts, Value: message.Attempts},
		opentracing.Tag{Key: nsqAddress, Value: message.NSQDAddress},
		opentracing.Tag{Key: component, Value: nsqComponent},
		ext.SpanKindRPCClient,
	}

	clientSpan, ctx := opentracing.StartSpanFromContextWithTracer(ctx, tracer, operationName, opts...)

	defer clientSpan.Finish()

	// TODO: inject from headers.
	// carrier := opentracing.HTTPHeadersCarrier(req.Header)
	// if err := tracer.Inject(clientSpan.Context(), opentracing.HTTPHeaders, carrier); err != nil {
	// 	return nil, err
	// }

	ctx, err := h.Handler.Handle(ctx, message)

	for k, v := range pkgMeta.Attributes(ctx) {
		clientSpan.SetTag(k, v)
	}

	clientSpan.SetTag(nsqDuration, time.ToMilliseconds(time.Since(start)))

	if err != nil {
		ext.Error.Set(clientSpan, true)
		clientSpan.LogFields(log.String("event", "error"), log.String("message", err.Error()))
	}

	return ctx, err
}
