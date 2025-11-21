package driver

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type logger struct{}

func (logger) Printf(_ context.Context, _ string, _ ...any) {
	// Do nothing here
}

func init() {
	redis.SetLogger(&logger{})
}
