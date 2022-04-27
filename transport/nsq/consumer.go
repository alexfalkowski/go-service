package nsq

import (
	"context"

	"github.com/alexfalkowski/go-service/transport/nsq/handler"
	"github.com/alexfalkowski/go-service/transport/nsq/logger"
	lzap "github.com/alexfalkowski/go-service/transport/nsq/logger/zap"
	"github.com/alexfalkowski/go-service/transport/nsq/marshaller"
	"github.com/alexfalkowski/go-service/transport/nsq/meta"
	"github.com/alexfalkowski/go-service/transport/nsq/trace/opentracing"
	"github.com/nsqio/go-nsq"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// ConsumerParams for NSQ.
type ConsumerParams struct {
	Lifecycle  fx.Lifecycle
	Config     *Config
	Logger     *zap.Logger
	Topic      string
	Channel    string
	Tracer     opentracing.Tracer
	Handler    handler.Handler
	Marshaller marshaller.Marshaller
}

// RegisterConsumer for NSQ.
func RegisterConsumer(params ConsumerParams) error {
	cfg := nsq.NewConfig()

	c, err := nsq.NewConsumer(params.Topic, params.Channel, cfg)
	if err != nil {
		return err
	}

	c.SetLogger(logger.NewLogger(), nsq.LogLevelInfo)

	var h handler.Handler = lzap.NewHandler(params.Topic, params.Channel, params.Logger, params.Handler)
	h = opentracing.NewHandler(params.Topic, params.Channel, params.Tracer, h)
	h = meta.NewHandler(h)

	c.AddHandler(handler.New(handler.Params{Handler: h, Marshaller: params.Marshaller}))

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
