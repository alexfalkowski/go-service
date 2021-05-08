package handler

import (
	"context"

	pkgMeta "github.com/alexfalkowski/go-service/pkg/meta"
	"github.com/alexfalkowski/go-service/pkg/transport/nsq/message"
	"github.com/alexfalkowski/go-service/pkg/transport/nsq/meta"
	"github.com/nsqio/go-nsq"
)

// Handler for NSQ.
type Handler interface {
	Handle(ctx context.Context, message *message.Message) (context.Context, error)
}

// NewHandler for NSQ.
func NewHandler(topic, channel string, h Handler) nsq.Handler {
	return &handler{topic: topic, channel: channel, Handler: h}
}

type handler struct {
	topic   string
	channel string

	Handler
}

func (h *handler) HandleMessage(m *nsq.Message) error {
	if m.Body == nil {
		return nil
	}

	var msg message.Message
	if err := message.Unmarshal(m.Body, &msg); err != nil {
		return nil
	}

	msg.Message = m

	ctx := context.Background()
	ctx = pkgMeta.WithAttribute(ctx, meta.Topic, h.topic)
	ctx = pkgMeta.WithAttribute(ctx, meta.Channel, h.channel)

	_, err := h.Handle(ctx, &msg)

	return err
}
