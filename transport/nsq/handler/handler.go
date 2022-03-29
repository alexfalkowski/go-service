package handler

import (
	"context"

	"github.com/alexfalkowski/go-service/transport/nsq/message"
	"github.com/nsqio/go-nsq"
)

// Handler for NSQ.
type Handler interface {
	Handle(ctx context.Context, message *message.Message) error
}

// New handler for NSQ.
// nolint:ireturn
func New(h Handler) nsq.Handler {
	return &handler{Handler: h}
}

type handler struct {
	Handler
}

func (h *handler) HandleMessage(m *nsq.Message) error {
	if m.Body == nil {
		return nil
	}

	var msg message.Message
	if err := message.Unmarshal(m.Body, &msg); err != nil {
		return err
	}

	msg.Message = m

	ctx := context.Background()

	return h.Handle(ctx, &msg)
}
