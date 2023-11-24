package metrics

import (
	"context"
	"time"

	"github.com/alexfalkowski/go-service/nsq"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// NewConsumer for metrics.
func NewConsumer(topic, channel string, meter metric.Meter, handler nsq.Consumer) (*Consumer, error) {
	started, err := meter.Int64Counter("nsq_consumer_started_total", metric.WithDescription("Total number of messages started to be consumed."))
	if err != nil {
		return nil, err
	}

	received, err := meter.Int64Counter("nsq_consumer_msg_received_total", metric.WithDescription("Total number of messages consumed."))
	if err != nil {
		return nil, err
	}

	handled, err := meter.Int64Counter("nsq_consumer_handled_total", metric.WithDescription("Total number of messages consumed, regardless of success or failure."))
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

	h := &Consumer{
		opts: opts, started: started, received: received, handled: handled, handledHist: handledHist,
		Consumer: handler,
	}

	return h, nil
}

// Consumer for metrics.
type Consumer struct {
	opts        metric.MeasurementOption
	started     metric.Int64Counter
	received    metric.Int64Counter
	handled     metric.Int64Counter
	handledHist metric.Float64Histogram

	nsq.Consumer
}

func (h *Consumer) Consume(ctx context.Context, message *nsq.Message) error {
	start := time.Now()

	h.started.Add(ctx, 1, h.opts)
	h.received.Add(ctx, 1, h.opts)

	if err := h.Consumer.Consume(ctx, message); err != nil {
		return err
	}

	h.handled.Add(ctx, 1, h.opts)
	h.handledHist.Record(ctx, time.Since(start).Seconds(), h.opts)

	return nil
}

// NewProducer for metrics.
func NewProducer(meter metric.Meter, producer nsq.Producer) (*Producer, error) {
	started, err := meter.Int64Counter("nsq_producer_started_total", metric.WithDescription("Total number of messages started by the producer."))
	if err != nil {
		return nil, err
	}

	sent, err := meter.Int64Counter("nsq_producer_msg_sent_total", metric.WithDescription("Total number of stream messages sent by the producer."))
	if err != nil {
		return nil, err
	}

	handled, err := meter.Int64Counter("nsq_producer_handled_total", metric.WithDescription("Total number of messages published, regardless of success or failure."))
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
	started     metric.Int64Counter
	sent        metric.Int64Counter
	handled     metric.Int64Counter
	handledHist metric.Float64Histogram

	nsq.Producer
}

func (p *Producer) Produce(ctx context.Context, topic string, message *nsq.Message) error {
	start := time.Now()
	opts := metric.WithAttributes(
		attribute.Key("nsq_topic").String(topic),
	)

	p.started.Add(ctx, 1, opts)
	p.sent.Add(ctx, 1, opts)

	err := p.Producer.Produce(ctx, topic, message)
	if err != nil {
		return err
	}

	p.handled.Add(ctx, 1, opts)
	p.handledHist.Record(ctx, time.Since(start).Seconds(), opts)

	return nil
}
