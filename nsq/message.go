package nsq

import (
	"github.com/nsqio/go-nsq"
)

// Headers of message.
type Headers map[string]string

// Message for NSQ.
type Message struct {
	Headers Headers `msgpack:"headers"`
	Body    []byte  `msgpack:"body"`

	*nsq.Message `msgpack:"-"`
}

// New message.
func NewMessage(body []byte) *Message {
	message := &Message{
		Headers: Headers{},
		Body:    body,
	}

	return message
}
