package nsq

import (
	"context"

	"github.com/alexfalkowski/go-service/pkg/config"
	nsq "github.com/nsqio/go-nsq"
	"go.uber.org/fx"
)

// NewProducer for NSQ.
func NewProducer(lc fx.Lifecycle, systemConfig *config.Config, nsqConfig *nsq.Config) (*nsq.Producer, error) {
	producer, err := nsq.NewProducer(systemConfig.NSQHost, nsqConfig)
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			return producer.Ping()
		},
		OnStop: func(context.Context) error {
			producer.Stop()

			return nil
		},
	})

	return producer, nil
}
