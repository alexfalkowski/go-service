package meta

import (
	"context"

	"github.com/alexfalkowski/go-service/nsq"
	"github.com/alexfalkowski/go-service/transport/meta"
	"github.com/google/uuid"
)

// NewConsumer for meta.
func NewConsumer(h nsq.Consumer) *Consumer {
	return &Consumer{Consumer: h}
}

// Consumer for meta.
type Consumer struct {
	nsq.Consumer
}

func (h *Consumer) Consume(ctx context.Context, message *nsq.Message) error {
	ctx = meta.WithUserAgent(ctx, extractUserAgent(ctx, message.Headers))

	requestID := extractRequestID(ctx, message.Headers)
	if requestID == "" {
		requestID = uuid.New().String()
	}

	ctx = meta.WithRequestID(ctx, requestID)

	return h.Consumer.Consume(ctx, message)
}

// NewProducer for meta.
func NewProducer(userAgent string, p nsq.Producer) *Producer {
	return &Producer{userAgent: userAgent, Producer: p}
}

// Producer for meta.
type Producer struct {
	userAgent string
	nsq.Producer
}

func (p *Producer) Produce(ctx context.Context, topic string, message *nsq.Message) error {
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

	message.Headers["request-id"] = requestID
	ctx = meta.WithRequestID(ctx, requestID)

	return p.Producer.Produce(ctx, topic, message)
}

func extractUserAgent(ctx context.Context, headers nsq.Headers) string {
	if userAgent := headers["user-agent"]; userAgent != "" {
		return userAgent
	}

	return meta.UserAgent(ctx)
}

func extractRequestID(ctx context.Context, headers nsq.Headers) string {
	if requestID := headers["request-id"]; requestID != "" {
		return requestID
	}

	return meta.RequestID(ctx)
}
