package test

import (
	"context"
	"sync"

	"github.com/alexfalkowski/go-service/meta"
	"github.com/alexfalkowski/go-service/nsq"
)

// NewConsumer for test.
func NewConsumer(err error) *Consumer {
	return &Consumer{err: err}
}

// Consumer for test.
type Consumer struct {
	m   *nsq.Message
	err error
	mux sync.RWMutex
}

func (h *Consumer) Message() *nsq.Message {
	h.mux.RLock()
	defer h.mux.RUnlock()

	return h.m
}

func (h *Consumer) Consume(ctx context.Context, m *nsq.Message) error {
	h.mux.Lock()
	defer h.mux.Unlock()

	h.m = m

	meta.WithAttribute(ctx, "test", "test")

	return h.err
}
