package nsq

import (
	"context"

	lzap "github.com/alexfalkowski/go-service/transport/nsq/logger/zap"
	"github.com/alexfalkowski/go-service/transport/nsq/message"
	"github.com/alexfalkowski/go-service/transport/nsq/meta"
	"github.com/alexfalkowski/go-service/transport/nsq/producer"
	"github.com/alexfalkowski/go-service/transport/nsq/trace/opentracing"
	"github.com/nsqio/go-nsq"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// ProducerParams for NSQ.
type ProducerParams struct {
	Config *Config
	Logger *zap.Logger
}

// NewProducer for NSQ.
// nolint:ireturn
func NewProducer(lc fx.Lifecycle, params *ProducerParams) (producer.Producer, error) {
	cfg := nsq.NewConfig()

	p, err := nsq.NewProducer(params.Config.Host, cfg)
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
	pr = lzap.NewProducer(params.Logger, pr)
	pr = opentracing.NewProducer(pr)
	pr = meta.NewProducer(params.Config.UserAgent, pr)

	return pr, nil
}

type nsqProducer struct {
	*nsq.Producer
}

// Publish a message to a topic.
func (p *nsqProducer) Publish(ctx context.Context, topic string, msg *message.Message) (context.Context, error) {
	bytes, err := message.Marshal(msg)
	if err != nil {
		return ctx, err
	}

	return ctx, p.Producer.Publish(topic, bytes)
}
