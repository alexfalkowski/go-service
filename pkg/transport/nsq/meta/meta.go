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
	requestID, ok := message.Headers["request-id"]
	if !ok {
		requestID = uuid.New().String()
	}

	ctx = meta.WithRequestID(ctx, requestID)
	message.Headers["request-id"] = requestID

	return h.Handler.Handle(ctx, message)
}

// NewProducer for meta.
func NewProducer(p producer.Producer) producer.Producer {
	return &metaProducer{Producer: p}
}

type metaProducer struct {
	producer.Producer
}

func (p *metaProducer) Publish(ctx context.Context, topic string, message *message.Message) (context.Context, error) {
	requestID := meta.RequestID(ctx)
	if requestID == "" {
		requestID = uuid.New().String()
	}

	ctx = meta.WithRequestID(ctx, requestID)
	message.Headers["request-id"] = requestID

	return p.Producer.Publish(ctx, topic, message)
}
