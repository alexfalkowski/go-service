package breaker

import (
	"context"

	"github.com/alexfalkowski/go-service/transport/nsq/message"
	"github.com/alexfalkowski/go-service/transport/nsq/producer"
	breaker "github.com/sony/gobreaker"
)

// NewProducer for retry.
func NewProducer(p producer.Producer) *Producer {
	cb := breaker.NewCircuitBreaker(breaker.Settings{})

	return &Producer{cb: cb, Producer: p}
}

// Producer for retry.
type Producer struct {
	cb *breaker.CircuitBreaker
	producer.Producer
}

func (p *Producer) Publish(ctx context.Context, topic string, message *message.Message) error {
	operation := func() (any, error) {
		return nil, p.Producer.Publish(ctx, topic, message)
	}

	_, err := p.cb.Execute(operation)

	return err
}
