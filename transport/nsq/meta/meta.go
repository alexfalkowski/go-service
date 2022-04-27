package meta

import (
	"context"
	"strings"

	"github.com/alexfalkowski/go-service/meta"
	"github.com/alexfalkowski/go-service/net"
	tmeta "github.com/alexfalkowski/go-service/transport/meta"
	"github.com/alexfalkowski/go-service/transport/nsq/handler"
	"github.com/alexfalkowski/go-service/transport/nsq/message"
	"github.com/alexfalkowski/go-service/transport/nsq/producer"
	"github.com/alexfalkowski/go-service/version"
	"github.com/google/uuid"
)

// NewHandler for meta.
func NewHandler(version version.Version, h handler.Handler) *Handler {
	return &Handler{version: version, Handler: h}
}

// Handler for meta.
type Handler struct {
	version version.Version
	handler.Handler
}

func (h *Handler) Handle(ctx context.Context, message *message.Message) error {
	ctx = meta.WithVersion(ctx, string(h.version))
	ctx = tmeta.WithUserAgent(ctx, extractUserAgent(ctx, message.Headers))

	requestID := extractRequestID(ctx, message.Headers)
	if requestID == "" {
		requestID = uuid.New().String()
	}

	ctx = tmeta.WithRequestID(ctx, requestID)
	ctx = tmeta.WithRemoteAddress(ctx, extractRemoteAddress(ctx, message.Headers))

	return h.Handler.Handle(ctx, message)
}

// NewProducer for meta.
func NewProducer(userAgent string, version version.Version, p producer.Producer) *Producer {
	return &Producer{userAgent: userAgent, version: version, Producer: p}
}

// Producer for meta.
type Producer struct {
	userAgent string
	version   version.Version
	producer.Producer
}

func (p *Producer) Publish(ctx context.Context, topic string, message *message.Message) error {
	message.Headers["version"] = string(p.version)
	ctx = meta.WithVersion(ctx, string(p.version))

	userAgent := tmeta.UserAgent(ctx)
	if userAgent == "" {
		userAgent = p.userAgent
	}

	message.Headers["user-agent"] = userAgent
	ctx = tmeta.WithUserAgent(ctx, userAgent)

	requestID := tmeta.RequestID(ctx)
	if requestID == "" {
		requestID = uuid.New().String()
	}

	message.Headers["request-id"] = requestID
	ctx = tmeta.WithRequestID(ctx, requestID)

	remoteAddress := tmeta.RemoteAddress(ctx)
	if remoteAddress == "" {
		remoteAddress = net.OutboundIP(ctx)
	}

	message.Headers["forwarded-for"] = remoteAddress
	ctx = tmeta.WithRemoteAddress(ctx, remoteAddress)

	return p.Producer.Publish(ctx, topic, message)
}

func extractUserAgent(ctx context.Context, headers message.Headers) string {
	if userAgent := headers["user-agent"]; userAgent != "" {
		return userAgent
	}

	return tmeta.UserAgent(ctx)
}

func extractRequestID(ctx context.Context, headers message.Headers) string {
	if requestID := headers["request-id"]; requestID != "" {
		return requestID
	}

	return tmeta.RequestID(ctx)
}

func extractRemoteAddress(ctx context.Context, headers message.Headers) string {
	if forwardedFor := headers["forwarded-for"]; forwardedFor != "" {
		return strings.Split(forwardedFor, ",")[0]
	}

	return tmeta.RemoteAddress(ctx)
}
