package tracer

import (
	"context"
	"fmt"

	"github.com/alexfalkowski/go-service/meta"
	"github.com/alexfalkowski/go-service/telemetry/tracer"
	"github.com/alexfalkowski/go-service/transport/nsq/handler"
	"github.com/alexfalkowski/go-service/transport/nsq/message"
	"github.com/alexfalkowski/go-service/transport/nsq/producer"
	"github.com/alexfalkowski/go-service/version"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.19.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

// Params for tracer.
type Params struct {
	fx.In

	Lifecycle fx.Lifecycle
	Config    *tracer.Config
	Version   version.Version
}

// NewTracer for tracer.
func NewTracer(params Params) (Tracer, error) {
	return tracer.NewTracer(params.Lifecycle, "nsq", params.Version, params.Config)
}

// Tracer for tracer.
type Tracer trace.Tracer

// NewHandler for tracer.
func NewHandler(topic, channel string, tracer Tracer, h handler.Handler) *Handler {
	return &Handler{topic: topic, channel: channel, tracer: tracer, Handler: h}
}

// Handler for tracer.
type Handler struct {
	topic, channel string
	tracer         Tracer

	handler.Handler
}

func (h *Handler) Handle(ctx context.Context, message *message.Message) error {
	ctx = extract(ctx, message.Headers)

	operationName := fmt.Sprintf("consume %s:%s", h.topic, h.channel)
	attrs := []attribute.KeyValue{
		semconv.MessagingSystem("nsq"),
		semconv.MessagingSourceKindTopic,
		semconv.MessagingSourceName(h.topic),
	}

	ctx, span := h.tracer.Start(
		trace.ContextWithRemoteSpanContext(ctx, trace.SpanContextFromContext(ctx)),
		operationName,
		trace.WithSpanKind(trace.SpanKindConsumer),
		trace.WithAttributes(attrs...),
	)
	defer span.End()

	err := h.Handler.Handle(ctx, message)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
	}

	for k, v := range meta.Attributes(ctx) {
		span.SetAttributes(attribute.Key(k).String(v))
	}

	return err
}

// NewProducer for tracer.
func NewProducer(tracer Tracer, p producer.Producer) *Producer {
	return &Producer{tracer: tracer, Producer: p}
}

// Producer for tracer.
type Producer struct {
	tracer Tracer
	producer.Producer
}

func (p *Producer) Publish(ctx context.Context, topic string, message *message.Message) error {
	operationName := fmt.Sprintf("publish %s", topic)
	attrs := []attribute.KeyValue{
		semconv.MessagingSystem("nsq"),
		semconv.MessagingDestinationKindTopic,
		semconv.MessagingDestinationName(topic),
	}

	ctx, span := p.tracer.Start(
		ctx,
		operationName,
		trace.WithSpanKind(trace.SpanKindProducer),
		trace.WithAttributes(attrs...),
	)
	defer span.End()

	inject(ctx, message.Headers)

	err := p.Producer.Publish(ctx, topic, message)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
	}

	for k, v := range meta.Attributes(ctx) {
		span.SetAttributes(attribute.Key(k).String(v))
	}

	return err
}
