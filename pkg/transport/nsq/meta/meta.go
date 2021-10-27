package meta

import (
	"context"

	"github.com/alexfalkowski/go-service/pkg/transport/meta"
	"github.com/alexfalkowski/go-service/pkg/transport/nsq/handler"
	"github.com/alexfalkowski/go-service/pkg/transport/nsq/message"
	"github.com/alexfalkowski/go-service/pkg/transport/nsq/producer"
	"github.com/google/uuid"
)

// NewHandler for meta.
func NewHandler(h handler.Handler) handler.Handler {
	return &metaHandler{Handler: h}
}

type metaHandler struct {
	handler.Handler
}

func (h *metaHandler) Handle(ctx context.Context, message *message.Message) (context.Context, error) {
	ctx = meta.WithUserAgent(ctx, message.Headers["user-agent"])

	requestID, ok := message.Headers["request-id"]
	if !ok {
		requestID = uuid.New().String()
	}

	ctx = meta.WithRequestID(ctx, requestID)
	message.Headers["request-id"] = requestID

	return h.Handler.Handle(ctx, message)
}

// NewProducer for meta.
func NewProducer(userAgent string, p producer.Producer) producer.Producer {
	return &metaProducer{userAgent: userAgent, Producer: p}
}

type metaProducer struct {
	userAgent string
	producer.Producer
}

func (p *metaProducer) Publish(ctx context.Context, topic string, message *message.Message) (context.Context, error) {
	message.Headers["user-agent"] = p.userAgent
	ctx = meta.WithUserAgent(ctx, p.userAgent)

	requestID := meta.RequestID(ctx)
	if requestID == "" {
		requestID = uuid.New().String()
	}

	ctx = meta.WithRequestID(ctx, requestID)
	message.Headers["request-id"] = requestID

	return p.Producer.Publish(ctx, topic, message)
}
