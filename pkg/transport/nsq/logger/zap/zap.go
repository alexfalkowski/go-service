package zap

import (
	"context"

	"github.com/alexfalkowski/go-service/pkg/meta"
	"github.com/alexfalkowski/go-service/pkg/time"
	"github.com/alexfalkowski/go-service/pkg/transport/nsq/handler"
	"github.com/nsqio/go-nsq"
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
	consumer     = "consumer"
)

// NewHandler for zap.
func NewHandler(logger *zap.Logger, h handler.Handler) handler.Handler {
	return &loggerHandler{logger: logger, Handler: h}
}

type loggerHandler struct {
	logger *zap.Logger

	handler.Handler
}

func (h *loggerHandler) Handle(ctx context.Context, message *nsq.Message) (context.Context, error) {
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
		zap.String("span.kind", consumer),
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
