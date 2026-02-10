package breaker

import (
	"github.com/alexfalkowski/go-service/v2/breaker"
	"github.com/alexfalkowski/go-service/v2/net/grpc/status"
	"github.com/alexfalkowski/go-service/v2/sync"
)

type registry struct {
	opts     *opts
	breakers sync.Map[string, *breaker.CircuitBreaker]
}

func (r *registry) get(fullMethod string) *breaker.CircuitBreaker {
	if cb, ok := r.breakers.Load(fullMethod); ok {
		return cb
	}

	s := r.opts.settings
	s.Name = fullMethod

	failureCodes := r.opts.failureCodes
	s.IsSuccessful = func(err error) bool {
		if err != nil {
			_, isFailure := failureCodes[status.Code(err)]
			return !isFailure
		}

		return true
	}

	cb := breaker.NewCircuitBreaker(s)
	actual, _ := r.breakers.LoadOrStore(fullMethod, cb)
	return actual
}
