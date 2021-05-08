package message

import (
	"github.com/golang/snappy"
	"github.com/nsqio/go-nsq"
	"github.com/vmihailenco/msgpack/v5"
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
func New(body []byte) *Message {
	message := &Message{
		Headers: Headers{},
		Body:    body,
	}

	return message
}

// Marshal a message.
func Marshal(message *Message) ([]byte, error) {
	b, err := msgpack.Marshal(message)
	if err != nil {
		return nil, err
	}

	b = snappy.Encode(nil, b)

	return b, nil
}

// Unmarshal a message.
func Unmarshal(data []byte, message *Message) error {
	b, err := snappy.Decode(nil, data)
	if err != nil {
		return err
	}

	return msgpack.Unmarshal(b, message)
}
