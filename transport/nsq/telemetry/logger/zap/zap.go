package zap

import (
	"context"
	"time"

	"github.com/alexfalkowski/go-service/meta"
	"github.com/alexfalkowski/go-service/nsq"
	stime "github.com/alexfalkowski/go-service/time"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
	kind         = "kind"
	nsqKind      = "nsq"
	consumerKind = "consumer"
	producerKind = "producer"
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
	start := time.Now().UTC()
	err := h.Consumer.Consume(ctx, message)
	fields := []zapcore.Field{
		zap.Int64(nsqDuration, stime.ToMilliseconds(time.Since(start))),
		zap.String(nsqStartTime, start.Format(time.RFC3339)),
		zap.String(nsqTopic, h.topic),
		zap.String(nsqChannel, h.channel),
		zap.ByteString(nsqID, message.ID[:]),
		zap.ByteString(nsqBody, message.Body),
		zap.Int64(nsqTimestamp, message.Timestamp),
		zap.Uint16(nsqAttempts, message.Attempts),
		zap.String(nsqAddress, message.NSQDAddress),
		zap.String("nsq.kind", consumerKind),
		zap.String(kind, nsqKind),
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
	start := time.Now().UTC()
	err := p.Producer.Produce(ctx, topic, message)
	fields := []zapcore.Field{
		zap.Int64(nsqDuration, stime.ToMilliseconds(time.Since(start))),
		zap.String(nsqStartTime, start.Format(time.RFC3339)),
		zap.String(nsqTopic, topic),
		zap.ByteString(nsqBody, message.Body),
		zap.String("nsq.kind", producerKind),
		zap.String(kind, nsqKind),
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
