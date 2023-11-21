package limiter

import (
	"context"

	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

// New limiter.
func New(formatted string) (*limiter.Limiter, error) {
	rate, err := limiter.NewRateFromFormatted(formatted)
	if err != nil {
		return nil, err
	}

	store := memory.NewStore()

	return limiter.New(store, rate), nil
}

// KeyFunc for the limiter.
type KeyFunc func(context.Context) string
