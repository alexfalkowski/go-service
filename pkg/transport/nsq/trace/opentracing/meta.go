package opentracing

import (
	"github.com/alexfalkowski/go-service/pkg/transport/nsq/message"
)

type headersTextMap message.Headers

func (c headersTextMap) ForeachKey(handler func(key, val string) error) error {
	for k, v := range c {
		if err := handler(k, v); err != nil {
			return err
		}
	}

	return nil
}

func (c headersTextMap) Set(key, val string) {
	c[key] = val
}
