package redis

import (
	"github.com/go-redis/cache/v9"
)

// NewCache from config.
// The cache is based on https://github.com/go-redis/cache
func NewCache(opts *cache.Options) *cache.Cache {
	return cache.New(opts)
}
