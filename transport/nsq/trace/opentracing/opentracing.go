package opentracing

import (
	"context"
	"fmt"
	"time"

	"github.com/alexfalkowski/go-service/meta"
	stime "github.com/alexfalkowski/go-service/time"
	"github.com/alexfalkowski/go-service/transport/nsq/handler"
	"github.com/alexfalkowski/go-service/transport/nsq/message"
	"github.com/alexfalkowski/go-service/transport/nsq/producer"
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
func NewHandler(topic, channel string, tracer opentracing.Tracer, h handler.Handler) *Handler {
	return &Handler{topic: topic, channel: channel, tracer: tracer, Handler: h}
}

// Handler for opentracing.
type Handler struct {
	topic, channel string
	tracer         opentracing.Tracer

	handler.Handler
}

func (h *Handler) Handle(ctx context.Context, message *message.Message) error {
	start := time.Now().UTC()
	traceCtx, _ := h.tracer.Extract(opentracing.TextMap, headersTextMap(message.Headers))
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

	span := h.tracer.StartSpan(operationName, opts...)
	defer span.Finish()

	ctx = opentracing.ContextWithSpan(ctx, span)
	err := h.Handler.Handle(ctx, message)

	for k, v := range meta.Attributes(ctx) {
		span.SetTag(k, v)
	}

	span.SetTag(nsqDuration, stime.ToMilliseconds(time.Since(start)))

	if err != nil {
		setError(span, err)
	}

	return err
}

// NewProducer for opentracing.
func NewProducer(tracer opentracing.Tracer, p producer.Producer) *Producer {
	return &Producer{tracer: tracer, Producer: p}
}

// Producer for opentracing.
type Producer struct {
	tracer opentracing.Tracer
	producer.Producer
}

func (p *Producer) Publish(ctx context.Context, topic string, message *message.Message) error {
	start := time.Now().UTC()
	operationName := fmt.Sprintf("publish %s", topic)
	opts := []opentracing.StartSpanOption{
		opentracing.Tag{Key: nsqStartTime, Value: start.Format(time.RFC3339)},
		opentracing.Tag{Key: nsqBody, Value: string(message.Body)},
		opentracing.Tag{Key: nsqTopic, Value: topic},
		opentracing.Tag{Key: component, Value: nsqComponent},
		ext.SpanKindProducer,
	}

	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, p.tracer, operationName, opts...)
	defer span.Finish()

	if err := p.tracer.Inject(span.Context(), opentracing.TextMap, headersTextMap(message.Headers)); err != nil {
		return err
	}

	err := p.Producer.Publish(ctx, topic, message)

	for k, v := range meta.Attributes(ctx) {
		span.SetTag(k, v)
	}

	span.SetTag(nsqDuration, stime.ToMilliseconds(time.Since(start)))

	if err != nil {
		setError(span, err)
	}

	return err
}

func setError(span opentracing.Span, err error) {
	ext.Error.Set(span, true)
	span.LogFields(log.String("event", "error"), log.String("message", err.Error()))
}
