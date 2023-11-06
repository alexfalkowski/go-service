package nsq

import (
	"context"
)

// Producer for NSQ.
type Producer interface {
	Produce(ctx context.Context, topic string, message *Message) error
}
