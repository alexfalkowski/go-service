package retry

import (
	"context"

	"github.com/alexfalkowski/go-service/nsq"
	retry "github.com/avast/retry-go/v3"
)

// NewProducer for retry.
func NewProducer(cfg *Config, p nsq.Producer) *Producer {
	return &Producer{cfg: cfg, Producer: p}
}

// Producer for retry.
type Producer struct {
	cfg *Config
	nsq.Producer
}

func (p *Producer) Produce(ctx context.Context, topic string, message *nsq.Message) error {
	operation := func() error {
		tctx, cancel := context.WithTimeout(ctx, p.cfg.Timeout)
		defer cancel()

		return p.Producer.Produce(tctx, topic, message)
	}

	return retry.Do(operation, retry.Attempts(p.cfg.Attempts))
}
