package zap

import (
	"context"
	"time"

	"github.com/alexfalkowski/go-service/meta"
	"github.com/alexfalkowski/go-service/nsq"
	stime "github.com/alexfalkowski/go-service/time"
	tm "github.com/alexfalkowski/go-service/transport/meta"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	service = "nsq"
)

// NewConsumer for zap.
func NewConsumer(topic, channel string, logger *zap.Logger, handler nsq.Consumer) *Consumer {
	return &Consumer{topic: topic, channel: channel, logger: logger, Consumer: handler}
}

// Consumer for zap.
type Consumer struct {
	topic, channel string
	logger         *zap.Logger

	nsq.Consumer
}

func (h *Consumer) Consume(ctx context.Context, message *nsq.Message) error {
	start := time.Now()
	err := h.Consumer.Consume(ctx, message)
	fields := []zapcore.Field{
		zap.Int64(tm.DurationKey, stime.ToMilliseconds(time.Since(start))),
		zap.String(tm.StartTimeKey, start.Format(time.RFC3339)),
		zap.String(tm.ServiceKey, service),
		zap.String(tm.PathKey, h.topic+":"+h.channel),
	}

	for k, v := range meta.Attributes(ctx) {
		fields = append(fields, zap.String(k, v))
	}

	if err != nil {
		fields = append(fields, zap.Error(err))
		h.logger.Error("finished call with error", fields...)

		return err
	}

	h.logger.Info("finished call with success", fields...)

	return nil
}

// NewProducer for zap.
func NewProducer(logger *zap.Logger, producer nsq.Producer) *Producer {
	return &Producer{logger: logger, Producer: producer}
}

// Producer for zap.
type Producer struct {
	logger *zap.Logger

	nsq.Producer
}

func (p *Producer) Produce(ctx context.Context, topic string, message *nsq.Message) error {
	start := time.Now()
	err := p.Producer.Produce(ctx, topic, message)
	fields := []zapcore.Field{
		zap.Int64(tm.DurationKey, stime.ToMilliseconds(time.Since(start))),
		zap.String(tm.StartTimeKey, start.Format(time.RFC3339)),
		zap.String(tm.ServiceKey, service),
		zap.String(tm.PathKey, topic),
	}

	for k, v := range meta.Attributes(ctx) {
		fields = append(fields, zap.String(k, v))
	}

	if err != nil {
		fields = append(fields, zap.Error(err))
		p.logger.Error("finished call with error", fields...)

		return err
	}

	p.logger.Info("finished call with success", fields...)

	return nil
}
