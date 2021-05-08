package zap

import (
	"context"

	"github.com/alexfalkowski/go-service/pkg/meta"
	"github.com/alexfalkowski/go-service/pkg/time"
	"github.com/alexfalkowski/go-service/pkg/transport/nsq/handler"
	"github.com/alexfalkowski/go-service/pkg/transport/nsq/message"
	pkgMeta "github.com/alexfalkowski/go-service/pkg/transport/nsq/meta"
	"github.com/alexfalkowski/go-service/pkg/transport/nsq/producer"
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
	component    = "component"
	nsqComponent = "nsq"
	consumerKind = "consumer"
	producerKind = "producer"
)

// NewHandler for zap.
func NewHandler(logger *zap.Logger, h handler.Handler) handler.Handler {
	return &loggerHandler{logger: logger, Handler: h}
}

type loggerHandler struct {
	logger *zap.Logger

	handler.Handler
}

func (h *loggerHandler) Handle(ctx context.Context, message *message.Message) (context.Context, error) {
	start := time.Now().UTC()
	ctx, err := h.Handler.Handle(ctx, message)
	fields := []zapcore.Field{
		zap.Int64(nsqDuration, time.ToMilliseconds(time.Since(start))),
		zap.String(nsqStartTime, start.Format(time.RFC3339)),
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

		return ctx, err
	}

	h.logger.Info("finished call with success", fields...)

	return ctx, nil
}

// NewProducer for zap.
func NewProducer(logger *zap.Logger, p producer.Producer) producer.Producer {
	return &loggerProducer{logger: logger, Producer: p}
}

type loggerProducer struct {
	logger *zap.Logger

	producer.Producer
}

func (p *loggerProducer) Publish(ctx context.Context, topic string, message *message.Message) (context.Context, error) {
	start := time.Now().UTC()
	ctx, err := p.Producer.Publish(ctx, topic, message)
	fields := []zapcore.Field{
		zap.Int64(nsqDuration, time.ToMilliseconds(time.Since(start))),
		zap.String(nsqStartTime, start.Format(time.RFC3339)),
		zap.String(pkgMeta.Topic, topic),
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

		return ctx, err
	}

	p.logger.Info("finished call with success", fields...)

	return ctx, nil
}
