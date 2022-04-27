package retry

import (
	"context"

	"github.com/alexfalkowski/go-service/transport/nsq/message"
	"github.com/alexfalkowski/go-service/transport/nsq/producer"
	retry "github.com/avast/retry-go/v3"
)

// NewProducer for retry.
func NewProducer(cfg *Config, p producer.Producer) *Producer {
	return &Producer{cfg: cfg, Producer: p}
}

// Producer for retry.
type Producer struct {
	cfg *Config
	producer.Producer
}

func (p *Producer) Publish(ctx context.Context, topic string, message *message.Message) error {
	operation := func() error {
		tctx, cancel := context.WithTimeout(ctx, p.cfg.Timeout)
		defer cancel()

		return p.Producer.Publish(tctx, topic, message)
	}

	return retry.Do(operation, retry.Attempts(p.cfg.Attempts))
}
