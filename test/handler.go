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
	mux sync.RWMutex
}

func (h *Handler) Message() *message.Message {
	h.mux.RLock()
	defer h.mux.RUnlock()

	return h.m
}

func (h *Handler) Handle(ctx context.Context, m *message.Message) error {
	h.mux.Lock()
	defer h.mux.Unlock()

	h.m = m

	meta.WithAttribute(ctx, "test", "test")

	return h.err
}
