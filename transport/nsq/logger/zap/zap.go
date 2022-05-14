package zap

import (
	"context"
	"time"

	"github.com/alexfalkowski/go-service/meta"
	stime "github.com/alexfalkowski/go-service/time"
	"github.com/alexfalkowski/go-service/transport/nsq/handler"
	"github.com/alexfalkowski/go-service/transport/nsq/message"
	"github.com/alexfalkowski/go-service/transport/nsq/producer"
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
	component    = "component"
	nsqComponent = "nsq"
	consumerKind = "consumer"
	producerKind = "producer"
)

// HandlerParams for zap.
type HandlerParams struct {
	Topic, Channel string
	Logger         *zap.Logger
	Handler        handler.Handler
}

// NewHandler for zap.
func NewHandler(params HandlerParams) *Handler {
	return &Handler{topic: params.Topic, channel: params.Channel, logger: params.Logger, Handler: params.Handler}
}

// Handler for zap.
type Handler struct {
	topic, channel string
	logger         *zap.Logger

	handler.Handler
}

func (h *Handler) Handle(ctx context.Context, message *message.Message) error {
	start := time.Now().UTC()
	err := h.Handler.Handle(ctx, message)
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
		zap.String("span.kind", consumerKind),
		zap.String(component, nsqComponent),
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

// ProducerParams for zap.
type ProducerParams struct {
	Logger   *zap.Logger
	Producer producer.Producer
}

// NewProducer for zap.
func NewProducer(params ProducerParams) *Producer {
	return &Producer{logger: params.Logger, Producer: params.Producer}
}

// Producer for zap.
type Producer struct {
	logger *zap.Logger

	producer.Producer
}

func (p *Producer) Publish(ctx context.Context, topic string, message *message.Message) error {
	start := time.Now().UTC()
	err := p.Producer.Publish(ctx, topic, message)
	fields := []zapcore.Field{
		zap.Int64(nsqDuration, stime.ToMilliseconds(time.Since(start))),
		zap.String(nsqStartTime, start.Format(time.RFC3339)),
		zap.String(nsqTopic, topic),
		zap.ByteString(nsqBody, message.Body),
		zap.String("span.kind", producerKind),
		zap.String(component, nsqComponent),
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
