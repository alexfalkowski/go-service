package nsq

import (
	nsq "github.com/nsqio/go-nsq"
)

// NewConfig for NSQ.
func NewConfig() *nsq.Config {
	cfg := nsq.NewConfig()

	return cfg
}
