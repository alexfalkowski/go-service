package nsq

import (
	"context"

	"github.com/nsqio/go-nsq"
)

// Consumer for NSQ.
type Consumer interface {
	Consume(ctx context.Context, message *Message) error
}

// NewConsumer handler for NSQ.
func NewConsumer(h Consumer, m Marshaller) nsq.Handler {
	return &handler{Marshaller: m, Consumer: h}
}

type handler struct {
	Marshaller Marshaller

	Consumer
}

func (h *handler) HandleMessage(m *nsq.Message) error {
	var msg Message
	if err := h.Marshaller.Unmarshal(m.Body, &msg); err != nil {
		return err
	}

	msg.Message = m

	ctx := context.Background()

	return h.Consume(ctx, &msg)
}
