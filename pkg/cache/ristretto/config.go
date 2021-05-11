package ristretto

import (
	"github.com/dgraph-io/ristretto"
)

// NewConfig for ristretto.
func NewConfig() *ristretto.Config {
	cfg := &ristretto.Config{
		NumCounters: 1e7,     // nolint:gomnd
		MaxCost:     1 << 30, // nolint:gomnd
		BufferItems: 64,      // nolint:gomnd
		Metrics:     true,
	}

	return cfg
}
