package test

import (
	sync "sync"

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

func (h *Handler) HandleMessage(m *nsq.Message) error {
	h.mux.Lock()
	defer h.mux.Unlock()

	h.m = m

	return nil
}
