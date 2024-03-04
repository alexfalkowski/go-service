package retry

import (
	"context"

	"github.com/alexfalkowski/go-service/nsq"
	"github.com/alexfalkowski/go-service/retry"
)

// NewProducer for retry.
func NewProducer(cfg *retry.Config, p nsq.Producer) *Producer {
	return &Producer{cfg: cfg, Producer: p}
}

// Producer for retry.
type Producer struct {
	cfg *retry.Config
	nsq.Producer
}

func (p *Producer) Produce(ctx context.Context, topic string, message *nsq.Message) error {
	operation := func() error {
		tctx, cancel := context.WithTimeout(ctx, p.cfg.Timeout)
		defer cancel()

		return p.Producer.Produce(tctx, topic, message)
	}

	return retry.Try(operation, p.cfg)
}
