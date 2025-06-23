package psutil

import (
	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/shirou/gopsutil/v4/net"
)

// NewNet for debug.
func NewNet(ctx context.Context) *Net {
	counters, _ := net.IOCountersWithContext(ctx, true)

	return &Net{Counters: counters}
}

// Net for debug.
type Net struct {
	Counters []net.IOCountersStat `yaml:"counters,omitempty" json:"counters,omitempty" toml:"counters,omitempty"`
}
