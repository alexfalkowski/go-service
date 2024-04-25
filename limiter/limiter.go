package limiter

import (
	"context"

	"github.com/alexfalkowski/go-service/meta"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

// New limiter.
//
//nolint:nilnil
func New(cfg *Config) (*limiter.Limiter, error) {
	if !IsEnabled(cfg) {
		return nil, nil
	}

	rate, err := limiter.NewRateFromFormatted(cfg.Pattern)
	if err != nil {
		return nil, err
	}

	store := memory.NewStore()

	return limiter.New(store, rate), nil
}

// KeyFunc for the limiter.
type KeyFunc func(context.Context) meta.Valuer
