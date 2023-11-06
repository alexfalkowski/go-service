package breaker

import (
	"context"

	"github.com/alexfalkowski/go-service/nsq"
	breaker "github.com/sony/gobreaker"
)

// NewProducer for retry.
func NewProducer(p nsq.Producer) *Producer {
	cb := breaker.NewCircuitBreaker(breaker.Settings{})

	return &Producer{cb: cb, Producer: p}
}

// Producer for retry.
type Producer struct {
	cb *breaker.CircuitBreaker
	nsq.Producer
}

func (p *Producer) Produce(ctx context.Context, topic string, message *nsq.Message) error {
	operation := func() (any, error) {
		return nil, p.Producer.Produce(ctx, topic, message)
	}

	_, err := p.cb.Execute(operation)

	return err
}
