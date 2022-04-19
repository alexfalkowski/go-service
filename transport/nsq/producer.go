package nsq

import (
	"context"

	lzap "github.com/alexfalkowski/go-service/transport/nsq/logger/zap"
	"github.com/alexfalkowski/go-service/transport/nsq/message"
	"github.com/alexfalkowski/go-service/transport/nsq/meta"
	"github.com/alexfalkowski/go-service/transport/nsq/producer"
	nopentracing "github.com/alexfalkowski/go-service/transport/nsq/trace/opentracing"
	"github.com/nsqio/go-nsq"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// ProducerParams for NSQ.
type ProducerParams struct {
	Lifecycle fx.Lifecycle
	Config    *Config
	Logger    *zap.Logger
	Tracer    opentracing.Tracer
}

// NewProducer for NSQ.
// nolint:ireturn
func NewProducer(params *ProducerParams) (producer.Producer, error) {
	cfg := nsq.NewConfig()

	p, err := nsq.NewProducer(params.Config.Host, cfg)
	if err != nil {
		return nil, err
	}

	p.SetLogger(lzap.NewLogger(params.Logger), nsq.LogLevelInfo)

	params.Lifecycle.Append(fx.Hook{
		OnStop: func(context.Context) error {
			p.Stop()

			return nil
		},
	})

	var pr producer.Producer = &nsqProducer{Producer: p}
	pr = lzap.NewProducer(params.Logger, pr)
	pr = nopentracing.NewProducer(params.Tracer, pr)
	pr = meta.NewProducer(params.Config.UserAgent, pr)

	return pr, nil
}

type nsqProducer struct {
	*nsq.Producer
}

// Publish a message to a topic.
func (p *nsqProducer) Publish(ctx context.Context, topic string, msg *message.Message) error {
	bytes, err := message.Marshal(msg)
	if err != nil {
		return err
	}

	return p.Producer.Publish(topic, bytes)
}
