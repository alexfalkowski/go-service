package redis

import (
	"github.com/go-redis/cache/v8"
)

// NewCache from config.
// The cache is based on https://github.com/go-redis/cache
func NewCache(opts *cache.Options) *cache.Cache {
	cache := cache.New(opts)

	return cache
}
