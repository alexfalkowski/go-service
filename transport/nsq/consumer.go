package nsq

import (
	"context"

	"github.com/alexfalkowski/go-service/transport/nsq/handler"
	lzap "github.com/alexfalkowski/go-service/transport/nsq/logger/zap"
	"github.com/alexfalkowski/go-service/transport/nsq/meta"
	"github.com/alexfalkowski/go-service/transport/nsq/trace/opentracing"
	"github.com/nsqio/go-nsq"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// ConsumerParams for NSQ.
type ConsumerParams struct {
	Config  *Config
	Logger  *zap.Logger
	Topic   string
	Channel string
	Handler handler.Handler
}

// RegisterConsumer for NSQ.
func RegisterConsumer(lc fx.Lifecycle, params *ConsumerParams) error {
	cfg := nsq.NewConfig()

	c, err := nsq.NewConsumer(params.Topic, params.Channel, cfg)
	if err != nil {
		return err
	}

	c.SetLogger(lzap.NewLogger(params.Logger), nsq.LogLevelInfo)

	lh := lzap.NewHandler(params.Topic, params.Channel, params.Logger, params.Handler)
	oh := opentracing.NewHandler(params.Topic, params.Channel, lh)
	mh := meta.NewHandler(oh)

	c.AddHandler(handler.New(mh))

	err = c.ConnectToNSQLookupd(params.Config.LookupHost)
	if err != nil {
		return err
	}

	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			c.Stop()

			return nil
		},
	})

	return nil
}
