package opentracing

import (
	"context"
	"fmt"

	"github.com/alexfalkowski/go-service/meta"
	"github.com/alexfalkowski/go-service/trace/opentracing"
	"github.com/alexfalkowski/go-service/transport/nsq/handler"
	"github.com/alexfalkowski/go-service/transport/nsq/message"
	"github.com/alexfalkowski/go-service/transport/nsq/producer"
	"github.com/alexfalkowski/go-service/version"
	otr "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"go.uber.org/fx"
)

const (
	nsqID        = "nsq.id"
	nsqBody      = "nsq.body"
	nsqTimestamp = "nsq.timestamp"
	nsqAttempts  = "nsq.attempts"
	nsqAddress   = "nsq.address"
	nsqTopic     = "nsq.topic"
	nsqChannel   = "nsq.channel"
	component    = "component"
	nsqComponent = "nsq"
)

// TracerParams for otr.
type TracerParams struct {
	fx.In

	Lifecycle fx.Lifecycle
	Config    *opentracing.Config
	Version   version.Version
}

// NewTracer for otr.
func NewTracer(params TracerParams) (Tracer, error) {
	return opentracing.NewTracer(opentracing.TracerParams{Lifecycle: params.Lifecycle, Name: "nsq", Config: params.Config, Version: params.Version})
}

// Tracer for otr.
type Tracer otr.Tracer

// NewHandler for otr.
func NewHandler(topic, channel string, tracer Tracer, h handler.Handler) *Handler {
	return &Handler{topic: topic, channel: channel, tracer: tracer, Handler: h}
}

// Handler for otr.
type Handler struct {
	topic, channel string
	tracer         Tracer

	handler.Handler
}

func (h *Handler) Handle(ctx context.Context, message *message.Message) error {
	traceCtx, _ := h.tracer.Extract(otr.TextMap, headersTextMap(message.Headers))
	operationName := fmt.Sprintf("consume %s:%s", h.topic, h.channel)
	opts := []otr.StartSpanOption{
		ext.RPCServerOption(traceCtx),
		otr.Tag{Key: nsqTopic, Value: h.topic},
		otr.Tag{Key: nsqChannel, Value: h.channel},
		otr.Tag{Key: nsqID, Value: string(message.ID[:])},
		otr.Tag{Key: nsqBody, Value: string(message.Body)},
		otr.Tag{Key: nsqTimestamp, Value: message.Timestamp},
		otr.Tag{Key: nsqAttempts, Value: message.Attempts},
		otr.Tag{Key: nsqAddress, Value: message.NSQDAddress},
		otr.Tag{Key: component, Value: nsqComponent},
		ext.SpanKindConsumer,
	}

	span := h.tracer.StartSpan(operationName, opts...)
	defer span.Finish()

	ctx = otr.ContextWithSpan(ctx, span)

	err := h.Handler.Handle(ctx, message)
	if err != nil {
		opentracing.SetError(span, err)
	}

	for k, v := range meta.Attributes(ctx) {
		span.SetTag(k, v)
	}

	return err
}

// NewProducer for otr.
func NewProducer(tracer Tracer, p producer.Producer) *Producer {
	return &Producer{tracer: tracer, Producer: p}
}

// Producer for otr.
type Producer struct {
	tracer Tracer
	producer.Producer
}

func (p *Producer) Publish(ctx context.Context, topic string, message *message.Message) error {
	operationName := fmt.Sprintf("publish %s", topic)
	opts := []otr.StartSpanOption{
		otr.Tag{Key: nsqBody, Value: string(message.Body)},
		otr.Tag{Key: nsqTopic, Value: topic},
		otr.Tag{Key: component, Value: nsqComponent},
		ext.SpanKindProducer,
	}

	span, ctx := otr.StartSpanFromContextWithTracer(ctx, p.tracer, operationName, opts...)
	defer span.Finish()

	if err := p.tracer.Inject(span.Context(), otr.TextMap, headersTextMap(message.Headers)); err != nil {
		return err
	}

	err := p.Producer.Publish(ctx, topic, message)
	if err != nil {
		opentracing.SetError(span, err)
	}

	for k, v := range meta.Attributes(ctx) {
		span.SetTag(k, v)
	}

	return err
}
