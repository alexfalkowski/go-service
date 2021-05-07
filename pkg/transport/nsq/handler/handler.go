package handler

import (
	"context"

	"github.com/nsqio/go-nsq"
)

// Handler for NSQ.
type Handler interface {
	Handle(ctx context.Context, message *nsq.Message) (context.Context, error)
}

// NewHandler for NSQ.
func NewHandler(h Handler) nsq.Handler {
	return &handler{Handler: h}
}

type handler struct {
	Handler
}

func (h *handler) HandleMessage(m *nsq.Message) error {
	_, err := h.Handle(context.Background(), m)

	return err
}
