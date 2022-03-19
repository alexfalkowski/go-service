package test

import (
	"context"
	"sync"

	"github.com/alexfalkowski/go-service/meta"
	"github.com/alexfalkowski/go-service/transport/nsq/message"
)

// NewHandler for test.
func NewHandler(err error) *Handler {
	return &Handler{err: err}
}

// Handler for test.
type Handler struct {
	m   *message.Message
	err error
	mux sync.Mutex
}

func (h *Handler) Message() *message.Message {
	h.mux.Lock()
	defer h.mux.Unlock()

	return h.m
}

func (h *Handler) Handle(ctx context.Context, m *message.Message) (context.Context, error) {
	h.mux.Lock()
	defer h.mux.Unlock()

	h.m = m

	ctx = meta.WithAttribute(ctx, "test", "test")

	return ctx, h.err
}
