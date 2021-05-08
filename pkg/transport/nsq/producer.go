package nsq

import (
	"context"
	"encoding/json"

	"github.com/alexfalkowski/go-service/pkg/config"
	"github.com/alexfalkowski/go-service/pkg/transport/nsq/message"
	nsq "github.com/nsqio/go-nsq"
	"go.uber.org/fx"
)

// NewProducer for NSQ.
func NewProducer(lc fx.Lifecycle, systemConfig *config.Config, nsqConfig *nsq.Config) (*Producer, error) {
	producer, err := nsq.NewProducer(systemConfig.NSQHost, nsqConfig)
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			producer.Stop()

			return nil
		},
	})

	return &Producer{producer: producer}, nil
}

// Producer for NSQ.
type Producer struct {
	producer *nsq.Producer
}

// Publish a message to a topic.
func (p *Producer) Publish(topic string, message *message.Message) error {
	bytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return p.producer.Publish(topic, bytes)
}
