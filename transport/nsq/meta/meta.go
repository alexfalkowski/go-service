package meta

import (
	"context"

	"github.com/alexfalkowski/go-service/transport/meta"
	"github.com/alexfalkowski/go-service/transport/nsq/handler"
	"github.com/alexfalkowski/go-service/transport/nsq/message"
	"github.com/alexfalkowski/go-service/transport/nsq/producer"
	"github.com/google/uuid"
)

// NewHandler for meta.
func NewHandler(h handler.Handler) *Handler {
	return &Handler{Handler: h}
}

// Handler for meta.
type Handler struct {
	handler.Handler
}

func (h *Handler) Handle(ctx context.Context, message *message.Message) (context.Context, error) {
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
func NewProducer(userAgent string, p producer.Producer) *Producer {
	return &Producer{userAgent: userAgent, Producer: p}
}

// Producer for meta.
type Producer struct {
	userAgent string
	producer.Producer
}

func (p *Producer) Publish(ctx context.Context, topic string, message *message.Message) (context.Context, error) {
	userAgent := meta.UserAgent(ctx)
	if userAgent == "" {
		userAgent = p.userAgent
	}

	message.Headers["user-agent"] = userAgent
	ctx = meta.WithUserAgent(ctx, userAgent)

	requestID := meta.RequestID(ctx)
	if requestID == "" {
		requestID = uuid.New().String()
	}

	ctx = meta.WithRequestID(ctx, requestID)
	message.Headers["request-id"] = requestID

	return p.Producer.Publish(ctx, topic, message)
}
