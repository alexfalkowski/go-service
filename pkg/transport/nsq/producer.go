package nsq

import (
	"context"
	"encoding/json"

	"github.com/alexfalkowski/go-service/pkg/config"
	pkgZap "github.com/alexfalkowski/go-service/pkg/transport/nsq/logger/zap"
	"github.com/alexfalkowski/go-service/pkg/transport/nsq/message"
	"github.com/alexfalkowski/go-service/pkg/transport/nsq/producer"
	"github.com/nsqio/go-nsq"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// ProducerParams for NSQ.
type ProducerParams struct {
	SystemConfig *config.Config
	NSQConfig    *nsq.Config
	Logger       *zap.Logger
}

// NewProducer for NSQ.
func NewProducer(lc fx.Lifecycle, params *ProducerParams) (producer.Producer, error) {
	p, err := nsq.NewProducer(params.SystemConfig.NSQHost, params.NSQConfig)
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			p.Stop()

			return nil
		},
	})

	var pr producer.Producer = &nsqProducer{Producer: p}
	pr = pkgZap.NewProducer(params.Logger, pr)

	return pr, nil
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
