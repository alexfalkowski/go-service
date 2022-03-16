package opentracing

import (
	"context"
	"fmt"

	"github.com/alexfalkowski/go-service/pkg/meta"
	"github.com/alexfalkowski/go-service/pkg/time"
	"github.com/alexfalkowski/go-service/pkg/transport/nsq/handler"
	"github.com/alexfalkowski/go-service/pkg/transport/nsq/message"
	"github.com/alexfalkowski/go-service/pkg/transport/nsq/producer"
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
	nsqTopic     = "nsq.topic"
	nsqChannel   = "nsq.channel"
	component    = "component"
	nsqComponent = "nsq"
)

// NewHandler for opentracing.
func NewHandler(topic, channel string, h handler.Handler) *Handler {
	return &Handler{topic: topic, channel: channel, Handler: h}
}

// Handler for opentracing.
type Handler struct {
	topic, channel string

	handler.Handler
}

func (h *Handler) Handle(ctx context.Context, message *message.Message) (context.Context, error) {
	start := time.Now().UTC()
	tracer := opentracing.GlobalTracer()
	traceCtx, _ := tracer.Extract(opentracing.TextMap, headersTextMap(message.Headers))
	operationName := fmt.Sprintf("consume %s:%s", h.topic, h.channel)
	opts := []opentracing.StartSpanOption{
		ext.RPCServerOption(traceCtx),
		opentracing.Tag{Key: nsqStartTime, Value: start.Format(time.RFC3339)},
		opentracing.Tag{Key: nsqTopic, Value: h.topic},
		opentracing.Tag{Key: nsqChannel, Value: h.channel},
		opentracing.Tag{Key: nsqID, Value: string(message.ID[:])},
		opentracing.Tag{Key: nsqBody, Value: string(message.Body)},
		opentracing.Tag{Key: nsqTimestamp, Value: message.Timestamp},
		opentracing.Tag{Key: nsqAttempts, Value: message.Attempts},
		opentracing.Tag{Key: nsqAddress, Value: message.NSQDAddress},
		opentracing.Tag{Key: component, Value: nsqComponent},
		ext.SpanKindConsumer,
	}

	span := tracer.StartSpan(operationName, opts...)
	defer span.Finish()

	ctx = opentracing.ContextWithSpan(ctx, span)
	ctx, err := h.Handler.Handle(ctx, message)

	for k, v := range meta.Attributes(ctx) {
		span.SetTag(k, v)
	}

	span.SetTag(nsqDuration, time.ToMilliseconds(time.Since(start)))

	if err != nil {
		setError(span, err)
	}

	return ctx, err
}

// NewProducer for opentracing.
func NewProducer(p producer.Producer) *Producer {
	return &Producer{Producer: p}
}

// Producer for opentracing.
type Producer struct {
	producer.Producer
}

func (p *Producer) Publish(ctx context.Context, topic string, message *message.Message) (context.Context, error) {
	start := time.Now().UTC()
	tracer := opentracing.GlobalTracer()
	operationName := fmt.Sprintf("publish %s", topic)
	opts := []opentracing.StartSpanOption{
		opentracing.Tag{Key: nsqStartTime, Value: start.Format(time.RFC3339)},
		opentracing.Tag{Key: nsqBody, Value: string(message.Body)},
		opentracing.Tag{Key: nsqTopic, Value: topic},
		opentracing.Tag{Key: component, Value: nsqComponent},
		ext.SpanKindProducer,
	}

	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, tracer, operationName, opts...)
	defer span.Finish()

	if err := tracer.Inject(span.Context(), opentracing.TextMap, headersTextMap(message.Headers)); err != nil {
		return ctx, err
	}

	ctx, err := p.Producer.Publish(ctx, topic, message)

	for k, v := range meta.Attributes(ctx) {
		span.SetTag(k, v)
	}

	span.SetTag(nsqDuration, time.ToMilliseconds(time.Since(start)))

	if err != nil {
		setError(span, err)
	}

	return ctx, err
}

func setError(span opentracing.Span, err error) {
	ext.Error.Set(span, true)
	span.LogFields(log.String("event", "error"), log.String("message", err.Error()))
}
