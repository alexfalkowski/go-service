package nsq

import (
	"context"
	"encoding/json"

	"github.com/alexfalkowski/go-service/pkg/config"
	"github.com/alexfalkowski/go-service/pkg/transport/nsq/message"
	"github.com/alexfalkowski/go-service/pkg/transport/nsq/producer"
	"github.com/nsqio/go-nsq"
	"go.uber.org/fx"
)

// NewProducer for NSQ.
func NewProducer(lc fx.Lifecycle, systemConfig *config.Config, nsqConfig *nsq.Config) (producer.Producer, error) {
	p, err := nsq.NewProducer(systemConfig.NSQHost, nsqConfig)
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			p.Stop()

			return nil
		},
	})

	producer := &nsqProducer{Producer: p}

	return producer, nil
}

type nsqProducer struct {
	*nsq.Producer
}

// Publish a message to a topic.
func (p *nsqProducer) Publish(topic string, message *message.Message) error {
	bytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return p.Producer.Publish(topic, bytes)
}
