package nsq

import (
	"context"

	"github.com/alexfalkowski/go-service/transport/nsq/handler"
	lzap "github.com/alexfalkowski/go-service/transport/nsq/logger/zap"
	"github.com/alexfalkowski/go-service/transport/nsq/meta"
	nopentracing "github.com/alexfalkowski/go-service/transport/nsq/trace/opentracing"
	"github.com/nsqio/go-nsq"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// ConsumerParams for NSQ.
type ConsumerParams struct {
	Lifecycle fx.Lifecycle
	Config    *Config
	Logger    *zap.Logger
	Topic     string
	Channel   string
	Tracer    opentracing.Tracer
	Handler   handler.Handler
}

// RegisterConsumer for NSQ.
func RegisterConsumer(params *ConsumerParams) error {
	cfg := nsq.NewConfig()

	c, err := nsq.NewConsumer(params.Topic, params.Channel, cfg)
	if err != nil {
		return err
	}

	c.SetLogger(lzap.NewLogger(params.Logger), nsq.LogLevelInfo)

	lh := lzap.NewHandler(params.Topic, params.Channel, params.Logger, params.Handler)
	oh := nopentracing.NewHandler(params.Topic, params.Channel, params.Tracer, lh)
	mh := meta.NewHandler(oh)

	c.AddHandler(handler.New(mh))

	err = c.ConnectToNSQLookupd(params.Config.LookupHost)
	if err != nil {
		return err
	}

	params.Lifecycle.Append(fx.Hook{
		OnStop: func(context.Context) error {
			c.Stop()

			return nil
		},
	})

	return nil
}
