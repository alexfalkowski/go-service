package nsq

import (
	"context"

	"github.com/alexfalkowski/go-service/pkg/config"
	"github.com/alexfalkowski/go-service/pkg/transport/nsq/handler"
	pkgZap "github.com/alexfalkowski/go-service/pkg/transport/nsq/logger/zap"
	"github.com/alexfalkowski/go-service/pkg/transport/nsq/trace/opentracing"
	"github.com/nsqio/go-nsq"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// ConsumerParams for NSQ.
type ConsumerParams struct {
	SystemConfig *config.Config
	NSQConfig    *nsq.Config
	Logger       *zap.Logger
	Topic        string
	Channel      string
	Handler      handler.Handler
}

// RegisterConsumer for NSQ.
func RegisterConsumer(lc fx.Lifecycle, params *ConsumerParams) error {
	consumer, err := nsq.NewConsumer(params.Topic, params.Channel, params.NSQConfig)
	if err != nil {
		return err
	}

	h := pkgZap.NewHandler(params.Logger, params.Handler)
	h = opentracing.NewHandler(h)

	consumer.AddHandler(handler.NewHandler(params.Topic, params.Channel, h))

	err = consumer.ConnectToNSQLookupd(params.SystemConfig.NSQLookupHost)
	if err != nil {
		return err
	}

	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			consumer.Stop()

			return nil
		},
	})

	return nil
}
