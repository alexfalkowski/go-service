package zap

import (
	"github.com/alexfalkowski/go-service/pkg/time"
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
func NewHandler(logger *zap.Logger, nh nsq.Handler) nsq.Handler {
	return &handler{logger: logger, Handler: nh}
}

type handler struct {
	logger *zap.Logger

	nsq.Handler
}

func (h *handler) HandleMessage(m *nsq.Message) error {
	start := time.Now().UTC()
	err := h.Handler.HandleMessage(m)
	fields := []zapcore.Field{
		zap.Int64(nsqDuration, time.ToMilliseconds(time.Since(start))),
		zap.String(nsqStartTime, start.Format(time.RFC3339)),
		zap.ByteString(nsqID, m.ID[:]),
		zap.ByteString(nsqBody, m.Body),
		zap.Int64(nsqTimestamp, m.Timestamp),
		zap.Uint16(nsqAttempts, m.Attempts),
		zap.String(nsqAddress, m.NSQDAddress),
		zap.String("span.kind", consumer),
		zap.String(component, nsqComponent),
	}

	if err != nil {
		fields = append(fields, zap.Error(err))
		h.logger.Error("finished call with error", fields...)

		return err
	}

	h.logger.Info("finished call with success", fields...)

	return nil
}
