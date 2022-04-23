package handler

import (
	"context"

	"github.com/alexfalkowski/go-service/compressor"
	"github.com/alexfalkowski/go-service/marshaller"
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
	Compressor compressor.Compressor
	Marshaller marshaller.Marshaller
}

// New handler for NSQ.
// nolint:ireturn
func New(params Params) nsq.Handler {
	return &handler{Compressor: params.Compressor, Marshaller: params.Marshaller, Handler: params.Handler}
}

type handler struct {
	Compressor compressor.Compressor
	Marshaller marshaller.Marshaller
	Handler
}

func (h *handler) HandleMessage(m *nsq.Message) error {
	if m.Body == nil {
		return nil
	}

	var msg message.Message

	bytes, err := h.Compressor.Decompress(m.Body)
	if err != nil {
		return err
	}

	if err := h.Marshaller.Unmarshal(bytes, &msg); err != nil {
		return err
	}

	msg.Message = m

	ctx := context.Background()

	return h.Handle(ctx, &msg)
}
