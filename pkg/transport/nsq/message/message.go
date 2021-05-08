package message

import "github.com/nsqio/go-nsq"

// Headers of message.
type Headers map[string]string

// Message for NSQ.
type Message struct {
	Headers Headers `json:"headers"`
	Body    []byte  `json:"body"`

	*nsq.Message `json:"-"`
}

// New message.
func New(body []byte) *Message {
	message := &Message{
		Headers: make(Headers),
		Body:    body,
	}

	return message
}
