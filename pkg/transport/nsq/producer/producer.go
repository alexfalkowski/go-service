package producer

import (
	"github.com/alexfalkowski/go-service/pkg/transport/nsq/message"
)

// Producer for NSQ.
type Producer interface {
	Publish(topic string, message *message.Message) error
}
