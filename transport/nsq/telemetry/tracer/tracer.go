package tracer

import (
	"context"
	"fmt"

	"github.com/alexfalkowski/go-service/env"
	"github.com/alexfalkowski/go-service/meta"
	"github.com/alexfalkowski/go-service/nsq"
	"github.com/alexfalkowski/go-service/telemetry/tracer"
	"github.com/alexfalkowski/go-service/version"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

// Params for tracer.
type Params struct {
	fx.In

	Lifecycle   fx.Lifecycle
	Config      *tracer.Config
	Version     version.Version
	Environment env.Environment
}

// NewTracer for tracer.
func NewTracer(params Params) (Tracer, error) {
	return tracer.NewTracer(context.Background(), params.Lifecycle, "nsq", params.Environment, params.Version, params.Config)
}

// Tracer for tracer.
type Tracer trace.Tracer

// NewConsumer for tracer.
func NewConsumer(topic, channel string, tracer Tracer, h nsq.Consumer) *Consumer {
	return &Consumer{topic: topic, channel: channel, tracer: tracer, Consumer: h}
}

// Consumer for tracer.
type Consumer struct {
	topic, channel string
	tracer         Tracer

	nsq.Consumer
}

func (h *Consumer) Consume(ctx context.Context, message *nsq.Message) error {
	ctx = extract(ctx, message.Headers)

	operationName := fmt.Sprintf("consume %s:%s", h.topic, h.channel)
	attrs := []attribute.KeyValue{
		semconv.MessagingSystemKey.String("nsq"),
		semconv.MessagingDestinationName(h.topic),
	}

	ctx, span := h.tracer.Start(
		trace.ContextWithRemoteSpanContext(ctx, trace.SpanContextFromContext(ctx)),
		operationName,
		trace.WithSpanKind(trace.SpanKindConsumer),
		trace.WithAttributes(attrs...),
	)
	defer span.End()

	err := h.Consumer.Consume(ctx, message)
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
func NewProducer(tracer Tracer, p nsq.Producer) *Producer {
	return &Producer{tracer: tracer, Producer: p}
}

// Producer for tracer.
type Producer struct {
	tracer Tracer
	nsq.Producer
}

func (p *Producer) Produce(ctx context.Context, topic string, message *nsq.Message) error {
	operationName := "publish " + topic
	attrs := []attribute.KeyValue{
		semconv.MessagingSystemKey.String("nsq"),
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

	err := p.Producer.Produce(ctx, topic, message)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
	}

	for k, v := range meta.Attributes(ctx) {
		span.SetAttributes(attribute.Key(k).String(v))
	}

	return err
}
