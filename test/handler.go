package test

import (
	"context"
	"sync"

	"github.com/alexfalkowski/go-service/pkg/meta"
	"github.com/nsqio/go-nsq"
)

// NewHandler for test.
func NewHandler() *Handler {
	return &Handler{}
}

// Handler for test.
type Handler struct {
	m   *nsq.Message
	mux sync.Mutex
}

func (h *Handler) Message() *nsq.Message {
	h.mux.Lock()
	defer h.mux.Unlock()

	return h.m
}

func (h *Handler) Handle(ctx context.Context, m *nsq.Message) (context.Context, error) {
	h.mux.Lock()
	defer h.mux.Unlock()

	h.m = m

	ctx = meta.WithAttribute(ctx, "test", "test")

	return ctx, nil
}
