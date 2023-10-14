package metrics

import (
	"context"
	"time"

	"github.com/alexfalkowski/go-service/transport/nsq/handler"
	"github.com/alexfalkowski/go-service/transport/nsq/message"
	"github.com/alexfalkowski/go-service/transport/nsq/producer"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// NewHandler for metrics.
func NewHandler(topic, channel string, meter metric.Meter, handler handler.Handler) (*Handler, error) {
	started, err := meter.Float64Counter("nsq_consumer_started_total", metric.WithDescription("Total number of messages started to be consumed."))
	if err != nil {
		return nil, err
	}

	received, err := meter.Float64Counter("nsq_consumer_msg_received_total", metric.WithDescription("Total number of messages consumed."))
	if err != nil {
		return nil, err
	}

	handled, err := meter.Float64Counter("nsq_consumer_handled_total", metric.WithDescription("Total number of messages consumed, regardless of success or failure."))
	if err != nil {
		return nil, err
	}

	handledHist, err := meter.Float64Histogram("nsq_consumer_handling_seconds",
		metric.WithDescription("Histogram of response latency (seconds) of messages that had been consumed."))
	if err != nil {
		return nil, err
	}

	opts := metric.WithAttributes(
		attribute.Key("nsq_topic").String(topic),
		attribute.Key("nsq_channel").String(channel),
	)

	h := &Handler{
		opts: opts, started: started, received: received, handled: handled, handledHist: handledHist,
		Handler: handler,
	}

	return h, nil
}

// Handler for metrics.
type Handler struct {
	opts        metric.MeasurementOption
	started     metric.Float64Counter
	received    metric.Float64Counter
	handled     metric.Float64Counter
	handledHist metric.Float64Histogram

	handler.Handler
}

func (h *Handler) Handle(ctx context.Context, message *message.Message) error {
	st := time.Now()

	h.started.Add(ctx, 1, h.opts)
	h.received.Add(ctx, 1, h.opts)

	if err := h.Handler.Handle(ctx, message); err != nil {
		return err
	}

	h.handled.Add(ctx, 1, h.opts)
	h.handledHist.Record(ctx, time.Since(st).Seconds(), h.opts)

	return nil
}

// NewProducer for metrics.
func NewProducer(meter metric.Meter, producer producer.Producer) (*Producer, error) {
	started, err := meter.Float64Counter("nsq_producer_started_total", metric.WithDescription("Total number of messages started by the producer."))
	if err != nil {
		return nil, err
	}

	sent, err := meter.Float64Counter("nsq_producer_msg_sent_total", metric.WithDescription("Total number of stream messages sent by the producer."))
	if err != nil {
		return nil, err
	}

	handled, err := meter.Float64Counter("nsq_producer_handled_total", metric.WithDescription("Total number of messages published, regardless of success or failure."))
	if err != nil {
		return nil, err
	}

	handledHist, err := meter.Float64Histogram("nsq_producer_handling_seconds",
		metric.WithDescription("Histogram of response latency (seconds) of messages that had been application-level handled by the producer."))
	if err != nil {
		return nil, err
	}

	p := &Producer{
		started: started, sent: sent, handled: handled, handledHist: handledHist,
		Producer: producer,
	}

	return p, nil
}

// Producer for metrics.
type Producer struct {
	started     metric.Float64Counter
	sent        metric.Float64Counter
	handled     metric.Float64Counter
	handledHist metric.Float64Histogram

	producer.Producer
}

func (p *Producer) Publish(ctx context.Context, topic string, message *message.Message) error {
	st := time.Now()
	opts := metric.WithAttributes(
		attribute.Key("nsq_topic").String(topic),
	)

	p.started.Add(ctx, 1, opts)
	p.sent.Add(ctx, 1, opts)

	err := p.Producer.Publish(ctx, topic, message)
	if err != nil {
		return err
	}

	p.handled.Add(ctx, 1, opts)
	p.handledHist.Record(ctx, time.Since(st).Seconds(), opts)

	return nil
}
