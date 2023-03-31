package meta

import (
	"context"
	"strings"

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

func (h *Handler) Handle(ctx context.Context, message *message.Message) error {
	ctx = meta.WithUserAgent(ctx, extractUserAgent(ctx, message.Headers))

	requestID := extractRequestID(ctx, message.Headers)
	if requestID == "" {
		requestID = uuid.New().String()
	}

	ctx = meta.WithRequestID(ctx, requestID)
	ctx = meta.WithRemoteAddress(ctx, extractRemoteAddress(ctx, message.Headers))

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

func (p *Producer) Publish(ctx context.Context, topic string, message *message.Message) error {
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

	ctx = meta.WithRemoteAddress(ctx, extractRemoteAddress(ctx, message.Headers))

	return p.Producer.Publish(ctx, topic, message)
}

func extractUserAgent(ctx context.Context, headers message.Headers) string {
	if userAgent := headers["user-agent"]; userAgent != "" {
		return userAgent
	}

	return meta.UserAgent(ctx)
}

func extractRequestID(ctx context.Context, headers message.Headers) string {
	if requestID := headers["request-id"]; requestID != "" {
		return requestID
	}

	return meta.RequestID(ctx)
}

func extractRemoteAddress(ctx context.Context, headers message.Headers) string {
	if forwardedFor := headers["forwarded-for"]; forwardedFor != "" {
		return strings.Split(forwardedFor, ",")[0]
	}

	return meta.RemoteAddress(ctx)
}
