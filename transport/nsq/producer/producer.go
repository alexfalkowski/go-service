package producer

import (
	"context"

	"github.com/alexfalkowski/go-service/transport/nsq/message"
)

// Producer for NSQ.
type Producer interface {
	Publish(ctx context.Context, topic string, message *message.Message) error
}
