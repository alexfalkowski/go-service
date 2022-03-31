package redis

import (
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"github.com/golang/snappy"
	"google.golang.org/protobuf/proto"
)

// NewOptions for redis.
// nolint:forcetypeassert
func NewOptions(ring *redis.Ring) *cache.Options {
	opts := &cache.Options{
		Redis:        ring,
		StatsEnabled: true,
		Marshal: func(v any) ([]byte, error) {
			m, err := proto.Marshal(v.(proto.Message))
			if err != nil {
				return nil, err
			}

			return snappy.Encode(nil, m), nil
		},
		Unmarshal: func(b []byte, v any) error {
			m, err := snappy.Decode(nil, b)
			if err != nil {
				return err
			}

			return proto.Unmarshal(m, v.(proto.Message))
		},
	}

	return opts
}
