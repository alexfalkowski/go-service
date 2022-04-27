package handler

import (
	"context"

	"github.com/alexfalkowski/go-service/transport/nsq/marshaller"
	"github.com/alexfalkowski/go-service/transport/nsq/message"
	"github.com/nsqio/go-nsq"
)

// Handler for NSQ.
type Handler interface {
	Handle(ctx context.Context, message *message.Message) error
}

// Params for handler.
type Params struct {
	Handler    Handler
	Marshaller marshaller.Marshaller
}

// New handler for NSQ.
func New(params Params) nsq.Handler {
	return &handler{Marshaller: params.Marshaller, Handler: params.Handler}
}

type handler struct {
	Marshaller marshaller.Marshaller
	Handler
}

func (h *handler) HandleMessage(m *nsq.Message) error {
	var msg message.Message
	if err := h.Marshaller.Unmarshal(m.Body, &msg); err != nil {
		return err
	}

	msg.Message = m

	ctx := context.Background()

	return h.Handle(ctx, &msg)
}
